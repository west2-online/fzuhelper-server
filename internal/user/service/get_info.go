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
	"time"

	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	db "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (s *UserService) GetUserInfo(loginData *loginmodel.LoginData) (*db.Student, error) {
	stuId := context.ExtractIDFromLoginData(loginData)
	// 查询cache
	existCache := s.cache.IsKeyExist(s.ctx, stuId)
	if existCache {
		stuInfo, err := s.cache.User.GetStuInfoCache(s.ctx, stuId)
		if err != nil {
			return nil, errno.Errorf(errno.InternalRedisErrorCode, "User.GetUserInfo: %v", err)
		}
		return stuInfo, nil
	}

	// 未命中cache，查询数据库是否存入此学生信息
	exist, stuInfo, err := s.db.User.GetStudentById(s.ctx, stuId)
	IsUpdate := false
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "User.GetUserInfo: %v", err)
	}
	if exist {
		if stuInfo.UpdatedAt.Add(constants.StuInfoExpireTime).After(time.Now()) {
			s.taskQueue.Add(fmt.Sprintf("setStuInfoCache:%s", stuId), taskqueue.QueueTask{Execute: func() error {
				return s.cache.User.SetStuInfoCache(s.ctx, stuId, stuInfo)
			}})
			return stuInfo, nil
		}
		IsUpdate = true
	}

	// 将学生信息插入/更新
	stu := jwch.NewStudent().WithLoginData(loginData.Id, utils.ParseCookies(loginData.GetCookies()))
	resp, err := stu.GetInfo()
	if err = base.HandleJwchError(err); err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "User.GetUserInfo: jwch failed: %v", err)
	}
	grade, _ := strconv.Atoi(resp.Grade)
	userModel := &db.Student{
		StuId:    stuId,
		Name:     resp.Name,
		Sex:      resp.Sex,
		Birthday: resp.Birthday,
		College:  resp.College,
		Grade:    int64(grade),
		Major:    resp.Major,
	}
	if IsUpdate {
		err = s.db.User.UpdateStudent(s.ctx, userModel)
	} else {
		err = s.db.User.CreateStudent(s.ctx, userModel)
	}
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "User.GetUserInfo: %v", err)
	}

	// 存入cache
	s.taskQueue.Add(fmt.Sprintf("setStuInfoCache:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		return s.cache.User.SetStuInfoCache(s.ctx, stuId, userModel)
	}})

	return userModel, nil
}

func (s *UserService) GetUserInfoYjsy(loginData *loginmodel.LoginData) (*db.Student, error) {
	stuId := context.ExtractIDFromLoginData(loginData)
	// 查询cache
	existCache := s.cache.IsKeyExist(s.ctx, stuId)
	if existCache {
		stuInfo, err := s.cache.User.GetStuInfoCache(s.ctx, stuId)
		if err != nil {
			return nil, errno.Errorf(errno.InternalRedisErrorCode, "User.GetUserInfoYjsy: %v", err)
		}
		return stuInfo, nil
	}

	// 未命中cache，查询数据库是否存入此学生信息
	exist, stuInfo, err := s.db.User.GetStudentById(s.ctx, stuId)
	IsUpdate := false
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "User.GetUserInfoYjsy: %v", err)
	}
	if exist {
		if stuInfo.UpdatedAt.Add(constants.StuInfoExpireTime).After(time.Now()) {
			s.taskQueue.Add(fmt.Sprintf("setStuInfoCache:%s", stuId), taskqueue.QueueTask{Execute: func() error {
				return s.cache.User.SetStuInfoCache(s.ctx, stuId, stuInfo)
			}})
			return stuInfo, nil
		}
		IsUpdate = true
	}

	// 将学生信息插入/更新
	stu := yjsy.NewStudent().WithLoginData(utils.ParseCookies(loginData.GetCookies()))
	resp, err := stu.GetStudentInfo()
	if err = base.HandleYjsyError(err); err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "User.GetUserInfoYjsy: yjsy failed: %v", err)
	}
	grade, _ := strconv.Atoi(resp.Grade)
	userModel := &db.Student{
		StuId:    stuId,
		Name:     resp.Name,
		Sex:      resp.Sex,
		Birthday: resp.Birthday,
		College:  resp.College,
		Grade:    int64(grade),
		Major:    resp.Major,
	}
	if IsUpdate {
		err = s.db.User.UpdateStudent(s.ctx, userModel)
	} else {
		err = s.db.User.CreateStudent(s.ctx, userModel)
	}
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "User.GetUserInfoYjsy: %v", err)
	}

	// 存入cache
	s.taskQueue.Add(fmt.Sprintf("setStuInfoCache:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		return s.cache.User.SetStuInfoCache(s.ctx, stuId, userModel)
	}})

	return userModel, nil
}
