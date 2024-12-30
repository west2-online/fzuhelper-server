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
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *UserService) GetLoginData(req *user.GetLoginDataRequest) (string, []string, error) {
	stu := jwch.NewStudent().WithUser(req.Id, req.Password)
	id, rawCookies, err := stu.GetIdentifierAndCookies()
	if err != nil {
		return "", nil, err
	}
	// 进行学生信息的存储
	go s.insertStudentInfo(req, stu)

	return id, utils.ParseCookiesToString(rawCookies), nil
}

func (s *UserService) insertStudentInfo(req *user.GetLoginDataRequest, stu *jwch.Student) {
	// 查询数据库是否存入此学生信息
	exist, _, err := s.db.User.GetStudentById(s.ctx, req.Id)
	if err != nil {
		logger.Errorf("service.insertStudentInfo: %v", err)
	}
	if exist {
		return
	}
	// 若未查询则将学生信息插入
	resp, err := stu.GetInfo()
	if err != nil {
		logger.Errorf("service.insertStudentInfo: jwch failed: %v", err)
	}
	grade, _ := strconv.Atoi(resp.Grade)
	userModel := &model.Student{
		StuId:     req.Id,
		Sex:       resp.Sex,
		College:   resp.College,
		Grade:     int64(grade),
		Major:     resp.Major,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
	err = s.db.User.CreateStudent(s.ctx, userModel)
	if err != nil {
		logger.Errorf("service.insertStudentInfo: %v", err)
	}
}
