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
	"fmt"

	"github.com/west2-online/fzuhelper-server/internal/classroom/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (s *ClassroomService) GetExamRoomInfo(req *classroom.ExamRoomInfoRequest, loginData *model.LoginData) ([]*model.ExamRoomInfo, error) {
	key := fmt.Sprintf("exam:user:%s:term:%s", context.ExtractIDFromLoginData(loginData), req.GetTerm())

	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		examRooms, err := s.cache.Classroom.GetExamRoom(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetExamRoomInfo: Get exam room fail %w", err)
		}
		return examRooms, nil
	}

	stu := jwch.NewStudent().WithLoginData(loginData.Id, utils.ParseCookies(loginData.Cookies))
	rawRooms, err := stu.GetExamRoom(jwch.ExamRoomReq{Term: req.Term})
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetExamRoomInfo: Get exam room info fail %w", err)
	}
	modelRooms := pack.BuildExamRoomInfo(rawRooms)
	go s.cache.Classroom.SetExamRoom(s.ctx, key, modelRooms)
	return modelRooms, nil
}

func (s *ClassroomService) GetExamRoomInfoYjsy(req *classroom.ExamRoomInfoRequest, loginData *model.LoginData) ([]*model.ExamRoomInfo, error) {
	key := fmt.Sprintf("exam:user:%s:term:%s", context.ExtractIDFromLoginData(loginData), req.GetTerm())

	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		examRooms, err := s.cache.Classroom.GetExamRoom(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetExamRoomInfo: Get exam room fail %w", err)
		}
		return examRooms, nil
	}

	stu := yjsy.NewStudent().WithLoginData(utils.ParseCookies(loginData.Cookies))
	rawRooms, err := stu.GetExamRoom(yjsy.ExamRoomReq{Term: req.Term})
	if err = base.HandleYjsyError(err); err != nil {
		return nil, fmt.Errorf("service.GetExamRoomInfo: Get exam room info fail %w", err)
	}
	modelRooms := pack.BuildExamRoomInfoYjsy(rawRooms)
	go s.cache.Classroom.SetExamRoom(s.ctx, key, modelRooms)
	return modelRooms, nil
}
