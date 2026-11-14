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
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
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
			expectContains: `{"code":"10000","message":"ok","data":[]}`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/common/classroom/empty?date=2025-01-01&startTime=1&endTime=2&campus=qishan",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/common/classroom/empty?startTime=1&endTime=2&campus=qishan",
			expectContains: `{"code":"20001","message":"参数错误`,
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

func TestGetEmptyClassroomsCompatibility(t *testing.T) {
	for _, tc := range []struct {
		name    string
		query   string
		include bool
	}{
		{name: "old request"},
		{name: "explicit false", query: "&includeOccupied=false"},
		{name: "include occupied", query: "&includeOccupied=true", include: true},
	} {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetEmptyRoomRPC).To(func(_ context.Context, req *classroom.EmptyRoomRequest) ([]*model.Classroom, error) {
				assert.Equal(t, "2026-10-11", req.Date)
				assert.Equal(t, "旗山校区", req.Campus)
				assert.Equal(t, "3", req.StartTime)
				assert.Equal(t, "5", req.EndTime)
				assert.Equal(t, tc.include, req.GetIncludeOccupied())
				return []*model.Classroom{{
					Build: "东2", Location: "旗山东2-100", Capacity: "318", Type: "智慧教室报告厅型",
					FreeRanges:  []*model.ClassroomSectionRange{{Start: 5, End: 11}},
					Occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 4, StartWeek: 6, EndWeek: 6, Event: "讲座", Occupant: "先造", RawText: "讲座()(06-06)"}},
				}}, nil
			}).Build()
			router := route.NewEngine(&config.Options{})
			router.GET("/api/v1/common/classroom/empty", GetEmptyClassrooms)
			url := "/api/v1/common/classroom/empty?date=2026-10-11&campus=旗山校区&startTime=3&endTime=5" + tc.query
			response := ut.PerformRequest(router, consts.MethodGet, url, nil)
			require.Equal(t, http.StatusOK, response.Code)
			var legacy struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				Data    []struct{ Build, Location, Capacity, Type string }
			}
			require.NoError(t, json.Unmarshal(response.Result().Body(), &legacy))
			assert.Equal(t, "10000", legacy.Code)
			assert.Equal(t, "ok", legacy.Message)
			require.Len(t, legacy.Data, 1)
			assert.Equal(t, "东2", legacy.Data[0].Build)
			assert.Equal(t, "旗山东2-100", legacy.Data[0].Location)
			assert.Equal(t, "318", legacy.Data[0].Capacity)
			assert.Equal(t, "智慧教室报告厅型", legacy.Data[0].Type)
			body := string(response.Result().Body())
			assert.NotContains(t, body, `"is_free"`)
			if tc.include {
				assert.NotContains(t, body, `"reason"`)
				assert.NotContains(t, body, `"detail"`)
				assert.Contains(t, body, `"free_ranges":[{"start":5,"end":11}]`)
				assert.Contains(t, body, `"event":"讲座","occupant":"先造","raw_text":"讲座()(06-06)"`)
				assert.Contains(t, body, `"start_week":6`)
				assert.Contains(t, body, `"end_week":6`)
			} else {
				// 精确比较 JSON，保证未开启占用详情时没有新增字段。
				assert.JSONEq(t, `{"code":"10000","message":"ok","data":[{
					"build":"东2","location":"旗山东2-100","capacity":"318","type":"智慧教室报告厅型"
				}]}`, body)
			}
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
			expectContains: `{"code":"10000","message":"ok","data":[]}`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/classroom/exam?term=2024-2025-1",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/jwch/classroom/exam",
			expectContains: `{"code":"20001","message":"参数错误`,
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
