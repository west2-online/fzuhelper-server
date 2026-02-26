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
	"sync"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	commonCache "github.com/west2-online/fzuhelper-server/pkg/cache/common"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestGetTermList(t *testing.T) {
	type TestCase struct {
		Name         string
		expectError  error
		expectResult *jwch.SchoolCalendar

		// 新增字段：用于控制缓存的场景
		cacheExist    bool                 // 是否在 Redis 中存在这个 Key
		cacheGetError error                // 获取缓存时是否模拟报错
		cacheCalendar *jwch.SchoolCalendar // 如果缓存命中时，要返回的缓存结果
		setCacheError error                // 设置缓存时是否模拟报错（异步操作）
	}

	expectedResult := &jwch.SchoolCalendar{
		CurrentTerm: "202401",
		Terms: []jwch.CalTerm{
			{
				TermId:     "2024012024082620250117",
				SchoolYear: "2024",
				Term:       "202401",
				StartDate:  "2024-08-26",
				EndDate:    "2025-01-17",
			},
			{
				TermId:     "2024022025022420250704",
				SchoolYear: "2024",
				Term:       "202402",
				StartDate:  "2025-02-24",
				EndDate:    "2025-07-04",
			},
		},
	}

	testCases := []TestCase{
		{
			Name:         "GetTermListSuccessfully",
			expectResult: expectedResult,
		},
		{
			Name:        "GetTermListError",
			expectError: fmt.Errorf("get term list failed"),
		},
		//// ------------------- 以下为缓存相关测试场景示例 -------------------
		{
			Name:          "cache exist success",
			cacheExist:    true, // 缓存里已存在
			cacheCalendar: expectedResult,
			expectResult:  expectedResult,
		},
		{
			Name:          "cache exist but get cache error",
			cacheExist:    true,
			cacheGetError: fmt.Errorf("redis get error"),
			expectError:   fmt.Errorf("redis get error"),
		},
		{
			Name:          "SetTermListCacheError",
			expectResult:  expectedResult,
			setCacheError: fmt.Errorf("cache set failed"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			shouldWait := !tc.cacheExist && tc.expectError == nil
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
				mockey.Mock((*commonCache.CacheCommon).GetTermListCache).Return(tc.cacheCalendar, tc.cacheGetError).Build()
			} else {
				// 如果缓存不存在，一般不会去调 GetTermListCache
				// 也可以不 Mock，或 Mock 一个默认返回
				mockey.Mock((*commonCache.CacheCommon).GetTermListCache).Return(nil, assert.AnError).Build()
			}
			mockey.Mock((*jwch.Student).GetSchoolCalendar).Return(tc.expectResult, tc.expectError).Build()
			setCacheGuard := mockey.Mock((*commonCache.CacheCommon).SetTermListCache).To(func(ctx context.Context, key string, calendar *jwch.SchoolCalendar) error {
				if shouldWait {
					wg.Done()
				}
				return tc.setCacheError
			}).Build()
			defer setCacheGuard.UnPatch()

			commonService := NewCommonService(context.Background(), mockClientSet)
			result, err := commonService.GetTermList()
			if shouldWait && err == nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Fatalf("async cache set did not finish in time")
				}
			}
			if tc.expectError != nil {
				assert.ErrorContains(t, err, tc.expectError.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestGetTerm(t *testing.T) {
	type TestCase struct {
		Name          string
		expectResult  *jwch.CalTermEvents
		expectGetInfo bool
		cacheExist    bool
		cacheGetError error
		apiError      error
		setCacheError error
	}

	expectedResult := &jwch.CalTermEvents{
		TermId:     "201501",
		Term:       "201501",
		SchoolYear: "2015",
		Events: []jwch.CalTermEvent{
			{
				Name:      "学生注册",
				StartDate: "2015-08-29",
				EndDate:   "2015-08-30",
			},
			{
				Name:      "学生补考",
				StartDate: "2015-08-29",
				EndDate:   "2015-09-06",
			},
			{
				Name:      "正式上课",
				StartDate: "2015-08-31",
				EndDate:   "2015-08-31",
			},
			{
				Name:      "新生报到",
				StartDate: "2015-09-07",
				EndDate:   "2015-09-07",
			},
			{
				Name:      "校运会",
				StartDate: "2015-11-12",
				EndDate:   "2015-11-14",
			},
			{
				Name:      "期末考试",
				StartDate: "2016-01-16",
				EndDate:   "2016-01-22",
			},
			{
				Name:      "寒假",
				StartDate: "2016-01-23",
				EndDate:   "2016-02-28",
			},
		},
	}

	testCases := []TestCase{
		{
			Name:          "GetTermSuccessfullyWithoutCache",
			expectResult:  expectedResult,
			expectGetInfo: true,
		},
		{
			Name:     "GetTermError",
			apiError: fmt.Errorf("get term events failed"),
		},
		{
			Name:          "GetTermFromCache",
			expectResult:  expectedResult,
			expectGetInfo: true,
			cacheExist:    true,
		},
		{
			Name:          "CachedButGetTermError",
			cacheExist:    true,
			cacheGetError: fmt.Errorf("Get term cache failed"),
		},
		{
			Name:          "SetCacheError",
			expectGetInfo: true,
			setCacheError: fmt.Errorf("Set term events failed in cache"),
		},
		{
			Name:          "SuccessWithCacheSaveNoError",
			expectResult:  expectedResult,
			expectGetInfo: true,
		},
	}

	req := &common.TermRequest{Term: "201501"}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockClientSet := &base.ClientSet{
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			mockey.Mock((*commonCache.CacheCommon).TermInfoKey).Return("key").Build()
			if tc.cacheExist {
				mockey.Mock((*commonCache.CacheCommon).GetTermInfo).Return(expectedResult, tc.cacheGetError).Build()
			}
			mockey.Mock((*jwch.Student).GetTermEvents).Return(expectedResult, tc.apiError).Build()
			mockey.Mock((*commonCache.CacheCommon).SetTermInfo).Return(tc.setCacheError).Build()

			commonService := NewCommonService(context.Background(), mockClientSet)
			success, result, err := commonService.GetTerm(req)
			if tc.expectResult == nil {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
			assert.Equal(t, tc.expectGetInfo, success)
		})
	}
}
