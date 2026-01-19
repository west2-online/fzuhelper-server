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

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	commonCache "github.com/west2-online/fzuhelper-server/pkg/cache/common"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func TestGetContributorInfo(t *testing.T) {
	type testCase struct {
		name           string
		mockFileResult map[string][]*model.Contributor
		mockFileError  error
		expectedResult map[string][]*model.Contributor
		expectingError error

		// 新增字段：用于控制缓存的场景
		cacheExist    bool                            // 是否在 Redis 中存在这个 Key
		cacheGetError error                           // 获取缓存时是否模拟报错
		cachedResult  map[string][]*model.Contributor // 如果缓存命中时，要返回的缓存结果
	}

	// 准备完整的 mock 数据，覆盖所有4个 key
	mockContributors := map[string][]*model.Contributor{
		constants.ContributorFzuhelperAppKey: {
			{Name: "Alice", AvatarUrl: "https://example.com/avatar1.png", Url: "app", Contributions: 1},
		},
		constants.ContributorFzuhelperServerKey: {
			{Name: "Bob", AvatarUrl: "https://example.com/avatar2.png", Url: "server", Contributions: 2},
		},
		constants.ContributorJwchKey: {
			{Name: "Charlie", AvatarUrl: "https://example.com/avatar3.png", Url: "Jwch", Contributions: 3},
		},
		constants.ContributorYJSYKey: {
			{Name: "David", AvatarUrl: "https://example.com/avatar4.png", Url: "YJSY", Contributions: 4},
		},
	}
	testCases := []testCase{
		{
			name:           "SuccessCase",
			cacheExist:     true, // 必须设置缓存存在
			cachedResult:   mockContributors,
			expectedResult: mockContributors,
			expectingError: nil,
		},
		{
			name:           "CacheKeyNotExist",
			cacheExist:     false, // 缓存不存在
			mockFileResult: nil,
			mockFileError:  fmt.Errorf("key not found"),
			expectedResult: nil,
			expectingError: fmt.Errorf("service.GetContributorInfo: %s not exist", constants.ContributorFzuhelperAppKey),
		},
		{
			name:           "CacheGetError",
			cacheExist:     true,
			cacheGetError:  fmt.Errorf("cache get error"),
			expectedResult: nil,
			expectingError: fmt.Errorf("service.GetContributorInfo: failed to get contributor info for key %s: %w",
				constants.ContributorFzuhelperAppKey, fmt.Errorf("cache get error")),
		},
	}

	defer mockey.UnPatchAll()

	mockey.Mock((*commonCache.CacheCommon).SetContributorInfo).To(
		func(ctx context.Context, key string, contributors []*model.Contributor) error {
			return nil
		}).Build()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				CacheClient: new(cache.Cache),
			}

			// Mock IsKeyExist: 当 cacheExist=false 时，第一个 key 就会返回 false，导致直接报错
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()

			// Mock GetContributorInfo
			mockey.Mock((*commonCache.CacheCommon).GetContributorInfo).To(
				func(ctx context.Context, key string) ([]*model.Contributor, error) {
					// 如果测试用例要求缓存获取报错
					if tc.cacheGetError != nil {
						return nil, tc.cacheGetError
					}
					// 如果提供了缓存结果，返回对应 key 的数据
					if tc.cachedResult != nil {
						return tc.cachedResult[key], nil
					}
					// 否则返回 mockFileResult 中的数据（如果有）
					if tc.mockFileResult != nil {
						return tc.mockFileResult[key], tc.mockFileError
					}
					return nil, tc.mockFileError
				},
			).Build()

			commonService := NewCommonService(context.Background(), mockClientSet)
			result, err := commonService.GetContributorInfo()

			if tc.expectingError != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectingError.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
