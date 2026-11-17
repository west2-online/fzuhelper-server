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
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *ClassroomService) GetEmptyRoom(req *classroom.EmptyRoomRequest) ([]*model.Classroom, error) {
	start, end, err := parseSections(req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	rooms, err := s.getAvailability(req.Date, req.Campus)
	if err != nil {
		return nil, err
	}
	emptyRooms := make([]*model.Classroom, 0, len(rooms))
	for _, room := range rooms {
		for _, free := range room.FreeRanges {
			if free.Start <= start && free.End >= end {
				emptyRooms = append(emptyRooms, &model.Classroom{
					Build: room.Build, Location: room.Location, Capacity: room.Capacity, Type: room.Type,
				})
				break
			}
		}
	}
	return emptyRooms, nil
}

func (s *ClassroomService) GetEmptyRoomV2(req *classroom.EmptyRoomV2Request) ([]*model.Classroom, error) {
	return s.getAvailability(req.Date, req.Campus)
}

func (s *ClassroomService) getAvailability(date, campus string) ([]*model.Classroom, error) {
	const daysPerWeek = 7
	requestDate, err := time.ParseInLocation(time.DateOnly, date, constants.ChinaTZ)
	if err != nil {
		return nil, fmt.Errorf("invalid classroom date: %w", err)
	}
	now := time.Now().In(constants.ChinaTZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, constants.ChinaTZ)
	maxDays := daysPerWeek + (daysPerWeek-int(today.Weekday()))%daysPerWeek
	if requestDate.Before(today) || requestDate.After(today.AddDate(0, 0, maxDays)) {
		return nil, fmt.Errorf("date out of range, date: %s", date)
	}
	key := s.cache.Classroom.AvailabilityKey(date, campus)
	rooms, err := s.cache.Classroom.GetAvailabilityCache(s.ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errno.BizError.WithMessage("暂无该日教室数据")
		}
		return nil, fmt.Errorf("service.GetEmptyRoom: get availability failed: %w", err)
	}
	if rooms == nil {
		return nil, errno.BizError.WithMessage("教室数据拉取失败，请稍后重试")
	}
	return rooms, nil
}

func parseSections(startTime, endTime string) (int32, int32, error) {
	const lastSection = 11
	start, startErr := strconv.Atoi(startTime)
	end, endErr := strconv.Atoi(endTime)
	if startErr != nil || endErr != nil || start < 1 || end > lastSection || start > end {
		return 0, 0, fmt.Errorf("invalid classroom sections: require 1 <= startTime <= endTime <= 11")
	}
	return int32(start), int32(end), nil
}
