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

package pack

import (
	"strings"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
)

// BuildUsageRooms 保留旧接口的教室信息，并转换当天占用记录和连续空闲节次。
func BuildUsageRooms(rooms []*jwch.ClassroomUsageRoom, campus string) []*model.Classroom {
	result := make([]*model.Classroom, 0, len(rooms))
	for _, room := range rooms {
		build := campus
		if strings.Contains(campus, "旗山") {
			build = location2Build(room.Location)
		}
		occupancies := make([]*model.ClassroomOccupancy, 0, len(room.Occupancies))
		for _, occupancy := range room.Occupancies {
			occupancies = append(occupancies, &model.ClassroomOccupancy{
				Start: occupancy.Start, End: occupancy.End, Event: occupancy.Event,
				Occupant: occupancy.Occupant, RawText: occupancy.RawText,
				StartWeek: occupancy.StartWeek, EndWeek: occupancy.EndWeek,
			})
		}
		result = append(result, &model.Classroom{
			Build: build, Location: room.Location, Capacity: room.Capacity, Type: room.Type,
			Occupancies: occupancies, FreeRanges: BuildFreeRanges(occupancies),
		})
	}
	return result
}

// BuildFreeRanges 根据占用记录计算当天连续空闲节次，起止节均包含在内。
func BuildFreeRanges(occupancies []*model.ClassroomOccupancy) []*model.ClassroomSectionRange {
	const sectionsPerDay = 11
	var occupied [sectionsPerDay + 1]bool
	for _, occupancy := range occupancies {
		start, end := max(occupancy.Start, 1), min(occupancy.End, sectionsPerDay)
		for section := start; section <= end; section++ {
			occupied[section] = true
		}
	}
	ranges := make([]*model.ClassroomSectionRange, 0)
	for section := 1; section <= sectionsPerDay; section++ {
		if occupied[section] {
			continue
		}
		start := section
		for section < sectionsPerDay && !occupied[section+1] {
			section++
		}
		ranges = append(ranges, &model.ClassroomSectionRange{Start: int32(start), End: int32(section)})
	}
	return ranges
}

func location2Build(location string) string {
	runes := []rune(location)
	if strings.Contains(location, "公语") {
		return "东1"
	} else if strings.Contains(location, "中") {
		return "中楼"
	}
	return string(runes[2:4])
}
