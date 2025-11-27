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

package course

import (
	"context"
	"fmt"
	"strings"

	"github.com/west2-online/fzuhelper-server/internal/course/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

// CourseServiceImpl implements the last service interface defined in the IDL.
type CourseServiceImpl struct {
	ClientSet *base.ClientSet
	taskQueue taskqueue.TaskQueue
}

func NewCourseService(clientSet *base.ClientSet, taskQueue taskqueue.TaskQueue) *CourseServiceImpl {
	return &CourseServiceImpl{
		ClientSet: clientSet,
		taskQueue: taskQueue,
	}
}

// GetCourseList implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetCourseList(ctx context.Context, req *course.CourseListRequest) (resp *course.CourseListResponse, err error) {
	resp = course.NewCourseListResponse()
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Academic.GetScores: Get login data fail %w", err)
	}
	if strings.HasPrefix(loginData.Id[:5], "00000") {
		res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetCourseListYjsy(req, loginData)
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return resp, nil
		}
		resp.Base = base.BuildSuccessResp()
		resp.Data = res
		return resp, nil
	} else {
		// 检查学期是否合法的逻辑在 service 里面实现了，这里不需要再检查
		// 原因：GetSemesterCourses() 要用到 jwch 里面的 GetTerms() 函数返回的 ViewState 和 EventValidation 参数，顺便检查可以减少请求次数
		res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetCourseList(req, loginData)
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return resp, nil
		}
		resp.Base = base.BuildSuccessResp()
		resp.Data = res
		return resp, nil
	}
}

func (s *CourseServiceImpl) GetTermList(ctx context.Context, req *course.TermListRequest) (resp *course.TermListResponse, err error) {
	resp = course.NewTermListResponse()
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Academic.GetScores: Get login data fail %w", err)
	}
	if strings.HasPrefix(loginData.Id[:5], "00000") {
		res, err := service.NewCourseService(ctx, s.ClientSet, nil).GetTermsListYjsy(loginData)
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return resp, nil
		}
		resp.Base = base.BuildSuccessResp()
		resp.Data = res
		return resp, nil
	} else {
		res, err := service.NewCourseService(ctx, s.ClientSet, nil).GetTermsList(loginData)
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return resp, nil
		}
		resp.Base = base.BuildSuccessResp()
		resp.Data = res
		return resp, nil
	}
}

func (s *CourseServiceImpl) GetCalendar(ctx context.Context, req *course.GetCalendarRequest) (resp *course.GetCalendarResponse, err error) {
	resp = course.NewGetCalendarResponse()

	resp.Ics, err = service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetCalendar(req.StuId)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()

	return resp, nil
}

func (s *CourseServiceImpl) GetLocateDate(ctx context.Context, _ *course.GetLocateDateRequest) (resp *course.GetLocateDateResponse, err error) {
	resp = course.NewGetLocateDateResponse()

	res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetLocateDate()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.LocateDate = res
	return resp, nil
}

func (s *CourseServiceImpl) GetFriendCourse(ctx context.Context, req *course.GetFriendCourseRequest) (
	resp *course.GetFriendCourseResponse, err error,
) {
	resp = new(course.GetFriendCourseResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Course.GetFriendCourse: Get login data fail %w", err)
	}
	res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetFriendCourse(req, loginData)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = res
	return resp, nil
}
