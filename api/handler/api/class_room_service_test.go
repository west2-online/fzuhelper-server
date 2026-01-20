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
	"errors"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestGetEmptyClassrooms(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/common/classroom/empty?date=2025-01-01&startTime=1&endTime=2&campus=qishan",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/common/classroom/empty?date=2025-01-01&startTime=1&endTime=2&campus=qishan",
			mockRPCError:   errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/common/classroom/empty?startTime=1&endTime=2&campus=qishan",
			mockRPCError:   nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/common/classroom/empty", GetEmptyClassrooms)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetEmptyRoomRPC).To(func(ctx context.Context, req *classroom.EmptyRoomRequest) ([]*model.Classroom, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return []*model.Classroom{}, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetExamRoomInfo(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/classroom/exam?term=2024-2025-1",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/classroom/exam?term=2024-2025-1",
			mockRPCError:   errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/jwch/classroom/exam",
			mockRPCError:   nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/classroom/exam", GetExamRoomInfo)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetExamRoomInfoRPC).To(func(ctx context.Context, req *classroom.ExamRoomInfoRequest) ([]*model.ExamRoomInfo, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return []*model.ExamRoomInfo{}, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
