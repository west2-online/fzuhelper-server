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
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestISOWeekCrossWeek(t *testing.T) {
	// 创建一个周日的时间（2025-03-23）
	sunday := time.Date(2025, 3, 23, 0, 0, 0, 0, time.UTC)
	sundayYear, sundayWeek := sunday.ISOWeek()

	// 创建一个周一的时间（2025-03-25）
	monday := time.Date(2025, 3, 24, 0, 0, 0, 0, time.UTC)
	mondayYear, mondayWeek := monday.ISOWeek()

	// 验证周日和周一的周数是否不同
	if sundayWeek == mondayWeek {
		t.Errorf("Expected different weeks for Sunday and Monday, but got same week: Sunday(Week %d) Monday(Week %d)",
			sundayWeek, mondayWeek)
	}

	t.Logf("Sunday: Year %d Week %d", sundayYear, sundayWeek)
	t.Logf("Monday: Year %d Week %d", mondayYear, mondayWeek)
}

func TestCourseService_GetLocateDate(t *testing.T) {
	type testCase struct {
		name          string
		cacheExist    bool
		cacheGetError error
		jwchReturn    *jwch.LocateDate
		jwchError     error
		expectYear    string
		expectWeek    string
		expectTerm    string
		expectError   bool
		expectErrMsg  string
	}

	testCases := []testCase{
		{
			name:       "cache exist success",
			cacheExist: true,
			expectYear: "2025",
			expectWeek: "10",
			expectTerm: "202501",
		},
		{
			name:          "cache exist but get error -> fallback",
			cacheExist:    true,
			cacheGetError: fmt.Errorf("redis error"),
			jwchReturn:    &jwch.LocateDate{Year: "2025", Week: "11", Term: "202501"},
			expectYear:    "2025",
			expectWeek:    "11",
			expectTerm:    "202501",
		},
		{
			name:       "cache miss -> fetch from jwch",
			cacheExist: false,
			jwchReturn: &jwch.LocateDate{Year: "2025", Week: "12", Term: "202501"},
			expectYear: "2025",
			expectWeek: "12",
			expectTerm: "202501",
		},
		{
			name:         "jwch error",
			cacheExist:   false,
			jwchError:    fmt.Errorf("network error"),
			expectError:  true,
			expectErrMsg: "Get locate date fail",
		},
	}

	mockLoginData := &model.LoginData{Id: "052106112", Cookies: "test_cookie"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer mockey.UnPatchAll()

			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			// Mock cache IsKeyExist
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()

			if tc.cacheExist {
				// Mock GetDateCache when cache exists
				mockey.Mock((*coursecache.CacheCourse).GetDateCache).To(
					func(ctx context.Context, key string) (*model.LocateDate, error) {
						if tc.cacheGetError != nil {
							return nil, tc.cacheGetError
						}
						return &model.LocateDate{
							Year: tc.expectYear,
							Week: tc.expectWeek,
							Term: tc.expectTerm,
							Date: time.Now().In(constants.ChinaTZ).Format(time.DateTime),
						}, nil
					},
				).Build()
			}

			// Mock jwch.NewStudent() and GetLocateDate
			mockStudent := &jwch.Student{}
			mockey.Mock(jwch.NewStudent).Return(mockStudent).Build()
			mockey.Mock((*jwch.Student).GetLocateDate).To(
				func(s *jwch.Student) (*jwch.LocateDate, error) {
					return tc.jwchReturn, tc.jwchError
				},
			).Build()

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))

			result, err := courseService.GetLocateDate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tc.expectErrMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectYear, result.Year)
				assert.Equal(t, tc.expectWeek, result.Week)
				assert.Equal(t, tc.expectTerm, result.Term)
				assert.NotEmpty(t, result.Date)
			}
		})
	}
}
