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
	"errors"
	"github.com/west2-online/fzuhelper-server/pkg/errno"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	jwchErrno "github.com/west2-online/jwch/errno"
)

func (s *ClassroomService) GetExamRoomInfo(req *classroom.ExamRoomInfoRequest) ([]*jwch.ExamRoomInfo, error) {
	stu := jwch.NewStudent().WithLoginData(req.LoginData.Id, utils.ParseCookies(req.LoginData.Cookies))
	rooms, err := stu.GetExamRoom(jwch.ExamRoomReq{Term: req.Term})
	if errors.Is(err, &jwchErrno.SessionExpiredError) {
		return nil, errno.Errorf(errno.AuthExpiredCode, "Classroom.GetExamRoomInfo: cookies expired")
	}
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "Classroom.GetExamRoomInfo: jwch error: %v", err)
	}
	return rooms, nil
}
