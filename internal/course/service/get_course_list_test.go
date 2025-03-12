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

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestCourseService_GetCourseList(t *testing.T) {
	mockTerm := &jwch.Term{
		Terms:           []string{"202401"},
		ViewState:       "viewstate123",
		EventValidation: "eventvalidation123",
	}

	mockCourses := []*jwch.Course{
		{
			Type:    "Required",
			Name:    "Mathematics",
			Credits: "3.0",
			Teacher: "Prof. John",
			ScheduleRules: []jwch.CourseScheduleRule{
				{
					Location:     "A-202",
					StartClass:   2,
					EndClass:     4,
					StartWeek:    1,
					EndWeek:      16,
					Weekday:      1,
					Single:       false,
					Double:       true,
					Adjust:       false,
					FromFullWeek: false,
				},
			},
		},
		{
			Type:    "Elective",
			Name:    "Physics",
			Credits: "3.0",
			Teacher: "Prof. Smith",
			ScheduleRules: []jwch.CourseScheduleRule{
				{
					Location:     "A-203",
					StartClass:   3,
					EndClass:     4,
					StartWeek:    2,
					EndWeek:      17,
					Weekday:      2,
					Single:       false,
					Double:       true,
					Adjust:       false,
					FromFullWeek: false,
				},
			},
		},
	}

	mockResult := []*model.Course{
		{
			Name:    "Mathematics",
			Teacher: "Prof. John",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "A-202",
					StartClass: 2,
					EndClass:   4,
					StartWeek:  1,
					EndWeek:    16,
					Weekday:    1,
					Single:     false,
					Double:     true,
					Adjust:     false,
				},
			},
		},
		{
			Name:    "Physics",
			Teacher: "Prof. Smith",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "A-203",
					StartClass: 3,
					EndClass:   4,
					StartWeek:  2,
					EndWeek:    17,
					Weekday:    2,
					Single:     false,
					Double:     true,
					Adjust:     false,
				},
			},
		},
	}
	type testCase struct {
		name              string
		mockTerms         *jwch.Term
		mockCourses       []*jwch.Course
		expectedResult    []*model.Course
		expectingError    bool
		expectedErrorMsg  string
		mockTermsReturn   *jwch.Term
		mockTermsError    error
		mockCoursesReturn []*jwch.Course
		mockCoursesError  error
		cacheExist        bool
		cacheGetError     error
	}

	// Test cases
	testCases := []testCase{
		{
			name:              "GetCourseListSuccess",
			mockTerms:         mockTerm,
			mockCourses:       mockCourses,
			expectedResult:    mockResult,
			expectingError:    false,
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: mockCourses,
		},
		{
			name:              "GetCourseListInvalidTerm",
			mockTerms:         mockTerm,
			mockCourses:       nil,
			expectedResult:    nil,
			expectingError:    true,
			expectedErrorMsg:  "Invalid term",
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: nil,
			mockCoursesError:  fmt.Errorf("Invalid term"),
		},
		{
			name:             "GetCourseListGetTermsFailed",
			mockTerms:        nil,
			mockCourses:      nil,
			expectedResult:   nil,
			expectingError:   true,
			expectedErrorMsg: "Get terms failed",
			mockTermsReturn:  nil,
			mockTermsError:   fmt.Errorf("Get terms failed"),
		},
		{
			name:              "GetCourseListGetCoursesFailed",
			mockTerms:         mockTerm,
			mockCourses:       nil,
			expectedResult:    nil,
			expectingError:    true,
			expectedErrorMsg:  "Get semester courses failed",
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: nil,
			mockCoursesError:  fmt.Errorf("Get semester courses failed"),
		},
		{
			name:           "cache exist success",
			cacheExist:     true, // 缓存里已存在
			cacheGetError:  nil,  // 获取缓存不报错
			expectedResult: mockResult,
		},
	}

	mockLoginData := &model.LoginData{
		Id:      "102301517",
		Cookies: "cookie1=value1; cookie2=value2",
	}
	req := &course.CourseListRequest{
		Term: "202401",
	}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*jwch.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()
			mockey.Mock((*jwch.Student).GetSemesterCourses).Return(tc.mockCoursesReturn, tc.mockCoursesError).Build()
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						if tc.cacheGetError != nil {
							return nil, tc.cacheGetError
						}
						return mockTerm.Terms, nil
					},
				).Build()
				mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).To(
					func(ctx context.Context, key string) ([]*jwch.Course, error) {
						if tc.cacheGetError != nil {
							return nil, tc.cacheGetError
						}
						return mockCourses, nil
					},
				).Build()
			} else {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						return nil, fmt.Errorf("should not be called if cache doesn't exist")
					},
				).Build()
			}
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))

			result, err := courseService.GetCourseList(req, &model.LoginData{Id: "123456789", Cookies: "cookie1=value1;cookie2=value2"})

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
