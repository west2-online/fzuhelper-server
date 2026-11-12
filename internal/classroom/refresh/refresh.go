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
	"time"

	"github.com/west2-online/fzuhelper-server/internal/classroom/pack"
	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/jwch"
)

type Source interface {
	GetClassroomBuildings(ctx context.Context, campus string) ([]string, error)
	GetClassroomUsage(ctx context.Context, req jwch.ClassroomUsageReq) (*jwch.ClassroomUsage, error)
}

type Cache interface {
	SetAvailability(ctx context.Context, campus string, day *classroomCache.DayAvailability, failedBuildings ...string) error
	MarkAvailabilityFailure(ctx context.Context, date, campus string, failedBuildings ...string) error
}

type Refresher struct {
	Source       Source
	Cache        Cache
	RequestDelay time.Duration
	RetryDelay   time.Duration
}

// target 是一个「日期 × 校区」的刷新单元；pending 始终表示尚未成功抓取的楼。
type target struct {
	date       string
	campus     string
	term       string
	pending    []string
	rooms      []*jwch.ClassroomUsageRoom
	fetched    bool // 至少有一栋楼抓取成功
	published  bool
	listErr    error // 楼列表加载失败，该目标整轮跳过
	publishErr error
}

// campusBuildings 是单个校区的教学楼列表查询结果。
type campusBuildings struct {
	buildings []string
	err       error
}

// attempts 是每个上游环节的最大尝试次数（首次 + 一次补抓）。
const attempts = 2

func (r *Refresher) Run(ctx context.Context, dates, campuses []string) error {
	// 楼列表是校区级数据：先统一取一次，再展开成各日期校区的刷新单元。
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
			// 首轮只发布完整快照；最后一轮才允许发布残缺快照并附带失败标记。
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

// reportFailures 汇总各类失败；列表失败、发布失败和抓取不全各自只报一次。
func (r *Refresher) reportFailures(ctx context.Context, targets []*target) error {
	failures := make([]error, 0)
	for _, item := range targets {
		switch {
		case item.listErr != nil:
			failures = append(failures, fmt.Errorf("refresh classrooms %s %s: %w", item.date, item.campus, item.listErr))
			if err := r.Cache.MarkAvailabilityFailure(ctx, item.date, item.campus); err != nil {
				failures = append(failures, err)
			}
		case !item.published:
			failures = append(failures, fmt.Errorf("publish classrooms %s %s: %w", item.date, item.campus, item.publishErr))
			if err := r.Cache.MarkAvailabilityFailure(ctx, item.date, item.campus, item.pending...); err != nil {
				failures = append(failures, err)
			}
		case len(item.pending) > 0:
			failures = append(failures, fmt.Errorf("classrooms %s %s failed buildings: %v", item.date, item.campus, item.pending))
		}
	}
	return errors.Join(failures...)
}

// loadCampusLists 每个校区只取一次教学楼列表；单个校区失败不影响其他校区。
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

// loadBuildings 获取单个校区的教学楼列表，失败时空转 RetryDelay 后补抓一次。
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
			buildings = uniqueBuildings(buildings)
			if len(buildings) > 0 {
				return buildings, nil
			}
			// 空列表按失败处理，避免把上游异常发布成"当天没有教室"。
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
		item.fetched = true
		item.rooms = append(item.rooms, usage.Rooms...)
	}
	item.pending = failed
}

func validateUsage(item *target, usage *jwch.ClassroomUsage, building string) error {
	if usage == nil || usage.Date != item.date || usage.Term == "" || usage.Rooms == nil {
		return fmt.Errorf("invalid classroom usage response")
	}
	if item.term != "" && item.term != usage.Term {
		return fmt.Errorf("classroom usage academic term changed during refresh")
	}
	// 当前楼先完整校验再合并，重复教室不能被静默覆盖或混入快照。
	seen := make(map[string]bool, len(item.rooms)+len(usage.Rooms))
	for _, room := range item.rooms {
		seen[room.Location] = true
	}
	for _, room := range usage.Rooms {
		if room == nil {
			return fmt.Errorf("nil classroom usage room")
		}
		// jwch 的调用成功不代表楼查询成功；单楼错误放在返回记录中。
		if room.Error != nil {
			return fmt.Errorf("query building %s: %w", building, room.Error)
		}
		if room.Building != building {
			return fmt.Errorf("unexpected building %q in response for %q", room.Building, building)
		}
		if room.Location == "" {
			return fmt.Errorf("empty classroom location in building %q", building)
		}
		if seen[room.Location] {
			return fmt.Errorf("duplicate classroom %q in building %q", room.Location, building)
		}
		for _, occupancy := range room.Occupancies {
			if occupancy == nil || occupancy.Start < 1 || occupancy.End > 11 || occupancy.Start > occupancy.End {
				return fmt.Errorf("invalid classroom occupancy for %s", room.Location)
			}
		}
		seen[room.Location] = true
	}
	return nil
}

func (r *Refresher) publish(ctx context.Context, item *target) error {
	if !item.fetched {
		// 全部楼都失败时保留原快照，只更新失败状态。
		return r.Cache.MarkAvailabilityFailure(ctx, item.date, item.campus, item.pending...)
	}
	return r.Cache.SetAvailability(ctx, item.campus, &classroomCache.DayAvailability{
		Date: item.date, Term: item.term, Rooms: pack.BuildUsageRooms(item.rooms, item.campus),
	}, item.pending...)
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

func uniqueBuildings(buildings []string) []string {
	result := make([]string, 0, len(buildings))
	seen := make(map[string]bool, len(buildings))
	for _, building := range buildings {
		if building != "" && !seen[building] {
			result = append(result, building)
			seen[building] = true
		}
	}
	return result
}
