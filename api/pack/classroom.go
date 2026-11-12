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
	classroomModel "github.com/west2-online/fzuhelper-server/api/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildClassroom(res *model.Classroom, includeOccupied bool) *classroomModel.Classroom {
	room := &classroomModel.Classroom{
		Build:    res.Build,
		Location: res.Location,
		Capacity: res.Capacity,
		Type:     res.Type,
	}
	// 普通查询仅返回旧版四个字段；占用详情由客户端显式开启。
	if !includeOccupied {
		return room
	}
	if res.FreeRanges != nil {
		room.FreeRanges = make([]*classroomModel.ClassroomSectionRange, 0, len(res.FreeRanges))
		for _, freeRange := range res.FreeRanges {
			room.FreeRanges = append(room.FreeRanges, &classroomModel.ClassroomSectionRange{Start: freeRange.Start, End: freeRange.End})
		}
	}
	if res.Occupancies != nil {
		room.Occupancies = make([]*classroomModel.ClassroomOccupancy, 0, len(res.Occupancies))
		for _, occupancy := range res.Occupancies {
			room.Occupancies = append(room.Occupancies, &classroomModel.ClassroomOccupancy{
				Start: occupancy.Start, End: occupancy.End, Event: occupancy.Event, Occupant: occupancy.Occupant, RawText: occupancy.RawText,
				StartWeek: occupancy.StartWeek, EndWeek: occupancy.EndWeek,
			})
		}
	}
	return room
}

func BuildClassroomList(res []*model.Classroom, includeOccupied bool) []*classroomModel.Classroom {
	list := make([]*classroomModel.Classroom, 0, len(res))
	for _, v := range res {
		list = append(list, BuildClassroom(v, includeOccupied))
	}
	return list
}

func BuildExamRoomInfo(rooms []*model.ExamRoomInfo) []*classroomModel.ExamRoomInfo {
	list := make([]*classroomModel.ExamRoomInfo, 0, len(rooms))
	for _, room := range rooms {
		list = append(list, &classroomModel.ExamRoomInfo{
			Name:     room.Name,
			Credit:   room.Credit,
			Teacher:  room.Teacher,
			Date:     room.Date,
			Time:     room.Time,
			Location: room.Location,
		})
	}
	return list
}
