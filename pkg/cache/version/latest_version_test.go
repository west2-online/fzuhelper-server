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

package version

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base/environment"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func makeVersionHistoryForCache() *model.VersionHistory {
	return &model.VersionHistory{
		Id:        1,
		Version:   "1.0.0",
		Code:      "100",
		Url:       "https://example.com/release.apk",
		Feature:   "init",
		Force:     true,
		Type:      "release",
		CreatedAt: time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC),
	}
}

type redisProcessHook struct {
	process func(context.Context, redis.Cmder) error
}

func (h redisProcessHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (h redisProcessHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if h.process != nil {
			return h.process(ctx, cmd)
		}
		return next(ctx, cmd)
	}
}

func (h redisProcessHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func newTestCacheVersion(process func(context.Context, redis.Cmder) error) *CacheVersion {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	client.AddHook(redisProcessHook{process: process})
	return NewCacheVersion(client)
}

func TestCacheVersion_GetLatestVersionCache(t *testing.T) {
	defer mockey.UnPatchAll()

	versionHistory := makeVersionHistoryForCache()
	data, err := sonic.Marshal(versionHistory)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		mockResult  *redis.StringCmd
		expectError string
	}{
		{
			name:       "success",
			mockResult: redis.NewStringResult(string(data), nil),
		},
		{
			name:        "redis error",
			mockResult:  redis.NewStringResult("", redis.Nil),
			expectError: "cache failed",
		},
		{
			name:        "unmarshal error",
			mockResult:  redis.NewStringResult("{invalid-json", nil),
			expectError: "Unmarshal failed",
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			cacheVersion := newTestCacheVersion(func(_ context.Context, cmd redis.Cmder) error {
				assert.Equal(t, []interface{}{"get", constants.LatestVersionCachePrefix + "release"}, cmd.Args())
				if tc.mockResult.Err() != nil {
					return tc.mockResult.Err()
				}
				stringCmd, ok := cmd.(*redis.StringCmd)
				assert.True(t, ok)
				stringCmd.SetVal(tc.mockResult.Val())
				return nil
			})
			defer cacheVersion.client.Close()

			res, err := cacheVersion.GetLatestVersionCache(context.Background(), "release")
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, versionHistory.Version, res.Version)
			assert.Equal(t, versionHistory.Code, res.Code)
			assert.Equal(t, versionHistory.Type, res.Type)
		})
	}
}

func TestCacheVersion_SetLatestVersionCache(t *testing.T) {
	defer mockey.UnPatchAll()

	versionHistory := makeVersionHistoryForCache()

	testCases := []struct {
		name        string
		testEnv     bool
		mockSetErr  error
		mockMarshal bool
		expectSet   bool
		expectError string
	}{
		{
			name:    "skip in test environment",
			testEnv: true,
		},
		{
			name:      "success",
			expectSet: true,
		},
		{
			name:        "marshal error",
			mockMarshal: true,
			expectError: "Marshal failed",
		},
		{
			name:        "redis error",
			mockSetErr:  fmt.Errorf("redis down"),
			expectSet:   true,
			expectError: "Set cache failed",
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			setCalled := false
			cacheVersion := newTestCacheVersion(func(_ context.Context, cmd redis.Cmder) error {
				setCalled = true
				args := cmd.Args()
				assert.GreaterOrEqual(t, len(args), 3)
				assert.Equal(t, "set", args[0])
				assert.Equal(t, constants.LatestVersionCachePrefix+"release", args[1])
				assert.NotEmpty(t, args[2])
				statusCmd, ok := cmd.(*redis.StatusCmd)
				assert.True(t, ok)
				statusCmd.SetVal("OK")
				return tc.mockSetErr
			})
			defer cacheVersion.client.Close()

			mockey.Mock(environment.IsTestEnvironment).Return(tc.testEnv).Build()
			if tc.mockMarshal {
				mockey.Mock(sonic.Marshal).Return(nil, fmt.Errorf("marshal failed")).Build()
			}

			err := cacheVersion.SetLatestVersionCache(context.Background(), versionHistory)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectSet, setCalled)
		})
	}
}

func TestCacheVersion_DeleteLatestVersionCache(t *testing.T) {
	defer mockey.UnPatchAll()

	testCases := []struct {
		name        string
		mockDelErr  error
		expectError string
	}{
		{
			name: "success",
		},
		{
			name:        "redis error",
			mockDelErr:  fmt.Errorf("redis down"),
			expectError: "cache failed",
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			cacheVersion := newTestCacheVersion(func(_ context.Context, cmd redis.Cmder) error {
				assert.Equal(t, []interface{}{"del", constants.LatestVersionCachePrefix + "release"}, cmd.Args())
				intCmd, ok := cmd.(*redis.IntCmd)
				assert.True(t, ok)
				intCmd.SetVal(1)
				return tc.mockDelErr
			})
			defer cacheVersion.client.Close()

			err := cacheVersion.DeleteLatestVersionCache(context.Background(), "release")
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
