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
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	courseCache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestCourseService_getLectures(t *testing.T) {
	mockLectures := []*jwch.Lecture{
		{
			Category:         "jxjt",
			IssueNumber:      123,
			Title:            "Cooking",
			Speaker:          "Walter White",
			Timestamp:        1234567890123,
			Location:         "New Mexico",
			AttendanceStatus: "123",
		},
	}

	mockResult := []*model.Lecture{
		{
			Category:         "jxjt",
			IssueNumber:      123,
			Title:            "Cooking",
			Speaker:          "Walter White",
			TimeStamp:        1234567890123,
			Location:         "New Mexico",
			AttendanceStatus: "123",
		},
	}

	const (
		GetLecturesFailedMsg = "Get lectures from jwch failed"
	)

	type testCase struct {
		name               string
		expectedResult     []*model.Lecture
		expectingError     bool
		expectedErrorMsg   string
		mockLecturesError  error
		mockLecturesReturn []*jwch.Lecture
		cacheExist         bool
		cacheGetError      error // 获取缓存时的报错（如果有）
	}
	testCases := []testCase{
		{
			name:               "getLecturesSuccess", // 成功
			mockLecturesReturn: mockLectures,
			mockLecturesError:  nil,
			expectedResult:     mockResult,
			expectingError:     false,
		},
		{
			name:               "getLecturesGetLecturesFailed", // 从 jwch 拉取失败
			mockLecturesReturn: nil,
			mockLecturesError:  errors.New(GetLecturesFailedMsg),
			expectedResult:     nil,
			expectingError:     true,
			expectedErrorMsg:   GetLecturesFailedMsg,
		},
		{
			name:           "cacheHits", // 缓存命中
			cacheExist:     true,
			cacheGetError:  nil,
			expectedResult: mockResult,
		},
	}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*jwch.Student).GetLectures).Return(tc.mockLecturesReturn, tc.mockLecturesError).Build()
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			if tc.cacheExist {
				mockey.Mock((*courseCache.CacheCourse).GetLecturesCache).To(
					func(ctx context.Context, key string) ([]*jwch.Lecture, error) {
						if tc.cacheGetError != nil {
							return nil, tc.cacheGetError
						}
						return mockLectures, nil
					},
				).Build()
			} else {
				mockey.Mock((*courseCache.CacheCourse).GetLecturesCache).To(
					func(ctx context.Context, key string) ([]*jwch.Lecture, error) {
						return nil, fmt.Errorf("should not be called since in current test case cache doesn't exist")
					},
				).Build()
			}
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			mockLoginData := model.LoginData{
				Id:      "102400215",
				Cookies: "cookie1=value1;cookie2=value2",
			}
			ctx := customContext.WithLoginData(context.Background(), &mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, taskqueue.NewBaseTaskQueue())

			result, err := courseService.getLectures(false, &mockLoginData)
			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
