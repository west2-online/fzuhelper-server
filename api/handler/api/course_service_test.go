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

package api

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/mw"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestGetCourseList(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       []*model.Course
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/course/list?term=202401",
			mockResp:       []*model.Course{},
			expectContains: `{"code":"10000","message":"ok","data":[]}`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/course/list?term=202401",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/jwch/course/list",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/course/list", GetCourseList)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetCourseListRPC).To(func(ctx context.Context, req *course.CourseListRequest) ([]*model.Course, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetTermList(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *course.TermListResponse
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/term/list",
			mockResp:       &course.TermListResponse{Data: []string{"202401"}},
			expectContains: `{"code":"10000","message":"ok","data":["202401"]}`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/term/list",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/term/list", GetTermList)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetCourseTermsListRPC).To(func(ctx context.Context, req *course.TermListRequest) (*course.TermListResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetLocateDate(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *model.LocateDate
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/course/date",
			mockResp:       &model.LocateDate{},
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/course/date",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/course/date", GetLocateDate)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetLocateDateRPC).To(func(ctx context.Context, req *course.GetLocateDateRequest) (*model.LocateDate, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestSubscribeCalendar(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       []byte
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/course/calendar/subscribe?stu=1",
			mockResp:       []byte("BEGIN:VCALENDAR"),
			expectContains: "BEGIN:VCALENDAR",
		},
		{
			name:           "rpc error",
			url:            "/api/v1/course/calendar/subscribe?stu=1",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "missing stu id",
			url:            "/api/v1/course/calendar/subscribe",
			expectContains: `{"code":"20001","message":"参数错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.Use(func(ctx context.Context, c *app.RequestContext) {
		// 从 query 参数读取并设置 stu_id 到 context
		stuIDStr := c.Query("stu")
		if stuIDStr != "" {
			c.Set("stu_id", stuIDStr)
		}
	})
	router.GET("/api/v1/course/calendar/subscribe", SubscribeCalendar)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetCalendarRPC).To(func(ctx context.Context, req *course.GetCalendarRequest) ([]byte, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetCalendar(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockToken      string
		mockLoginData  *model.LoginData
		mockLoginErr   error
		mockTokenErr   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/course/calendar/token",
			mockToken:      "token123",
			mockLoginData:  &model.LoginData{Id: "202400001"},
			expectContains: `{"code":"10000","message":"Success","data":"token123"}`,
		},
		{
			name:           "login error",
			url:            "/api/v1/jwch/course/calendar/token",
			mockLoginData:  &model.LoginData{Id: "202400001"},
			mockLoginErr:   errno.AuthInvalid,
			expectContains: `{"code":"30001","message":"鉴权失败, [30002] 鉴权无效"}`,
		},
		{
			name:           "create token error",
			url:            "/api/v1/jwch/course/calendar/token",
			mockLoginData:  &model.LoginData{Id: "202400001"},
			mockTokenErr:   errno.AuthError,
			expectContains: `{"code":"30001","message":"鉴权失败, [30001] 鉴权失败"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/course/calendar/token", GetCalendar)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(metainfoContext.GetLoginData).To(func(ctx context.Context) (*model.LoginData, error) {
				return tc.mockLoginData, tc.mockLoginErr
			}).Build()
			mockey.Mock(mw.CreateToken).To(func(tokenType int64, stuID string) (string, error) {
				return tc.mockToken, tc.mockTokenErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetFriendCourse(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       []*model.Course
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/course/friend?term=202401&student_id=102300001",
			mockResp:       []*model.Course{},
			expectContains: `{"code":"10000","message":"ok","data":[]}`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/course/friend?term=202401&student_id=102300001",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/course/friend?term=202401",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/course/friend", GetFriendCourse)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetFriendCourseRPC).To(func(ctx context.Context, req *course.GetFriendCourseRequest) ([]*model.Course, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
