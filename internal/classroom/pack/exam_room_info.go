package pack

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
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
