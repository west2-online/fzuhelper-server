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

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	versioncache "github.com/west2-online/fzuhelper-server/pkg/cache/version"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	dbversion "github.com/west2-online/fzuhelper-server/pkg/db/version"
)

func TestGetReleaseVersion(t *testing.T) {
	type testCase struct {
		name            string
		mockCacheReturn *model.VersionHistory
		mockCacheError  error
		mockDBReturn    *model.VersionHistory
		mockDBError     error
		expectVersion   string
		expectError     string
	}

	mockVersion := &model.VersionHistory{
		Id: 1, Version: "1.0.0", Code: "100", Url: "http://example.com/release.apk",
		Feature: "new", Force: false, Type: apkTypeRelease,
	}

	testCases := []testCase{
		{
			name:            "cache hit",
			mockCacheReturn: mockVersion,
			mockCacheError:  nil,
			expectVersion:   "1.0.0",
		},
		{
			name:            "cache miss — db success",
			mockCacheReturn: nil,
			mockCacheError:  fmt.Errorf("redis: nil"),
			mockDBReturn:    mockVersion,
			mockDBError:     nil,
			expectVersion:   "1.0.0",
		},
		{
			name:            "db not found",
			mockCacheReturn: nil,
			mockCacheError:  fmt.Errorf("redis: nil"),
			mockDBReturn:    nil,
			mockDBError:     nil,
			expectError:     "no release version found",
		},
		{
			name:            "db error",
			mockCacheReturn: nil,
			mockCacheError:  fmt.Errorf("redis: nil"),
			mockDBReturn:    nil,
			mockDBError:     fmt.Errorf("db down"),
			expectError:     "db query error",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			cs := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*versioncache.CacheVersion).GetLatestVersionCache).Return(tc.mockCacheReturn, tc.mockCacheError).Build()
			mockey.Mock((*dbversion.DBVersion).GetLatestVersionByType).Return(tc.mockDBReturn, tc.mockDBError).Build()
			mockey.Mock((*versioncache.CacheVersion).SetLatestVersionCache).Return(nil).Build()

			svc := NewVersionService(context.Background(), cs)
			result, err := svc.GetReleaseVersion()

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectVersion, result.Version)
			}
		})
	}
}

func TestGetBetaVersion(t *testing.T) {
	type testCase struct {
		name            string
		mockCacheReturn *model.VersionHistory
		mockCacheError  error
		mockDBReturn    *model.VersionHistory
		mockDBError     error
		expectVersion   string
		expectError     string
	}

	mockVersion := &model.VersionHistory{
		Id: 1, Version: "2.0.0", Code: "200", Url: "http://example.com/beta.apk",
		Feature: "beta feature", Force: true, Type: apkTypeBeta,
	}

	testCases := []testCase{
		{
			name:            "cache hit",
			mockCacheReturn: mockVersion,
			mockCacheError:  nil,
			expectVersion:   "2.0.0",
		},
		{
			name:            "cache miss — db success",
			mockCacheReturn: nil,
			mockCacheError:  fmt.Errorf("redis: nil"),
			mockDBReturn:    mockVersion,
			mockDBError:     nil,
			expectVersion:   "2.0.0",
		},
		{
			name:            "db not found",
			mockCacheReturn: nil,
			mockCacheError:  fmt.Errorf("redis: nil"),
			mockDBReturn:    nil,
			mockDBError:     nil,
			expectError:     "no beta version found",
		},
		{
			name:            "db error",
			mockCacheReturn: nil,
			mockCacheError:  fmt.Errorf("redis: nil"),
			mockDBReturn:    nil,
			mockDBError:     fmt.Errorf("db down"),
			expectError:     "db query error",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			cs := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*versioncache.CacheVersion).GetLatestVersionCache).Return(tc.mockCacheReturn, tc.mockCacheError).Build()
			mockey.Mock((*dbversion.DBVersion).GetLatestVersionByType).Return(tc.mockDBReturn, tc.mockDBError).Build()
			mockey.Mock((*versioncache.CacheVersion).SetLatestVersionCache).Return(nil).Build()

			svc := NewVersionService(context.Background(), cs)
			result, err := svc.GetBetaVersion()

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectVersion, result.Version)
			}
		})
	}
}
