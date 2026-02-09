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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func TestGetExamRoomInfo(t *testing.T) {
	type testCase struct {
		name          string
		mockReturn    []*jwch.ExamRoomInfo
		mockError     error
		expectResult  []*model.ExamRoomInfo
		expectError   bool
		expectCached  bool
		cacheGetError error
	}

	tests := []testCase{
		{
			name: "GetExamRoomInfoWithoutCache",
			mockReturn: []*jwch.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectResult: []*model.ExamRoomInfo{
				{Location: "旗山东1"},
			},
		},
		{
			name: "GetExamRoomInfoFromCache",
			mockReturn: []*jwch.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectResult: []*model.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectCached: true,
		},
		{
			name:          "GetExamRoomInfoCacheGetError",
			expectError:   true,
			expectCached:  true,
			cacheGetError: assert.AnError,
		},
		{
			name:        "GetExamRoomInfoJwchError",
			mockError:   assert.AnError,
			expectError: true,
		},
		{
			name:         "GetExamRoomInfoEmptyResult",
			mockReturn:   []*jwch.ExamRoomInfo{},
			expectResult: []*model.ExamRoomInfo(nil),
		},
	}

	req := &classroom.ExamRoomInfoRequest{
		Term: "202401",
	}

	defer mockey.UnPatchAll()
	// 运行所有测试用例
	for _, tc := range tests {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				CacheClient: new(cache.Cache),
			}
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.expectCached).Build()
			mockey.Mock((*classroomCache.CacheClassroom).GetExamRoom).Return(tc.expectResult, tc.cacheGetError).Build()
			mockey.Mock((*jwch.Student).WithLoginData).Return(jwch.NewStudent()).Build()
			mockey.Mock((*jwch.Student).GetExamRoom).Return(tc.mockReturn, tc.mockError).Build()
			// mock login data
			loginData := &model.LoginData{
				Id:      "123456789",
				Cookies: "cookie1=value1;cookie2=value2",
			}

			ctx := customContext.WithLoginData(context.Background(), loginData)
			classroomService := NewClassroomService(ctx, mockClientSet)
			result, err := classroomService.GetExamRoomInfo(req, loginData)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestGetExamRoomInfoYjsy(t *testing.T) {
	type testCase struct {
		name          string
		mockReturn    []*yjsy.ExamRoomInfo
		mockError     error
		expectResult  []*model.ExamRoomInfo
		expectError   bool
		expectCached  bool
		cacheGetError error
	}

	tests := []testCase{
		{
			name: "GetExamRoomInfoYjsyWithoutCache",
			mockReturn: []*yjsy.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectResult: []*model.ExamRoomInfo{
				{Location: "旗山东1"},
			},
		},
		{
			name: "GetExamRoomInfoYjsyFromCache",
			mockReturn: []*yjsy.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectResult: []*model.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectCached: true,
		},
		{
			name:          "GetExamRoomInfoYjsyCacheGetError",
			expectError:   true,
			expectCached:  true,
			cacheGetError: assert.AnError,
		},
		{
			name:        "GetExamRoomInfoYjsyError",
			mockError:   assert.AnError,
			expectError: true,
		},
		{
			name:         "GetExamRoomInfoYjsyEmptyResult",
			mockReturn:   []*yjsy.ExamRoomInfo{},
			expectResult: []*model.ExamRoomInfo(nil),
		},
	}

	req := &classroom.ExamRoomInfoRequest{
		Term: "202401",
	}

	defer mockey.UnPatchAll()
	// 运行所有测试用例
	for _, tc := range tests {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.expectCached).Build()
			mockey.Mock((*classroomCache.CacheClassroom).GetExamRoom).Return(tc.expectResult, tc.cacheGetError).Build()
			mockey.Mock((*yjsy.Student).WithLoginData).Return(yjsy.NewStudent()).Build()
			mockey.Mock((*yjsy.Student).GetExamRoom).Return(tc.mockReturn, tc.mockError).Build()
			// mock login data
			loginData := &model.LoginData{
				Id:      "123456789",
				Cookies: "cookie1=value1;cookie2=value2",
			}

			ctx := customContext.WithLoginData(context.Background(), loginData)
			classroomService := NewClassroomService(ctx, mockClientSet)
			result, err := classroomService.GetExamRoomInfoYjsy(req, loginData)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
