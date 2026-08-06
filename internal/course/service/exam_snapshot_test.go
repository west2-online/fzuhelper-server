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

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/umeng"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestBuildCourseExamInfo(t *testing.T) {
	testCases := []struct {
		name     string
		courses  []*jwch.Course
		expected []CourseExamInfo
	}{
		{
			name: "build exam info and ignore invalid courses",
			courses: []*jwch.Course{
				{
					Name:        "数据结构",
					Teacher:     "张老师",
					Credits:     "4.0",
					RawExamTime: " 2026年6月20日 09:00-11:00 旗山校区 ",
				},
				nil,
				{
					Name:        "无考试课程",
					RawExamTime: " ",
				},
				{
					Name:        "高等数学",
					Teacher:     "李老师",
					Credits:     "5.0",
					RawExamTime: "2026年6月21日 09:00-11:00",
				},
			},
			expected: []CourseExamInfo{
				{
					Name:     "数据结构",
					Teacher:  "张老师",
					Credit:   "4.0",
					ExamTime: "2026年6月20日 09:00-11:00 旗山校区",
				},
				{
					Name:     "高等数学",
					Teacher:  "李老师",
					Credit:   "5.0",
					ExamTime: "2026年6月21日 09:00-11:00",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, buildCourseExamInfo(tc.courses))
		})
	}
}

func TestCourseExamIdentityAndTag(t *testing.T) {
	testCases := []struct {
		name     string
		exam     CourseExamInfo
		identity string
	}{
		{
			name:     "course exam identity",
			exam:     CourseExamInfo{Name: "数据结构", Teacher: "张老师", Credit: "4.0"},
			identity: "数据结构|张老师|4.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.identity, courseExamIdentity(tc.exam))
			assert.Equal(t, utils.MD5(tc.identity), courseExamTag(tc.exam))
		})
	}
}

func TestCourseExamInfoHashIsOrderIndependent(t *testing.T) {
	testCases := []struct {
		name   string
		first  []CourseExamInfo
		second []CourseExamInfo
	}{
		{
			name: "exam order does not affect hash",
			first: []CourseExamInfo{
				{Name: "A", Teacher: "T", Credit: "1", ExamTime: "time-a"},
				{Name: "B", Teacher: "T", Credit: "2", ExamTime: "time-b"},
			},
			second: []CourseExamInfo{
				{Name: "B", Teacher: "T", Credit: "2", ExamTime: "time-b"},
				{Name: "A", Teacher: "T", Credit: "1", ExamTime: "time-a"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			firstHash, err := courseExamInfoHash(tc.first)
			assert.NoError(t, err)
			secondHash, err := courseExamInfoHash(tc.second)
			assert.NoError(t, err)
			assert.Equal(t, firstHash, secondHash)
		})
	}
}

func TestBuildCourseExamChanges(t *testing.T) {
	testCases := []struct {
		name     string
		term     string
		oldExams []CourseExamInfo
		newExams []CourseExamInfo
	}{
		{
			name: "exam time changed",
			term: "202401",
			oldExams: []CourseExamInfo{
				{Name: "数据结构", Teacher: "张老师", Credit: "4.0", ExamTime: "旧时间"},
			},
			newExams: []CourseExamInfo{
				{Name: "数据结构", Teacher: "张老师", Credit: "4.0", ExamTime: "新时间"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changes := buildCourseExamChanges(tc.term, tc.oldExams, tc.newExams)
			if assert.Len(t, changes, 1) {
				assert.Equal(t, courseExamTag(tc.newExams[0]), changes[0].Tag)
				assert.Equal(t, utils.SHA256("数据结构|202401|张老师|4.0|旧时间|新时间"), changes[0].ExamHash)
			}
		})
	}
}

func TestCourseServiceSendExamNotification(t *testing.T) {
	tests := []struct {
		name       string
		androidErr error
		iosErr     error
	}{
		{
			name:       "AndroidErrorIgnored",
			androidErr: assert.AnError,
		},
		{
			name:   "IOSErrorIgnored",
			iosErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		mockey.PatchConvey(tt.name, t, func() {
			mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).To(
				func(pushType, title, text, ticker, tag, description, deeplink string) error {
					assert.Equal(t, constants.UmengPushTypeExam, pushType)
					assert.Equal(t, "fzuhelper://exam-room", deeplink)
					return tt.androidErr
				},
			).Build()
			mockey.Mock(umeng.SendIOSGroupcast).To(
				func(title, subtitle, body, tag, description, deeplink string) error {
					assert.Equal(t, "fzuhelper://exam-room", deeplink)
					return tt.iosErr
				},
			).Build()

			new(CourseService).sendExamNotification(courseExamChange{
				Tag:  utils.MD5("数据结构|张老师|4.0"),
				Exam: CourseExamInfo{Name: "数据结构"},
			})
		})
	}
}
