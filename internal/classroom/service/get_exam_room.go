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

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *ClassroomService) GetExamRoomInfo(req *classroom.ExamRoomInfoRequest) ([]*jwch.ExamRoomInfo, error) {
	// login data 从 ctx 中获取

	loginData, err := context.GetLoginData(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetExamRoomInfo: Get login data fail %w", err)
	}
	key := fmt.Sprintf("exam:user:%s:term:%s", loginData.GetId(), req.GetTerm())

	if s.cache.IsKeyExist(s.ctx, key) {
		examRooms, err := s.cache.Classroom.GetExamRoom(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetExamRoomInfo: Get exam room fail %w", err)
		}
		return examRooms, nil
	}
	stu := jwch.NewStudent().WithLoginData(loginData.Id, utils.ParseCookies(loginData.Cookies))
	examRooms, err := stu.GetExamRoom(jwch.ExamRoomReq{Term: req.Term})
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetExamRoomInfo: Get exam room info fail %w", err)
	}
	go s.cache.Classroom.SetExamRoom(s.ctx, key, examRooms)
	return examRooms, nil
}
