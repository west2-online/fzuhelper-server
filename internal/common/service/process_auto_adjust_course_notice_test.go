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

	"github.com/west2-online/fzuhelper-server/pkg/ai"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbcourse "github.com/west2-online/fzuhelper-server/pkg/db/course"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestProcessAutoAdjustCourseNotice(t *testing.T) {
	type testCase struct {
		name            string
		info            *jwch.NoticeInfo
		noticeDetailErr error
		aiResult        *ai.AutoAdjustCourseOutput
		aiErr           error
		termList        *jwch.SchoolCalendar
		termListErr     error
		findTermResult  jwch.CalTerm
		findTermFound   bool
		weekdayErr      error
		createErr       error
		getListErr      error
		setCacheErr     error
		expectError     string
	}

	mockTerm := jwch.CalTerm{
		Term:      "202501",
		StartDate: "2025-02-17",
		EndDate:   "2025-06-30",
	}
	mockCalendar := &jwch.SchoolCalendar{
		CurrentTerm: "202501",
		Terms:       []jwch.CalTerm{mockTerm},
	}
	mockAiResult := &ai.AutoAdjustCourseOutput{
		Items: []ai.AutoAdjustCourseItem{
			{FromDate: "2025-05-01", ToDate: "2025-05-08"},
		},
	}
	mockAiResultCancelled := &ai.AutoAdjustCourseOutput{
		Items: []ai.AutoAdjustCourseItem{
			{FromDate: "2025-05-01", ToDate: ""},
		},
	}
	mockNoticeInfo := &jwch.NoticeInfo{
		Title:    "关于课程调整的通知",
		WbTreeId: "1036",
		WbNewsId: "12345",
	}
	mockNoticeDetail := &jwch.NoticeDetail{
		Content: "课程调整内容",
	}

	testCases := []testCase{
		{
			name: "not a course adjust notice",
			info: &jwch.NoticeInfo{Title: "其他通知"},
		},
		{
			name:            "get notice detail error",
			info:            mockNoticeInfo,
			noticeDetailErr: assert.AnError,
			expectError:     "failed to get notice detail",
		},
		{
			name:        "ai error",
			info:        mockNoticeInfo,
			aiErr:       assert.AnError,
			expectError: "failed to auto adjust course",
		},
		{
			name:        "get term list error",
			info:        mockNoticeInfo,
			aiResult:    mockAiResult,
			termListErr: assert.AnError,
			expectError: "failed to get term list",
		},
		{
			name:     "success with no items",
			info:     mockNoticeInfo,
			aiResult: &ai.AutoAdjustCourseOutput{Items: []ai.AutoAdjustCourseItem{}},
			termList: mockCalendar,
		},
		{
			name: "item skipped due to invalid from_date",
			info: mockNoticeInfo,
			aiResult: &ai.AutoAdjustCourseOutput{
				Items: []ai.AutoAdjustCourseItem{
					{FromDate: "not-a-date", ToDate: ""},
				},
			},
			termList: mockCalendar,
		},
		{
			name: "item skipped due to invalid to_date",
			info: mockNoticeInfo,
			aiResult: &ai.AutoAdjustCourseOutput{
				Items: []ai.AutoAdjustCourseItem{
					{FromDate: "2025-05-01", ToDate: "not-a-date"},
				},
			},
			termList:       mockCalendar,
			findTermResult: mockTerm,
			findTermFound:  true,
		},
		{
			name:          "item skipped due to no term found",
			info:          mockNoticeInfo,
			aiResult:      mockAiResult,
			termList:      mockCalendar,
			findTermFound: false,
		},
		{
			name:           "item skipped due to get weekday error for from_date",
			info:           mockNoticeInfo,
			aiResult:       mockAiResult,
			termList:       mockCalendar,
			findTermResult: mockTerm,
			findTermFound:  true,
			weekdayErr:     assert.AnError,
		},
		{
			name:           "create auto adjust course error",
			info:           mockNoticeInfo,
			aiResult:       mockAiResult,
			termList:       mockCalendar,
			findTermResult: mockTerm,
			findTermFound:  true,
			createErr:      assert.AnError,
			expectError:    "failed to create auto adjust course",
		},
		{
			name:           "get auto adjust course list error during cache refresh",
			info:           mockNoticeInfo,
			aiResult:       mockAiResult,
			termList:       mockCalendar,
			findTermResult: mockTerm,
			findTermFound:  true,
			getListErr:     assert.AnError,
			expectError:    "failed to get auto adjust course list",
		},
		{
			name:           "set cache error during cache refresh",
			info:           mockNoticeInfo,
			aiResult:       mockAiResult,
			termList:       mockCalendar,
			findTermResult: mockTerm,
			findTermFound:  true,
			setCacheErr:    assert.AnError,
			expectError:    "failed to cache auto adjust course list",
		},
		{
			name:           "success with to_date set",
			info:           mockNoticeInfo,
			aiResult:       mockAiResult,
			termList:       mockCalendar,
			findTermResult: mockTerm,
			findTermFound:  true,
		},
		{
			name:           "success with to_date empty (course canceled)",
			info:           mockNoticeInfo,
			aiResult:       mockAiResultCancelled,
			termList:       mockCalendar,
			findTermResult: mockTerm,
			findTermFound:  true,
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*jwch.Student).GetNoticeDetail).Return(mockNoticeDetail, tc.noticeDetailErr).Build()
			mockey.Mock(ai.AutoAdjustCourse).Return(tc.aiResult, tc.aiErr).Build()
			mockey.Mock((*CommonService).GetTermList).Return(tc.termList, tc.termListErr).Build()
			mockey.Mock(utils.FindTermByDate).Return(tc.findTermResult, tc.findTermFound).Build()
			mockey.Mock(utils.GetWeekdayByDate).Return(18, 1, tc.weekdayErr).Build()
			mockey.Mock((*dbcourse.DBCourse).CreateAutoAdjustCourse).Return(nil, tc.createErr).Build()
			mockey.Mock((*dbcourse.DBCourse).GetAutoAdjustCourseListByTerm).Return(nil, tc.getListErr).Build()
			mockey.Mock((*coursecache.CacheCourse).SetAutoAdjustCourseListCache).Return(tc.setCacheErr).Build()
			mockey.Mock((*coursecache.CacheCourse).AutoAdjustCourseKey).Return("key").Build()

			commonService := NewCommonService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			err := commonService.ProcessAutoAdjustCourseNotice(tc.info)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
