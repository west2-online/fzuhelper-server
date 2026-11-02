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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func TestParseExamTimeAndLocation(t *testing.T) {
	testCases := []struct {
		name         string
		raw          string
		expectDate   string
		expectTime   string
		expectLoc    string
		expectParsed bool
	}{
		{
			name:         "real format without spaces",
			raw:          "2026年06月26日14:00-18:00旗山东3-106",
			expectDate:   "2026年06月26日",
			expectTime:   "14:00-18:00",
			expectLoc:    "旗山东3-106",
			expectParsed: true,
		},
		{
			name:         "real format with building name",
			raw:          "2024年12月10日14:00-16:00铜盘科报厅B",
			expectDate:   "2024年12月10日",
			expectTime:   "14:00-16:00",
			expectLoc:    "铜盘科报厅B",
			expectParsed: true,
		},
		{
			name:         "month and day without zero padding",
			raw:          "2026年6月20日09:00-11:00旗山校区",
			expectDate:   "2026年6月20日",
			expectTime:   "09:00-11:00",
			expectLoc:    "旗山校区",
			expectParsed: true,
		},
		{
			name:         "spaces between parts are tolerated",
			raw:          "2026年6月20日 09:00-11:00 旗山校区",
			expectDate:   "2026年6月20日",
			expectTime:   "09:00-11:00",
			expectLoc:    "旗山校区",
			expectParsed: true,
		},
		{
			name:         "exam without location",
			raw:          "2026年6月20日09:00-11:00",
			expectDate:   "2026年6月20日",
			expectTime:   "09:00-11:00",
			expectLoc:    "",
			expectParsed: true,
		},
		{
			name:         "empty text",
			raw:          "",
			expectParsed: false,
		},
		{
			name:         "blank text",
			raw:          "   ",
			expectParsed: false,
		},
		{
			name:         "placeholder text",
			raw:          "暂无",
			expectParsed: false,
		},
		{
			name:         "date only",
			raw:          "2026年6月20日",
			expectParsed: false,
		},
		{
			name:         "time without range",
			raw:          "2026年6月20日09:00",
			expectParsed: false,
		},
		{
			name:         "location only",
			raw:          "旗山东3-106",
			expectParsed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			date, examTime, location, parsed := parseExamTimeAndLocation(tc.raw)
			assert.Equal(t, tc.expectParsed, parsed)
			assert.Equal(t, tc.expectDate, date)
			assert.Equal(t, tc.expectTime, examTime)
			assert.Equal(t, tc.expectLoc, location)
		})
	}
}

func TestBuildExamRoomInfo(t *testing.T) {
	t.Run("empty courses", func(t *testing.T) {
		assert.Equal(t, []*model.ExamRoomInfo{}, buildExamRoomInfo(nil))
	})

	t.Run("courses without exam keep placeholder rows", func(t *testing.T) {
		courses := []*jwch.Course{
			nil,
			{Name: "无考试课程", Credits: "2.0", Teacher: "张老师"},
			{Name: "空白考试课程", Credits: "2.0", Teacher: "李老师", RawExamTime: "  "},
		}
		expected := []*model.ExamRoomInfo{
			{
				Name:     "无考试课程",
				Credit:   "2.0",
				Teacher:  "张老师",
				Location: "暂无考场数据",
			},
			{
				Name:     "空白考试课程",
				Credit:   "2.0",
				Teacher:  "李老师",
				Location: "暂无考场数据",
			},
		}
		assert.Equal(t, expected, buildExamRoomInfo(courses))
	})

	t.Run("build exam room info", func(t *testing.T) {
		courses := []*jwch.Course{
			{
				Name:        "数据结构",
				Credits:     "4.0",
				Teacher:     "张老师",
				RawExamTime: "2026年06月26日14:00-18:00旗山东3-106",
			},
			{
				Name:        "高等数学",
				Credits:     "5.0",
				Teacher:     "李老师",
				RawExamTime: "2026年07月02日09:00-11:00旗山文3-101",
			},
			{
				Name:    "形势与政策（五）",
				Credits: "0.25",
				Teacher: "冯杰,潘雯,陈睿莹,张龙燕",
			},
		}

		expected := []*model.ExamRoomInfo{
			{
				Name:     "数据结构",
				Credit:   "4.0",
				Teacher:  "张老师",
				Date:     "2026年06月26日",
				Time:     "14:00-18:00",
				Location: "旗山东3-106",
			},
			{
				Name:     "高等数学",
				Credit:   "5.0",
				Teacher:  "李老师",
				Date:     "2026年07月02日",
				Time:     "09:00-11:00",
				Location: "旗山文3-101",
			},
			{
				Name:     "形势与政策（五）",
				Credit:   "0.25",
				Teacher:  "冯杰,潘雯,陈睿莹,张龙燕",
				Location: "暂无考场数据",
			},
		}

		assert.Equal(t, expected, buildExamRoomInfo(courses))
	})

	t.Run("duplicate courses produce a single exam room", func(t *testing.T) {
		course := &jwch.Course{
			Name:        "数据结构",
			Credits:     "4.0",
			Teacher:     "张老师",
			RawExamTime: "2026年06月26日14:00-18:00旗山东3-106",
		}

		courses := []*jwch.Course{course, course, course}
		assert.Len(t, buildExamRoomInfo(courses), 1)
	})
}

func TestBuildExamRoomInfoYjsy(t *testing.T) {
	t.Run("empty courses", func(t *testing.T) {
		assert.Equal(t, []*model.ExamRoomInfo{}, buildExamRoomInfoYjsy(nil))
	})

	t.Run("build exam room info", func(t *testing.T) {
		courses := []*yjsy.Course{
			{
				Name:        "数值分析",
				Credits:     "2.0",
				Teacher:     "王老师",
				RawExamTime: "2026年06月22日14:00-16:00旗山东3-101",
			},
			{
				Name:    "无考试课程",
				Credits: "1.0",
				Teacher: "赵老师",
			},
		}

		expected := []*model.ExamRoomInfo{
			{
				Name:     "数值分析",
				Credit:   "2.0",
				Teacher:  "王老师",
				Date:     "2026年06月22日",
				Time:     "14:00-16:00",
				Location: "旗山东3-101",
			},
			{
				Name:     "无考试课程",
				Credit:   "1.0",
				Teacher:  "赵老师",
				Location: "暂无考场数据",
			},
		}

		assert.Equal(t, expected, buildExamRoomInfoYjsy(courses))
	})
}
