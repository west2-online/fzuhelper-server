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
	"sync"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbcourse "github.com/west2-online/fzuhelper-server/pkg/db/course"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestGetAutoAdjustCourseList(t *testing.T) {
	type testCase struct {
		name                  string
		mockTerm              string
		mockAutoAdjustCourses []*model.AutoAdjustCourse
		mockDBErr             error
		cacheExist            bool
		cacheGetErr           error
		expectResult          []*model.AutoAdjustCourse
		expectError           string
	}

	mockTerm := "202501"

	mockAutoAdjustCourses := []*model.AutoAdjustCourse{
		{
			Year:        "2025",
			FromDate:    "2025-05-01",
			ToDate:      new("2025-05-05"),
			Term:        mockTerm,
			FromWeek:    18,
			ToWeek:      new(int64(19)),
			FromWeekday: 1,
			ToWeekday:   new(int64(3)),
			Enabled:     true,
		},
		{
			Year:        "2025",
			FromDate:    "2025-06-04",
			ToDate:      new("2025-06-07"),
			Term:        mockTerm,
			FromWeek:    20,
			ToWeek:      new(int64(2)),
			FromWeekday: 2,
			ToWeekday:   new(int64(4)),
			Enabled:     true,
		},
	}

	testCases := []testCase{
		{
			name:                  "success",
			mockTerm:              mockTerm,
			mockAutoAdjustCourses: mockAutoAdjustCourses,
			expectResult:          mockAutoAdjustCourses,
		},
		{
			name:                  "cache exist success",
			mockTerm:              mockTerm,
			mockAutoAdjustCourses: mockAutoAdjustCourses,
			cacheExist:            true,
			expectResult:          mockAutoAdjustCourses,
		},
		{
			name:        "cache exist but get cache error",
			mockTerm:    mockTerm,
			cacheExist:  true,
			cacheGetErr: assert.AnError,
			expectError: "service.GetAutoAdjustCourseList: Get cache failed",
		},
		{
			name:        "cache miss get from db error",
			mockTerm:    mockTerm,
			cacheExist:  false,
			mockDBErr:   assert.AnError,
			expectError: "service.GetAutoAdjustCourseList: Get from db failed",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			shouldWait := !tc.cacheExist && tc.mockDBErr == nil
			var wg sync.WaitGroup
			if shouldWait {
				wg.Add(1)
			}

			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()

			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetAutoAdjustCourseListCache).
					Return(tc.mockAutoAdjustCourses, tc.cacheGetErr).
					Build()
			}

			mockey.Mock((*dbcourse.DBCourse).GetAutoAdjustCourseListByTerm).
				Return(tc.mockAutoAdjustCourses, tc.mockDBErr).
				Build()

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetAutoAdjustCourseList(tc.mockTerm)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
