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
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
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

	type testCase struct {
		name              string
		mockTerms         *jwch.Term
		mockCourses       []*jwch.Course
		mockPutToDbError  error
		expectedResult    []*jwch.Course
		expectingError    bool
		expectedErrorMsg  string
		mockTermsReturn   *jwch.Term
		mockTermsError    error
		mockCoursesReturn []*jwch.Course
		mockCoursesError  error
	}

	// Test cases
	testCases := []testCase{
		{
			name:              "GetCourseListSuccess",
			mockTerms:         mockTerm,
			mockCourses:       mockCourses,
			expectedResult:    mockCourses,
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
	}

	mockLoginData := &model.LoginData{
		Id:      "102301517",
		Cookies: "cookie1=value1; cookie2=value2",
	}
	req := &course.CourseListRequest{
		LoginData: mockLoginData,
		Term:      "202401",
	}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*jwch.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()
			mockey.Mock((*jwch.Student).GetSemesterCourses).Return(tc.mockCoursesReturn, tc.mockCoursesError).Build()
			mockey.Mock((*CourseService).putCourseListToDatabase).Return(tc.mockPutToDbError).Build()

			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			courseService := NewCourseService(context.Background(), mockClientSet)

			result, err := courseService.GetCourseList(req)

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
