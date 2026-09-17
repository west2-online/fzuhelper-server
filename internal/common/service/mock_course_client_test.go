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
	"context"
	"errors"

	"github.com/cloudwego/kitex/client/callopt"

	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	rpcmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// mockCourseClient implements courseservice.Client for testing ProcessAutoAdjustCourseNotice.
type mockCourseClient struct {
	createResp *course.CreateAdjustCourseResponse
	createErr  error
	// 记录调用参数，便于断言透传给 course 服务的内容
	createCalled bool
	createReq    *course.CreateAdjustCourseRequest
}

func (m *mockCourseClient) CreateAdjustCourse(
	ctx context.Context,
	req *course.CreateAdjustCourseRequest,
	opts ...callopt.Option,
) (*course.CreateAdjustCourseResponse, error) {
	m.createCalled = true
	m.createReq = req
	if m.createResp == nil && m.createErr == nil {
		// 兜底：未显式配置返回值时给出一个成功响应，避免用例漏配时在 resp.Base 处 panic
		return &course.CreateAdjustCourseResponse{Base: &rpcmodel.BaseResp{Code: errno.SuccessCode}}, nil
	}
	return m.createResp, m.createErr
}

// unused methods
func (m *mockCourseClient) GetCourseList(
	context.Context,
	*course.CourseListRequest,
	...callopt.Option,
) (*course.CourseListResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) GetTermList(
	context.Context,
	*course.TermListRequest,
	...callopt.Option,
) (*course.TermListResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) GetCalendar(
	context.Context,
	*course.GetCalendarRequest,
	...callopt.Option,
) (*course.GetCalendarResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) GetLocateDate(
	context.Context,
	*course.GetLocateDateRequest,
	...callopt.Option,
) (*course.GetLocateDateResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) GetFriendCourse(
	context.Context,
	*course.GetFriendCourseRequest,
	...callopt.Option,
) (*course.GetFriendCourseResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) GetAutoAdjustCourseList(
	context.Context,
	*course.GetAutoAdjustCourseListRequest,
	...callopt.Option,
) (*course.GetAutoAdjustCourseListResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) UpdateAdjustCourse(
	context.Context,
	*course.UpdateAdjustCourseRequest,
	...callopt.Option,
) (*course.UpdateAdjustCourseResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) UpsertCustomCourse(
	context.Context,
	*course.UpsertCustomCourseRequest,
	...callopt.Option,
) (*course.UpsertCustomCourseResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCourseClient) DeleteCustomCourse(
	context.Context,
	*course.DeleteCustomCourseRequest,
	...callopt.Option,
) (*course.DeleteCustomCourseResponse, error) {
	return nil, errors.New("not implemented")
}
