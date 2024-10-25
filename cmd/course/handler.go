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

package main

import (
	"context"

	"github.com/west2-online/fzuhelper-server/cmd/course/pack"
	"github.com/west2-online/fzuhelper-server/cmd/course/service"
	course "github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// CourseServiceImpl implements the last service interface defined in the IDL.
type CourseServiceImpl struct{}

// GetCourseList implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetCourseList(ctx context.Context, req *course.CourseListRequest) (resp *course.CourseListResponse, err error) {
	resp = course.NewCourseListResponse()

	// 检查学期是否合法的逻辑在 service 里面实现了，这里不需要再检查
	// 原因：GetSemesterCourses() 要用到 jwch 里面的 GetTerms() 函数返回的 ViewState 和 EventValidation 参数，顺便检查可以减少请求次数

	res, err := service.NewCourseService(ctx).GetCourseList(req)
	if err != nil {
		logger.Infof("Course.GetCourseList: GetCourseList failed, err: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Data = pack.BuildCourse(res)

	// logger.Infof("Course.GetCourseList: GetCourseList success, data: %v", resp.Data)

	return resp, nil
}
