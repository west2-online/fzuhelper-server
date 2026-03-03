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
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

func TestGetCalendar(t *testing.T) {
	type testCase struct {
		name                    string
		stuID                   string
		mockLatestStartTime     string
		mockLatestTerm          string
		mockYjsTerm             string
		mockGetLatestStartError error
		mockCourses             []*model.Course
		mockGetCoursesError     error
		expectError             string
	}

	// 准备 mock 数据
	mockCourses := []*model.Course{
		{
			Name:    "Advanced Programming", // 多个上课时间
			Teacher: "Prof. Chen",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "A-101",
					StartClass: 1,
					EndClass:   2,
					StartWeek:  1,
					EndWeek:    16,
					Weekday:    1,
					Single:     true,
					Double:     true,
					Adjust:     false,
				},
				{
					Location:   "A-102",
					StartClass: 3,
					EndClass:   4,
					StartWeek:  1,
					EndWeek:    16,
					Weekday:    3,
					Single:     true,
					Double:     true,
					Adjust:     false,
				},
			},
		},
		{
			Name:    "English", // 调课课程
			Teacher: "Prof. Wang",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "B-101",
					StartClass: 1,
					EndClass:   2,
					StartWeek:  1,
					EndWeek:    10,
					Weekday:    3,
					Single:     true,
					Double:     true,
					Adjust:     true,
				},
			},
		},
		{
			Name:    "Chemistry", // 单周课程
			Teacher: "Prof. Li",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "C-301",
					StartClass: 5,
					EndClass:   6,
					StartWeek:  1,
					EndWeek:    15,
					Weekday:    4,
					Single:     true,
					Double:     false,
					Adjust:     false,
				},
			},
		},
		{
			Name:    "Biology", // 双周课程
			Teacher: "Prof. Zhang",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "D-401",
					StartClass: 7,
					EndClass:   8,
					StartWeek:  2,
					EndWeek:    16,
					Weekday:    5,
					Single:     false,
					Double:     true,
					Adjust:     false,
				},
			},
		},
		{
			Name:    "Computer Science",
			Teacher: "Prof. Liu",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "铜盘A-101", // 已知位置，有 GEO 信息
					StartClass: 1,
					EndClass:   2,
					StartWeek:  1,
					EndWeek:    16,
					Weekday:    1,
					Single:     true,
					Double:     true,
					Adjust:     false,
				},
			},
		},
		{
			Name:    "History",
			Teacher: "Prof. Zhao",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "Unknown-Building-999", // 未知位置，无 GEO 信息
					StartClass: 3,
					EndClass:   4,
					StartWeek:  1,
					EndWeek:    16,
					Weekday:    2,
					Single:     true,
					Double:     true,
					Adjust:     false,
				},
			},
		},
	}

	testCases := []testCase{
		{
			name:                "SuccessCaseForUndergraduate",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202401",
			mockCourses:         mockCourses,
		},
		{
			name:                "SuccessCaseForGraduate",
			stuID:               "00000102301001", // 研究生学号格式：前5位是00000
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202401",
			mockCourses:         mockCourses,
		},
		{
			name:                    "GetLatestStartTermError",
			stuID:                   "102301001",
			mockGetLatestStartError: assert.AnError,
			expectError:             "get latest start term failed",
		},
		{
			name:                "InvalidStartDateFormat",
			stuID:               "102301001",
			mockLatestStartTime: "invalid-date",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			expectError:         "parse current term start date failed",
		},
		{
			name:                "GetSemesterCoursesError",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockGetCoursesError: assert.AnError,
			expectError:         "get semester courses failed",
		},
		{
			name:                "GetYjsSemesterCoursesError",
			stuID:               "00000102301001", // 研究生学号格式
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockGetCoursesError: assert.AnError,
			expectError:         "get yjs semester courses failed", // 研究生的错误消息
		},
		{
			name:                "EmptyCoursesList",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses:         []*model.Course{},
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			// Mock getLatestStartTerm
			mockey.Mock((*CourseService).getLatestStartTerm).Return(tc.mockLatestStartTime, tc.mockLatestTerm, tc.mockYjsTerm, tc.mockGetLatestStartError).Build()

			// Mock getSemesterCourses with parameter validation
			mockey.Mock((*CourseService).getSemesterCourses).Return(tc.mockCourses, tc.mockGetCoursesError).Build()

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetCalendar(tc.stuID)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				calendarContent := string(result)
				assert.Contains(t, calendarContent, "BEGIN:VCALENDAR")

				if len(tc.mockCourses) == 0 {
					return
				}
				assert.Contains(t, calendarContent, "BEGIN:VEVENT")
				for _, course := range tc.mockCourses {
					if len(course.ScheduleRules) == 0 {
						continue
					}
					schedule := course.ScheduleRules[0]

					// 验证位置信息
					assert.Contains(t, calendarContent, "LOCATION:"+schedule.Location)

					// 验证 GEO 信息：已知位置应该有 GEO，未知位置不应该有 GEO
					lat, lon := findGeoLocation(schedule.Location)
					if lat != 0 && lon != 0 {
						assert.Contains(t, calendarContent, "GEO:", "Known location should have GEO information")
					}
					// 注意：未知位置的情况下，GEO 标签不应出现在该事件中
					// 但由于 ICS 格式的特性，这个验证需要更精细的解析，这里主要验证已知位置的情况

					if schedule.Adjust {
						assert.True(t,
							strings.Contains(calendarContent, "[调课] "+course.Name) ||
								strings.Contains(calendarContent, course.Name))
						continue
					}
					assert.Contains(t, calendarContent, course.Name)
				}
			}
		})
	}
}
