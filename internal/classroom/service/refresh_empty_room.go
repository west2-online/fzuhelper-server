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

package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/classroom/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

const (
	attempts      = 2
	retryInterval = time.Second
)

func (s *ClassroomService) RefreshEmptyRoom(dates []string) error {
	stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
	if err := utils.RetryLogin(stu); err != nil {
		return fmt.Errorf("service.RefreshEmptyRoom: failed to login: %w", err)
	}
	for _, date := range dates {
		s.refreshEmptyRoomDate(stu, date, constants.CampusArray)
	}
	logger.WithCtx(s.ctx).Infof("refresh classrooms: complete %d dates", len(dates))
	return nil
}

func (s *ClassroomService) refreshEmptyRoomDate(stu *jwch.Student, date string, campuses []string) {
	pending := make([]jwch.ClassroomUsageReq, 0, len(campuses))
	roomsByCampus := make(map[string][]*jwch.ClassroomUsageRoom, len(campuses))
	for _, campus := range campuses {
		buildings, err := s.loadBuildings(stu, campus)
		if err != nil {
			logger.WithCtx(s.ctx).Errorf("refresh classrooms %s %s: %v", date, campus, err)
			s.saveAvailability(date, campus, nil)
			continue
		}
		pending = append(pending, jwch.ClassroomUsageReq{Date: date, Campus: campus, Buildings: buildings})
	}

	for attempt := 0; attempt < attempts && len(pending) > 0; attempt++ {
		if attempt > 0 {
			time.Sleep(retryInterval)
		}
		retry := make([]jwch.ClassroomUsageReq, 0, len(pending))
		for _, req := range pending {
			usage, err := stu.GetClassroomUsage(s.ctx, req)
			if err != nil {
				logger.WithCtx(s.ctx).Errorf("query classrooms %s %s: %v", date, req.Campus, err)
				if attempt < attempts-1 {
					retry = append(retry, req)
					continue
				}
				logger.WithCtx(s.ctx).Errorf("classrooms %s %s failed buildings: %v", date, req.Campus, req.Buildings)
				s.saveAvailability(date, req.Campus, nil)
				continue
			}

			req.Term = usage.Term
			req.Buildings = nil
			for _, room := range usage.Rooms {
				if room.Error != nil {
					req.Buildings = append(req.Buildings, room.Building)
					continue
				}
				roomsByCampus[req.Campus] = append(roomsByCampus[req.Campus], room)
			}
			if len(req.Buildings) == 0 {
				s.saveAvailability(date, req.Campus, pack.BuildUsageRooms(roomsByCampus[req.Campus], req.Campus))
				continue
			}
			if attempt < attempts-1 {
				retry = append(retry, req)
				continue
			}
			logger.WithCtx(s.ctx).Errorf("classrooms %s %s failed buildings: %v", date, req.Campus, req.Buildings)
			s.saveAvailability(date, req.Campus, nil)
		}
		pending = retry
	}
}

func (s *ClassroomService) loadBuildings(stu *jwch.Student, campus string) ([]string, error) {
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			time.Sleep(retryInterval)
		}
		buildings, err := stu.GetClassroomBuildings(s.ctx, campus)
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

func (s *ClassroomService) saveAvailability(date, campus string, rooms []*model.Classroom) {
	c := s.cache.Classroom
	if campus != "厦门工艺美院" {
		key := c.AvailabilityKey(date, campus)
		if err := c.SetAvailabilityCache(s.ctx, key, rooms); err != nil {
			logger.WithCtx(s.ctx).Errorf("save classroom cache %s: %v", key, err)
		}
		return
	}

	var guLangYu, jiMei []*model.Classroom
	if rooms != nil {
		guLangYu = make([]*model.Classroom, 0)
		jiMei = make([]*model.Classroom, 0)
	}
	for _, room := range rooms {
		if strings.Contains(room.Location, "鼓浪屿") {
			room.Build = "鼓浪屿校区"
			guLangYu = append(guLangYu, room)
		} else if strings.Contains(room.Location, "集美") {
			room.Build = "集美校区"
			jiMei = append(jiMei, room)
		}
	}
	guLangYuKey := c.AvailabilityKey(date, "鼓浪屿校区")
	jiMeiKey := c.AvailabilityKey(date, "集美校区")
	if err := c.SetAvailabilityCache(s.ctx, guLangYuKey, guLangYu); err != nil {
		logger.WithCtx(s.ctx).Errorf("save classroom cache %s: %v", guLangYuKey, err)
	}
	if err := c.SetAvailabilityCache(s.ctx, jiMeiKey, jiMei); err != nil {
		logger.WithCtx(s.ctx).Errorf("save classroom cache %s: %v", jiMeiKey, err)
	}
}
