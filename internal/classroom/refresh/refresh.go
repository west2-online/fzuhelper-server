/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package refresh

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/internal/classroom/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

type Source interface {
	GetClassroomBuildings(ctx context.Context, campus string) ([]string, error)
	GetClassroomUsage(ctx context.Context, req jwch.ClassroomUsageReq) (*jwch.ClassroomUsage, error)
}

type Cache interface {
	SetAvailability(ctx context.Context, campus string, day *classroomCache.DayAvailability) error
	SetAvailabilityFailure(ctx context.Context, date, campus string) error
}

type Refresher struct {
	Source       Source
	Cache        Cache
	RequestDelay time.Duration
	RetryDelay   time.Duration
}

// target ：日期 x 校区；pending 表示尚未成功抓取的楼。
type target struct {
	date       string
	campus     string
	term       string
	pending    []string
	rooms      []*jwch.ClassroomUsageRoom
	published  bool
	listErr    error
	publishErr error
}

type campusBuildings struct {
	buildings []string
	err       error
}

const attempts = 2

func (r *Refresher) RunScheduled(ctx context.Context, now time.Time, campuses []string) error {
	dates := Dates(now)
	var eg errgroup.Group
	for _, date := range dates {
		currentDate := date
		eg.Go(func() error {
			dateString := currentDate.Format(time.DateOnly)
			if err := r.Run(ctx, []string{dateString}, campuses); err != nil {
				logger.WithCtx(ctx).Errorf("refresh classrooms: date %s: %v", dateString, err)
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("refresh scheduled classrooms: %w", err)
	}
	logger.WithCtx(ctx).Infof("refresh scheduled classrooms: complete %d dates from %s to %s",
		len(dates), dates[0].Format(time.DateOnly), dates[len(dates)-1].Format(time.DateOnly))
	return nil
}

func (r *Refresher) Run(ctx context.Context, dates, campuses []string) error {
	lists, err := r.loadCampusLists(ctx, campuses)
	if err != nil {
		return err
	}
	targets := make([]*target, 0, len(dates)*len(campuses))
	for _, date := range dates {
		for _, campus := range campuses {
			item := &target{date: date, campus: campus, listErr: lists[campus].err}
			if item.listErr == nil {
				item.pending = append([]string(nil), lists[campus].buildings...)
			}
			targets = append(targets, item)
		}
	}
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			if err := r.wait(ctx, r.RetryDelay); err != nil {
				return err
			}
		}
		for _, item := range targets {
			if err := ctx.Err(); err != nil {
				return err
			}
			if item.published || item.listErr != nil {
				continue
			}
			r.fetchBuildings(ctx, item)
			if err := ctx.Err(); err != nil {
				return err
			}
			if len(item.pending) == 0 || attempt == attempts-1 {
				item.publishErr = r.publish(ctx, item)
				item.published = item.publishErr == nil
			}
		}
		if allPublished(targets) {
			return nil
		}
	}
	return r.reportFailures(ctx, targets)
}

func (r *Refresher) reportFailures(ctx context.Context, targets []*target) error {
	failures := make([]error, 0)
	for _, item := range targets {
		switch {
		case item.listErr != nil:
			logger.WithCtx(ctx).Errorf("refresh classrooms %s %s: %v", item.date, item.campus, item.listErr)
			if err := r.setAvailabilityFailure(ctx, item.date, item.campus); err != nil {
				failures = append(failures, err)
			}
		case !item.published:
			failures = append(failures, fmt.Errorf("publish classrooms %s %s: %w", item.date, item.campus, item.publishErr))
			if err := r.setAvailabilityFailure(ctx, item.date, item.campus); err != nil {
				failures = append(failures, err)
			}
		case len(item.pending) > 0:
			logger.WithCtx(ctx).Errorf("classrooms %s %s failed buildings: %v", item.date, item.campus, item.pending)
		}
	}
	return errors.Join(append(failures, ctx.Err())...)
}

func (r *Refresher) loadCampusLists(ctx context.Context, campuses []string) (map[string]campusBuildings, error) {
	result := make(map[string]campusBuildings, len(campuses))
	for _, campus := range campuses {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		buildings, err := r.loadBuildings(ctx, campus)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			result[campus] = campusBuildings{err: err}
			continue
		}
		result[campus] = campusBuildings{buildings: buildings}
	}
	return result, nil
}

// loadBuildings 获取单个校区的教学楼列表
func (r *Refresher) loadBuildings(ctx context.Context, campus string) ([]string, error) {
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if attempt > 0 {
			if err := r.wait(ctx, r.RetryDelay); err != nil {
				return nil, err
			}
		}
		buildings, err := r.Source.GetClassroomBuildings(ctx, campus)
		if err == nil {
			if len(buildings) > 0 {
				return buildings, nil
			}
			err = fmt.Errorf("no classroom buildings for %s", campus)
		}
		lastErr = err
	}
	return nil, lastErr
}

func (r *Refresher) fetchBuildings(ctx context.Context, item *target) {
	failed := make([]string, 0)
	for _, building := range item.pending {
		if err := r.wait(ctx, r.RequestDelay); err != nil {
			return
		}
		usage, err := r.Source.GetClassroomUsage(ctx, jwch.ClassroomUsageReq{
			Date: item.date, Campus: item.campus, Term: item.term, Buildings: []string{building},
		})
		if err == nil {
			err = validateUsage(item, usage, building)
		}
		if err != nil {
			failed = append(failed, building)
			continue
		}
		item.term = usage.Term
		item.rooms = append(item.rooms, usage.Rooms...)
	}
	item.pending = failed
}

func validateUsage(item *target, usage *jwch.ClassroomUsage, building string) error {
	seen := make(map[string]bool, len(item.rooms))
	for _, room := range item.rooms {
		seen[room.Location] = true
	}
	for _, room := range usage.Rooms {
		if room.Error != nil {
			return fmt.Errorf("query building %s: %w", building, room.Error)
		}
		if seen[room.Location] {
			return fmt.Errorf("duplicate classroom %q in building %q", room.Location, building)
		}
	}
	return nil
}

func (r *Refresher) publish(ctx context.Context, item *target) error {
	if len(item.pending) > 0 {
		return r.setAvailabilityFailure(ctx, item.date, item.campus)
	}
	day := &classroomCache.DayAvailability{
		Date: item.date, Term: item.term, Rooms: pack.BuildUsageRooms(item.rooms, item.campus),
	}
	for campus, snapshot := range splitAvailability(item.campus, day) {
		if err := r.Cache.SetAvailability(ctx, campus, snapshot); err != nil {
			return err
		}
	}
	return nil
}

func (r *Refresher) setAvailabilityFailure(ctx context.Context, date, campus string) error {
	failures := make([]error, 0)
	for _, targetCampus := range availabilityCampuses(campus) {
		if err := r.Cache.SetAvailabilityFailure(ctx, date, targetCampus); err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}

func availabilityCampuses(campus string) []string {
	if campus == "厦门工艺美院" {
		return []string{"鼓浪屿校区", "集美校区"}
	}
	return []string{campus}
}

func splitAvailability(campus string, day *classroomCache.DayAvailability) map[string]*classroomCache.DayAvailability {
	if campus != "厦门工艺美院" {
		return map[string]*classroomCache.DayAvailability{campus: day}
	}
	guLangYu := &classroomCache.DayAvailability{Date: day.Date, Term: day.Term, Rooms: make([]*model.Classroom, 0)}
	jiMei := &classroomCache.DayAvailability{Date: day.Date, Term: day.Term, Rooms: make([]*model.Classroom, 0)}
	for _, room := range day.Rooms {
		if strings.Contains(room.Location, "鼓浪屿") {
			copyRoom := *room
			copyRoom.Build = "鼓浪屿校区"
			guLangYu.Rooms = append(guLangYu.Rooms, &copyRoom)
		} else if strings.Contains(room.Location, "集美") {
			copyRoom := *room
			copyRoom.Build = "集美校区"
			jiMei.Rooms = append(jiMei.Rooms, &copyRoom)
		}
	}
	return map[string]*classroomCache.DayAvailability{"鼓浪屿校区": guLangYu, "集美校区": jiMei}
}

func (r *Refresher) wait(ctx context.Context, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func allPublished(targets []*target) bool {
	for _, item := range targets {
		if !item.published || len(item.pending) > 0 {
			return false
		}
	}
	return true
}
