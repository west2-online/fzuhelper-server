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
	"errors"
	"fmt"
	"testing"

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
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    *jwch.SchoolCalendar

		// 新增字段：用于控制缓存的场景
		cacheExist    bool                 // 是否在 Redis 中存在这个 Key
		cacheGetError error                // 获取缓存时是否模拟报错
		cacheCalendar *jwch.SchoolCalendar // 如果缓存命中时，要返回的缓存结果
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
	defer mockey.UnPatchAll()
	testCases := []TestCase{
		{
			Name:              "GetTermListSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
		},
		{
			Name:              "GetTermListError",
			expectedError:     true,
			expectedErrorInfo: errors.New("get term list failed"),
			expectedResult:    nil,
		},
		//// ------------------- 以下为缓存相关测试场景示例 -------------------
		{
			Name:           "cache exist success",
			cacheExist:     true, // 缓存里已存在
			cacheGetError:  nil,  // 获取缓存不报错
			cacheCalendar:  expectedResult,
			expectedResult: expectedResult,
		},
		{
			Name:              "cache exist but get cache error",
			cacheExist:        true,
			cacheGetError:     fmt.Errorf("redis get error"),
			expectedError:     true,
			expectedErrorInfo: errors.New("redis get error"),
		},
	}
	mockey.Mock((*commonCache.CacheCommon).SetTermListCache).To(
		func(ctx context.Context, key string, list *jwch.SchoolCalendar) error {
			return nil
		},
	).Build()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()
			if tc.cacheExist {
				mockey.Mock((*commonCache.CacheCommon).GetTermListCache).To(
					func(ctx context.Context, key string) (*jwch.SchoolCalendar, error) {
						if tc.cacheGetError != nil {
							return nil, tc.cacheGetError
						}
						return tc.cacheCalendar, nil
					},
				).Build()
			} else {
				// 如果缓存不存在，一般不会去调 GetStuInfoCache
				// 也可以不 Mock，或 Mock 一个默认返回
				mockey.Mock((*commonCache.CacheCommon).GetTermListCache).To(
					func(ctx context.Context, key string) (*jwch.SchoolCalendar, error) {
						return nil, fmt.Errorf("should not be called if cache doesn't exist")
					},
				).Build()
			}
			mockey.Mock((*jwch.Student).GetSchoolCalendar).To(func() (*jwch.SchoolCalendar, error) {
				return tc.expectedResult, tc.expectedErrorInfo
			}).Build()
			commonService := NewCommonService(context.Background(), mockClientSet)
			result, err := commonService.GetTermList()
			if tc.expectedError {
				assert.Contains(t, err.Error(), tc.expectedErrorInfo.Error())
				assert.Nil(t, result)
			} else {
				assert.Nil(t, tc.expectedErrorInfo, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetTerm(t *testing.T) {
	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    *jwch.CalTermEvents
		expectedGetInfo   bool
		expectedCached    bool
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
			Name:              "GetTermSuccessfullyWithoutCache",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
			expectedGetInfo:   true,
			expectedCached:    false,
		},
		{
			Name:              "GetTermError",
			expectedError:     true,
			expectedErrorInfo: errors.New("get term events failed"),
			expectedResult:    nil,
			expectedGetInfo:   false,
			expectedCached:    false,
		},
		{
			Name:              "GetTermFromCache",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
			expectedGetInfo:   true,
			expectedCached:    true,
		},
		{
			Name:              "CachedButGetTermError",
			expectedError:     true,
			expectedErrorInfo: errors.New("Get term events failed"),
			expectedResult:    nil,
			expectedGetInfo:   false,
			expectedCached:    true,
		},
		{
			Name:              "SetCacheError",
			expectedError:     true,
			expectedErrorInfo: errors.New("Set term events failed in cache"),
			expectedResult:    nil,
			expectedGetInfo:   false,
			expectedCached:    false,
		},
	}

	defer mockey.UnPatchAll()
	req := &common.TermRequest{Term: "201501"}
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			ClientSet := new(base.ClientSet)
			ClientSet.CacheClient = new(cache.Cache)
			commonService := NewCommonService(context.Background(), ClientSet)
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.expectedCached
			}).Build()
			mockey.Mock((*commonCache.CacheCommon).TermInfoKey).To(func(term string) string {
				return "key"
			}).Build()
			mockey.Mock((*commonCache.CacheCommon).GetTermInfo).To(func(ctx context.Context, key string) (*jwch.CalTermEvents, error) {
				return tc.expectedResult, tc.expectedErrorInfo
			}).Build()
			mockey.Mock((*jwch.Student).GetTermEvents).To(func(termId string) (*jwch.CalTermEvents, error) {
				return tc.expectedResult, tc.expectedErrorInfo
			}).Build()
			mockey.Mock((*commonCache.CacheCommon).SetTermInfo).To(func(ctx context.Context, key string, value *jwch.CalTermEvents) error {
				return tc.expectedErrorInfo
			}).Build()

			success, result, err := commonService.GetTerm(req)
			if tc.expectedError {
				assert.EqualError(t, err, "service.GetTerm: Get term  failed "+tc.expectedErrorInfo.Error())
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedGetInfo, success)
			} else {
				assert.Nil(t, err, tc.expectedErrorInfo)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedGetInfo, success)
			}
		})
	}
}
