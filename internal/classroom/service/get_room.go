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

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *ClassroomService) GetEmptyRoom(req *classroom.EmptyRoomRequest) ([]*model.Classroom, error) {
	start, end, err := parseSections(req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	day, err := s.cache.Classroom.GetAvailability(s.ctx, req.Date, req.Campus)
	if err != nil {
		if errors.Is(err, classroomCache.ErrAvailabilityFailure) {
			return nil, errno.BizError.WithMessage("教室数据拉取失败，请稍后重试")
		}
		if errors.Is(err, redis.Nil) {
			return nil, errno.BizError.WithMessage("暂无该日教室数据")
		}
		return nil, fmt.Errorf("service.GetEmptyRoom: get availability failed: %w", err)
	}
	if day.Rooms == nil {
		return nil, fmt.Errorf("invalid classroom availability snapshot")
	}
	return filterRooms(day.Rooms, start, end, req.GetIncludeOccupied()), nil
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

func filterRooms(rooms []*model.Classroom, start, end int32, includeOccupied bool) []*model.Classroom {
	result := make([]*model.Classroom, 0, len(rooms))
	for _, room := range rooms {
		if includeOccupied {
			result = append(result, room)
			continue
		}
		free := true
		for _, occupancy := range room.Occupancies {
			if occupancy.Start <= end && occupancy.End >= start {
				free = false
				break
			}
		}
		if !free {
			continue
		}
		result = append(result, room)
	}
	return result
}
