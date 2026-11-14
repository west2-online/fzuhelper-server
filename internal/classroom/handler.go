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

	"github.com/west2-online/fzuhelper-server/internal/classroom/refresh"
	"github.com/west2-online/fzuhelper-server/internal/classroom/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/singleflight"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const (
	HoursInADay = 24
	MinDateDiff = 0
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
	// 今天至下周日。
	requestDate, err := time.ParseInLocation(time.DateOnly, req.Date, constants.ChinaTZ)
	if err != nil {
		logger.WithCtx(ctx).Errorf("Classroom.GetEmptyRoom: date format error, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	now := time.Now().In(constants.ChinaTZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, constants.ChinaTZ)
	dateDiff := requestDate.Sub(today).Hours() / HoursInADay
	maxDateDiff := len(refresh.Dates(now)) - 1
	if dateDiff < MinDateDiff || dateDiff > float64(maxDateDiff) {
		err = fmt.Errorf("date out of range, date: %v", req.Date)
		logger.WithCtx(ctx).Infof("Classroom.GetEmptyRoom: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}

	res, err := service.NewClassroomService(ctx, s.ClientSet).GetEmptyRoom(req)
	if err != nil {
		logger.WithCtx(ctx).Infof("Classroom.GetEmptyRoom: GetEmptyRoom failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Rooms = res
	// logger.WithCtx(ctx).Info("Classroom.GetEmptyRoom: GetEmptyRoom success")
	return resp, nil
}

func (s *ClassroomServiceImpl) GetExamRoomInfo(ctx context.Context, req *classroom.ExamRoomInfoRequest) (resp *classroom.ExamRoomInfoResponse, err error) {
	resp = classroom.NewExamRoomInfoResponse()
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Classroom.GetExamRoomInfo: Get login data fail %w", err)
	}
	stuId := metainfoContext.ExtractIDFromLoginData(loginData)
	isGraduate := utils.IsGraduate(loginData.Id)
	key := singleflight.Key(constants.SingleflightExamRoomsPrefix, stuId, req.GetTerm(), isGraduate)

	rooms, err := singleflight.Do(key, func() ([]*model.ExamRoomInfo, error) {
		svc := service.NewClassroomService(ctx, s.ClientSet)
		if isGraduate {
			return svc.GetExamRoomInfoYjsy(req, loginData)
		} else {
			return svc.GetExamRoomInfo(req, loginData)
		}
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Rooms = rooms
	return resp, nil
}
