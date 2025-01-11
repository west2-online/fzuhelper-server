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
	"strconv"

	db "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/jwch"
)

func (s *UserService) GetUserInfo(stuId string) (*db.Student, error) {
	// 查询数据库是否存入此学生信息
	exist, stuInfo, err := s.db.User.GetStudentById(s.ctx, stuId)
	if err != nil {
		return nil, fmt.Errorf("service.GetUserInfo: %w", err)
	}
	if exist {
		return stuInfo, nil
	}

	// 若未查询则将学生信息插入
	stu := jwch.NewStudent().WithLoginData(s.Identifier, s.cookies)
	resp, err := stu.GetInfo()
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "service.GetUserInfo: jwch failed: %v", err)
	}
	grade, _ := strconv.Atoi(resp.Grade)
	userModel := &db.Student{
		StuId:    stuId,
		Sex:      resp.Sex,
		Birthday: resp.Birthday,
		College:  resp.College,
		Grade:    int64(grade),
		Major:    resp.Major,
	}
	err = s.db.User.CreateStudent(s.ctx, userModel)
	if err != nil {
		return nil, fmt.Errorf("service.GetUserInfo: %w", err)
	}
	return userModel, nil
}
