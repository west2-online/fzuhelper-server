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

func BuildClassroom(str string, campus string) *model.Classroom {
	// 旗山东1-103 0(0) 机房
	// 晋江A102 150(75) 多媒体
	// 铜盘A109 120(60) 多媒体
	// 泉港教-110 63(40) 多媒体
	// 怡山北301 92(0) 多媒体
	// 鼓浪屿多媒体1 0(0) 多媒体
	// 集美1-311 287(287) 多媒体
	parts := strings.Fields(str)

	location := parts[0]
	capacityWithParentheses := parts[1]
	roomType := strings.Join(parts[2:], " ")

	// Remove the parentheses from capacity
	capacity := strings.Split(capacityWithParentheses, "(")[0]

	// 只有旗山校区拥有build字段，其余build返回campus
	// 接下来通过location来判断build
	// TODO: 可能有些笨拙，不过没有什么好办法----
	if strings.Contains(campus, "旗山") {
		return &model.Classroom{
			Build:    location2Build(location),
			Location: location,
			Capacity: strings.TrimSpace(capacity),
			Type:     roomType,
		}
	} else {
		return &model.Classroom{
			Build:    campus,
			Location: location,
			Capacity: strings.TrimSpace(capacity),
			Type:     roomType,
		}
	}
}

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

func BuildClassRooms(strs []string, campus string) (res []*model.Classroom) {
	for _, str := range strs {
		res = append(res, BuildClassroom(str, campus))
	}
	return res
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
