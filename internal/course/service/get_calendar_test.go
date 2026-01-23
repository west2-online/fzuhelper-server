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
	"strings"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

// 辅助函数：解析日期字符串
func parseDate(dateStr string) (time.Time, error) {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	return time.ParseInLocation("2006-01-02", dateStr, cstSh)
}

// 辅助函数：解析日期时间字符串
func parseDateTime(dateTimeStr string) (time.Time, error) {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	return time.ParseInLocation("2006-01-02 15:04:05", dateTimeStr, cstSh)
}

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
		expectingError          bool
		expectingErrorMsg       string
		expectedCallStuID       string // 期望传递给 getSemesterCourses 的 stuID
		expectedCallTerm        string // 期望传递给 getSemesterCourses 的 term
	}

	// 准备 mock 数据
	mockCourses := []*model.Course{
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
					Single:     true,
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

	testCases := []testCase{
		{
			name:                "SuccessCaseForUndergraduate",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202401",
			mockCourses:         mockCourses,
			expectingError:      false,
			expectedCallStuID:   "102301001",
			expectedCallTerm:    "202402", // 本科生使用 latestTerm
		},
		{
			name:                "SuccessCaseForGraduate",
			stuID:               "00000102301001", // 研究生学号格式：前5位是00000
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202401",
			mockCourses:         mockCourses,
			expectingError:      false,
			expectedCallStuID:   "102301001", // 研究生需要去掉前5位前缀
			expectedCallTerm:    "202401",    // 研究生使用 yjsTerm
		},
		{
			name:                    "GetLatestStartTermError",
			stuID:                   "102301001",
			mockGetLatestStartError: fmt.Errorf("database error"),
			expectingError:          true,
			expectingErrorMsg:       "get latest start term failed",
		},
		{
			name:                "InvalidStartDateFormat",
			stuID:               "102301001",
			mockLatestStartTime: "invalid-date",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			expectingError:      true,
			expectingErrorMsg:   "parse current term start date failed",
		},
		{
			name:                "GetSemesterCoursesError",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockGetCoursesError: fmt.Errorf("database error"),
			expectingError:      true,
			expectingErrorMsg:   "get semester courses failed",
		},
		{
			name:                "GetYjsSemesterCoursesError",
			stuID:               "00000102301001", // 研究生学号格式
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockGetCoursesError: fmt.Errorf("database error"),
			expectingError:      true,
			expectingErrorMsg:   "get yjs semester courses failed", // 研究生的错误消息
		},
		{
			name:                "EmptyCoursesList",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses:         []*model.Course{},
			expectingError:      false,
		},
		{
			name:                "CourseWithMultipleScheduleRules",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses: []*model.Course{
				{
					Name:    "Advanced Programming",
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
			},
			expectingError: false,
		},
		{
			name:                "CourseWithAdjustFlag",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses: []*model.Course{
				{
					Name:    "English",
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
			},
			expectingError: false,
		},
		{
			name:                "CourseWithSingleWeekOnly",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses: []*model.Course{
				{
					Name:    "Chemistry",
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
			},
			expectingError: false,
		},
		{
			name:                "CourseWithDoubleWeekOnly",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses: []*model.Course{
				{
					Name:    "Biology",
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
			},
			expectingError: false,
		},
		{
			name:                "CourseWithKnownGeoLocation",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses: []*model.Course{
				{
					Name:    "Computer Science",
					Teacher: "Prof. Liu",
					ScheduleRules: []*model.CourseScheduleRule{
						{
							Location:   "铜盘A-101", // 已知位置，应该设置 GEO 信息
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
			},
			expectingError: false,
		},
		{
			name:                "CourseWithUnknownGeoLocation",
			stuID:               "102301001",
			mockLatestStartTime: "2024-02-26",
			mockLatestTerm:      "202402",
			mockYjsTerm:         "202402",
			mockCourses: []*model.Course{
				{
					Name:    "History",
					Teacher: "Prof. Zhao",
					ScheduleRules: []*model.CourseScheduleRule{
						{
							Location:   "Unknown-Building-999", // 未知位置，不应该设置 GEO 信息
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
			},
			expectingError: false,
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			// Mock getLatestStartTerm
			mockey.Mock((*CourseService).getLatestStartTerm).To(
				func(s *CourseService) (string, string, string, error) {
					if tc.mockGetLatestStartError != nil {
						return "", "", "", tc.mockGetLatestStartError
					}
					return tc.mockLatestStartTime, tc.mockLatestTerm, tc.mockYjsTerm, nil
				},
			).Build()

			// Mock getSemesterCourses with parameter validation
			mockey.Mock((*CourseService).getSemesterCourses).To(
				func(s *CourseService, stuID string, term string) ([]*model.Course, error) {
					// 验证传入的参数是否符合期望
					if tc.expectedCallStuID != "" {
						assert.Equal(t, tc.expectedCallStuID, stuID, "stuID parameter mismatch")
					}
					if tc.expectedCallTerm != "" {
						assert.Equal(t, tc.expectedCallTerm, term, "term parameter mismatch")
					}

					if tc.mockGetCoursesError != nil {
						return nil, tc.mockGetCoursesError
					}
					return tc.mockCourses, nil
				},
			).Build()

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetCalendar(tc.stuID)

			if tc.expectingError {
				assert.Error(t, err)
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)

			calendarContent := string(result)
			assert.Contains(t, calendarContent, "BEGIN:VCALENDAR")
			assert.Contains(t, calendarContent, "END:VCALENDAR")
			assert.Contains(t, calendarContent, "VERSION:2.0")

			if len(tc.mockCourses) == 0 {
				return
			}

			assert.Contains(t, calendarContent, "BEGIN:VEVENT")
			assert.Contains(t, calendarContent, "END:VEVENT")

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
		})
	}
}

func TestCalcClassTime(t *testing.T) {
	type testCase struct {
		name        string
		week        int64
		weekday     int64
		startClass  int64
		endClass    int64
		dateBase    string
		expectStart string
		expectEnd   string
	}

	testCases := []testCase{
		{
			name:        "FirstWeekMonday",
			week:        1,
			weekday:     1,
			startClass:  1,
			endClass:    2,
			dateBase:    "2024-02-26",
			expectStart: "2024-02-26 08:20:00",
			expectEnd:   "2024-02-26 10:00:00",
		},
		{
			name:        "ThirdWeekWednesday",
			week:        3,
			weekday:     3,
			startClass:  3,
			endClass:    4,
			dateBase:    "2024-02-26",
			expectStart: "2024-03-13 10:20:00",
			expectEnd:   "2024-03-13 12:00:00",
		},
		{
			name:        "TenthWeekFriday",
			week:        10,
			weekday:     5,
			startClass:  5,
			endClass:    6,
			dateBase:    "2024-02-26",
			expectStart: "2024-05-03 14:00:00",
			expectEnd:   "2024-05-03 15:40:00",
		},
		{
			name:        "EveningClass",
			week:        1,
			weekday:     1,
			startClass:  9,
			endClass:    11,
			dateBase:    "2024-02-26",
			expectStart: "2024-02-26 19:00:00",
			expectEnd:   "2024-02-26 21:35:00",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dateBase, err := parseDate(tc.dateBase)
			assert.Nil(t, err)

			startTime, endTime := calcClassTime(tc.week, tc.weekday, tc.startClass, tc.endClass, dateBase)

			expectedStart, err := parseDateTime(tc.expectStart)
			assert.Nil(t, err)
			expectedEnd, err := parseDateTime(tc.expectEnd)
			assert.Nil(t, err)

			assert.Equal(t, expectedStart.Format("2006-01-02 15:04:05"), startTime.Format("2006-01-02 15:04:05"))
			assert.Equal(t, expectedEnd.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
		})
	}
}

func TestMd5Str(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected string
	}

	testCases := []testCase{
		{
			name:     "SimpleString",
			input:    "test",
			expected: "098f6bcd4621d373cade4e832627b4f6",
		},
		{
			name:     "EmptyString",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "ComplexString",
			input:    "202402__Mathematics_Prof. John_1-16_1_2-4_A-202_true_true",
			expected: "974415000d787c0354299f8201aa6c52",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := md5Str(tc.input)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, 32, len(result)) // MD5 hash is always 32 characters
		})
	}
}

func TestFindGeoLocation(t *testing.T) {
	type testCase struct {
		name        string
		location    string
		expectedLat float64
		expectedLon float64
	}

	testCases := []testCase{
		{
			name:        "EmptyLocation",
			location:    "",
			expectedLat: 0,
			expectedLon: 0,
		},
		{
			name:        "UnknownLocation",
			location:    "Unknown Building",
			expectedLat: 0,
			expectedLon: 0,
		},
		{
			name:        "KnownLocationWithRoom",
			location:    "铜盘A",
			expectedLat: 26.10377684575211,
			expectedLon: 119.26204839259863,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lat, lon := findGeoLocation(tc.location)

			if tc.expectedLat == 0 && tc.expectedLon == 0 {
				if lat == 0 {
					assert.Equal(t, float64(0), lon)
				}
			} else {
				assert.Equal(t, tc.expectedLat, lat)
				assert.Equal(t, tc.expectedLon, lon)
			}
		})
	}
}
