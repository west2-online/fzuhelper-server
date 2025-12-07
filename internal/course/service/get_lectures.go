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

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *CourseService) GetLectures(req *course.GetLecturesRequest, loginData *model.LoginData) ([]*model.Lecture, error) {
	isRefresh := false
	if req.IsRefresh != nil {
		isRefresh = *req.IsRefresh
	}
	return s.getLectures(isRefresh, loginData)
}

func (s *CourseService) getLectures(isRefresh bool, loginData *model.LoginData) ([]*model.Lecture, error) {
	lectKey := s.lectureKey(context.ExtractIDFromLoginData(loginData))

	// cache hits
	if !isRefresh && s.cache.IsKeyExist(s.ctx, lectKey) {
		lects, err := s.cache.Course.GetLecturesCache(s.ctx, lectKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetLectures: Get lectures failed: %w", err)
		}
		return pack.BuildLectures(lects), nil
	}

	return s.fetchAndCacheLectures(loginData)
}

func (s *CourseService) fetchAndCacheLectures(loginData *model.LoginData) ([]*model.Lecture, error) {
	var err error

	// fetch from jwch
	stu := jwch.NewStudent().WithLoginData(loginData.GetId(), utils.ParseCookies(loginData.GetCookies()))
	lects, err := stu.GetLectures()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetLectures: Get lectures from jwch failed: %w", err)
	}

	// async put lectures to cache
	lectKey := s.lectureKey(context.ExtractIDFromLoginData(loginData))
	s.taskQueue.Add(lectKey, taskqueue.QueueTask{Execute: func() error {
		return cache.SetSliceCache(s.cache, s.ctx, lectKey, lects, constants.CourseTermsKeyExpire, "Course.SetLecturesCache")
	}})

	return pack.BuildLectures(lects), err
}

func (s *CourseService) lectureKey(stuId string) string {
	return fmt.Sprintf("lecture:%s", stuId)
}
