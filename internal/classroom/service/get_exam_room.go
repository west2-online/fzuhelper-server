package service

import (
	"fmt"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *ClassroomService) GetExamRoomInfo(req *classroom.ExamRoomInfoRequest) ([]*jwch.ExamRoomInfo, error) {
	stu := jwch.NewStudent().WithLoginData(req.LoginData.Id, utils.ParseCookies(req.LoginData.Cookies))
	rooms, err := stu.GetExamRoom(jwch.ExamRoomReq{Term: req.Term})
	if err != nil {
		return nil, fmt.Errorf("service.GetExamRoomInfo: Get exam room info fail %w", err)
	}
	return rooms, nil
}
