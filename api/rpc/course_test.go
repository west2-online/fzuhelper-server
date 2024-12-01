package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestGetCourseListRPC(t *testing.T) {
	type testCase struct {
		name           string
		mockResp       []*model.Course
		mockError      error
		expectedResult []*model.Course
		expectingError bool
	}

	testCases := []testCase{
		{
			name: "GetCourseListSuccess",
			mockResp: []*model.Course{
				{
					Name:    "高级语言程序设计实践",
					Teacher: "孙岚",
					ScheduleRules: []*model.CourseScheduleRule{
						{
							Location:   "铜盘A205",
							StartClass: 3,
							EndClass:   4,
							StartWeek:  8,
							EndWeek:    16,
							Weekday:    3,
							Single:     true,
							Double:     true,
							Adjust:     false,
						},
						{
							Location:   "铜盘A508",
							StartClass: 7,
							EndClass:   8,
							StartWeek:  8,
							EndWeek:    16,
							Weekday:    5,
							Single:     true,
							Double:     true,
							Adjust:     false,
						},
					},
					RawScheduleRules: "08-16 星期3:3-4节 铜盘A205\\n08-16 星期5:7-8节 铜盘A508\\n",
					RawAdjust:        "",
					Remark:           "实践课教师必须与理论课教师相同。",
					Syllabus:         "https://jwcjwxt2.fzu.edu.cn:81/pyfa/jxdg/TeachingProgram_view.aspx?kcdm=01000100",
					Lessonplan:       "https://jwcjwxt2.fzu.edu.cn:81/pyfa/skjh/TeachingPlan_view.aspx?kkhm=20240101000100003",
				},
			},
			mockError: nil,
			expectedResult: []*model.Course{
				{
					Name:    "高级语言程序设计实践",
					Teacher: "孙岚",
					ScheduleRules: []*model.CourseScheduleRule{
						{
							Location:   "铜盘A205",
							StartClass: 3,
							EndClass:   4,
							StartWeek:  8,
							EndWeek:    16,
							Weekday:    3,
							Single:     true,
							Double:     true,
							Adjust:     false,
						},
						{
							Location:   "铜盘A508",
							StartClass: 7,
							EndClass:   8,
							StartWeek:  8,
							EndWeek:    16,
							Weekday:    5,
							Single:     true,
							Double:     true,
							Adjust:     false,
						},
					},
					RawScheduleRules: "08-16 星期3:3-4节 铜盘A205\\n08-16 星期5:7-8节 铜盘A508\\n",
					RawAdjust:        "",
					Remark:           "实践课教师必须与理论课教师相同。",
					Syllabus:         "https://jwcjwxt2.fzu.edu.cn:81/pyfa/jxdg/TeachingProgram_view.aspx?kcdm=01000100",
					Lessonplan:       "https://jwcjwxt2.fzu.edu.cn:81/pyfa/skjh/TeachingPlan_view.aspx?kkhm=20240101000100003",
				},
			},
		},
		{
			name:           "GetCourseListRPCError",
			mockResp:       nil,
			mockError:      fmt.Errorf("RPC call failed"),
			expectedResult: nil,
			expectingError: true,
		},
	}

	req := &course.CourseListRequest{
		Term: "202401",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(GetCourseListRPC).Return(tc.mockResp, tc.mockError).Build()
			result, err := GetCourseListRPC(context.Background(), req)
			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
