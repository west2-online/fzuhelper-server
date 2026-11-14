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
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func BuildExamRoomInfo(rooms []*jwch.ExamRoomInfo) []*model.ExamRoomInfo {
	var res []*model.ExamRoomInfo
	for _, room := range rooms {
		res = append(res, &model.ExamRoomInfo{
			Name:     room.CourseName,
			Credit:   room.Credit,
			Teacher:  room.Teacher,
			Time:     room.Time,
			Date:     room.Date,
			Location: room.Location,
		})
	}
	return res
}

func BuildExamRoomInfoYjsy(rooms []*yjsy.ExamRoomInfo) []*model.ExamRoomInfo {
	var res []*model.ExamRoomInfo
	for _, room := range rooms {
		res = append(res, &model.ExamRoomInfo{
			Name:     room.CourseName,
			Credit:   room.Credit,
			Teacher:  room.Teacher,
			Time:     room.Time,
			Date:     room.Date,
			Location: room.Location,
		})
	}
	return res
}
