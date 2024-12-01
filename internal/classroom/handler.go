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

package classroom

import (
	"context"
	"fmt"
	"time"

	"github.com/west2-online/fzuhelper-server/internal/classroom/pack"
	"github.com/west2-online/fzuhelper-server/internal/classroom/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const (
	HoursInADay = 24
	MinDateDiff = 0
	MaxDateDiff = 30
)

// ClassroomServiceImpl implements the last service interface defined in the IDL.
type ClassroomServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewClassroomService(clientSet *base.ClientSet) *ClassroomServiceImpl {
	return &ClassroomServiceImpl{
		ClientSet: clientSet,
	}
}

// GetEmptyRoom implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetEmptyRoom(ctx context.Context, req *classroom.EmptyRoomRequest) (resp *classroom.EmptyRoomResponse, err error) {
	resp = classroom.NewEmptyRoomResponse()
	// 判断req.date只能从今天开始的一个月内，在当前日期前或超过 30 天则报错
	// 首先判断date的格式是否符合要求
	requestDate, err := utils.TimeParse(req.Date)
	if err != nil {
		logger.Errorf("Classroom.GetEmptyRoom: date format error, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	now := time.Now().Truncate(constants.OneDay)
	requestDate = requestDate.Truncate(constants.OneDay)
	dateDiff := requestDate.Sub(now).Hours() / HoursInADay
	if dateDiff < MinDateDiff || dateDiff > MaxDateDiff {
		err = fmt.Errorf("date out of range, date: %v", req.Date)
		logger.Infof("Classroom.GetEmptyRoom: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}

	res, err := service.NewClassroomService(ctx, s.ClientSet).GetEmptyRoom(req)
	if err != nil {
		logger.Infof("Classroom.GetEmptyRoom: GetEmptyRoom failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Rooms = pack.BuildClassRooms(res, req.Campus)
	// logger.Info("Classroom.GetEmptyRoom: GetEmptyRoom success")
	return resp, nil
}

func (s *ClassroomServiceImpl) GetExamRoomInfo(ctx context.Context, req *classroom.ExamRoomInfoRequest) (resp *classroom.ExamRoomInfoResponse, err error) {
	resp = classroom.NewExamRoomInfoResponse()
	rooms, err := service.NewClassroomService(ctx, s.ClientSet).GetExamRoomInfo(req)
	if err != nil {
		logger.Infof("Classroom.GetExamRoomInfo: Get exam room info fail %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Rooms = pack.BuildExamRoomInfo(rooms)
	return resp, nil
}
