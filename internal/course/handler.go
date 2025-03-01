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
	"sync"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/internal/course/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

// CourseServiceImpl implements the last service interface defined in the IDL.
type CourseServiceImpl struct {
	ClientSet     *base.ClientSet
	taskQueue     taskqueue.TaskQueue
	courseLockMap *sync.Map
}

func NewCourseService(clientSet *base.ClientSet, taskQueue taskqueue.TaskQueue, courseLockMap *sync.Map) *CourseServiceImpl {
	return &CourseServiceImpl{
		ClientSet:     clientSet,
		taskQueue:     taskQueue,
		courseLockMap: courseLockMap,
	}
}

// GetCourseList implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetCourseList(ctx context.Context, req *course.CourseListRequest) (resp *course.CourseListResponse, err error) {
	resp = course.NewCourseListResponse()
	// 检查学期是否合法的逻辑在 service 里面实现了，这里不需要再检查
	// 原因：GetSemesterCourses() 要用到 jwch 里面的 GetTerms() 函数返回的 ViewState 和 EventValidation 参数，顺便检查可以减少请求次数
	res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue, s.courseLockMap).GetCourseList(req)
	if err != nil {
		logger.Infof("Course.GetCourseList: GetCourseList failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = pack.BuildCourse(res)
	return resp, nil
}

func (s *CourseServiceImpl) GetTermList(ctx context.Context, req *course.TermListRequest) (resp *course.TermListResponse, err error) {
	resp = course.NewTermListResponse()

	res, err := service.NewCourseService(ctx, s.ClientSet, nil, nil).GetTermsList(req)
	if err != nil {
		logger.Infof("Course.GetTermList: GetTermList failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = res
	return resp, nil
}

func (s *CourseServiceImpl) GetCalendar(ctx context.Context, req *course.GetCalendarRequest) (resp *course.GetCalendarResponse, err error) {
	resp = course.NewGetCalendarResponse()

	res, err := service.NewCourseService(ctx, s.ClientSet, nil, nil).GetCalendar(req)
	if err != nil {
		logger.Infof("Course.GetCalendar: GetCalendar failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Content = res
	return resp, nil
}

func (s *CourseServiceImpl) GetLocateDate(ctx context.Context, _ *course.GetLocateDateRequest) (resp *course.GetLocateDateResponse, err error) {
	resp = course.NewGetLocateDateResponse()

	res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue, nil).GetLocateDate()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.LocateDate = res
	return resp, nil
}
