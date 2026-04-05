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

	"github.com/west2-online/fzuhelper-server/internal/classroom/pack"
	"github.com/west2-online/fzuhelper-server/internal/classroom/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
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
	resp = new(classroom.EmptyRoomResponse)
	res, err := service.NewClassroomService(ctx, s.ClientSet).GetEmptyRoom(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Rooms = pack.BuildClassRooms(res, req.Campus)
	return resp, nil
}

func (s *ClassroomServiceImpl) GetExamRoomInfo(ctx context.Context, req *classroom.ExamRoomInfoRequest) (resp *classroom.ExamRoomInfoResponse, err error) {
	resp = new(classroom.ExamRoomInfoResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.ErrNoWithPreMessage(err, "Classroom.GetExamRoomInfo: Get login data failed"))
		return resp, nil
	}
	if utils.IsGraduate(loginData.Id) {
		rooms, err := service.NewClassroomService(ctx, s.ClientSet).GetExamRoomInfoYjsy(req, loginData)
		resp.Base = base.BuildBaseResp(err)
		if err != nil {
			return resp, nil
		}
		resp.Rooms = rooms
		return resp, nil
	} else {
		rooms, err := service.NewClassroomService(ctx, s.ClientSet).GetExamRoomInfo(req, loginData)
		resp.Base = base.BuildBaseResp(err)
		if err != nil {
			return resp, nil
		}
		resp.Rooms = rooms
		return resp, nil
	}
}
