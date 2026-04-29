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
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	versioncache "github.com/west2-online/fzuhelper-server/pkg/cache/version"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	dbversion "github.com/west2-online/fzuhelper-server/pkg/db/version"
)

func mockTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", "2026-04-26 12:00:00")
	return t
}

func TestGetVersionHistoryList(t *testing.T) {
	type testCase struct {
		name        string
		mockReturn  []*model.VersionHistory
		mockError   error
		expectLen   int
		expectError string
	}

	now := mockTime()

	testCases := []testCase{
		{
			name:        "empty list — no versions uploaded yet",
			mockReturn:  []*model.VersionHistory{},
			mockError:   nil,
			expectLen:   0,
			expectError: "",
		},
		{
			name: "single record",
			mockReturn: []*model.VersionHistory{
				{Id: 1, Version: "1.0.0", Code: "100", Url: "http://a.apk", Feature: "init", Force: false, Type: "release", CreatedAt: now},
			},
			mockError:   nil,
			expectLen:   1,
			expectError: "",
		},
		{
			name: "multiple records — ordered by created_at desc",
			mockReturn: []*model.VersionHistory{
				{Id: 3, Version: "3.0.0", Code: "300", Type: "release", CreatedAt: now.Add(1)},
				{Id: 2, Version: "2.0.0", Code: "200", Type: "beta", CreatedAt: now},
			},
			mockError:   nil,
			expectLen:   2,
			expectError: "",
		},
		{
			name:        "database error",
			mockReturn:  nil,
			mockError:   fmt.Errorf("connection refused"),
			expectLen:   0,
			expectError: "get version history list error",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			cs := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			mockey.Mock((*dbversion.DBVersion).GetVersionHistoryList).Return(tc.mockReturn, tc.mockError).Build()

			svc := NewVersionService(context.Background(), cs)
			res, err := svc.GetVersionHistoryList()
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tc.expectLen)
				if tc.expectLen > 0 && res != nil {
					assert.Equal(t, tc.mockReturn[0].Version, res[0].Version)
					assert.Equal(t, tc.mockReturn[0].CreatedAt.Format("2006-01-02 15:04:05"), res[0].CreatedAt)
				}
			}
		})
	}
}

func TestGetLatestVersion_WithCacheHit(t *testing.T) {
	defer mockey.UnPatchAll()

	now := mockTime()
	cached := &model.VersionHistory{
		Id: 1, Version: "2.0.0", Code: "200", Url: "http://a.apk",
		Feature: "cached version", Force: false, Type: "release", CreatedAt: now,
	}

	mockey.PatchConvey("cache hit — returns cached version without DB query", t, func() {
		cs := &base.ClientSet{
			DBClient:    new(db.Database),
			CacheClient: new(cache.Cache),
		}

		mockey.Mock((*versioncache.CacheVersion).GetLatestVersionCache).Return(cached, nil).Build()

		svc := NewVersionService(context.Background(), cs)
		res, err := svc.GetReleaseVersion()
		assert.NoError(t, err)
		assert.Equal(t, "2.0.0", res.Version)
		assert.Equal(t, "cached version", res.Feature)
	})
}

func TestGetLatestVersion_WithCacheMiss(t *testing.T) {
	defer mockey.UnPatchAll()

	now := mockTime()
	dbRecord := &model.VersionHistory{
		Id: 2, Version: "3.0.0", Code: "300", Url: "http://b.apk",
		Feature: "from db", Force: true, Type: "beta", CreatedAt: now,
	}

	mockey.PatchConvey("cache miss — falls back to DB and populates cache", t, func() {
		cs := &base.ClientSet{
			DBClient:    new(db.Database),
			CacheClient: new(cache.Cache),
		}

		mockey.Mock((*versioncache.CacheVersion).GetLatestVersionCache).Return(nil, fmt.Errorf("redis: nil")).Build()
		mockey.Mock((*dbversion.DBVersion).GetLatestVersionByType).Return(dbRecord, nil).Build()
		mockey.Mock((*versioncache.CacheVersion).SetLatestVersionCache).Return(nil).Build()

		svc := NewVersionService(context.Background(), cs)
		res, err := svc.GetBetaVersion()
		assert.NoError(t, err)
		assert.Equal(t, "3.0.0", res.Version)
		assert.Equal(t, "from db", res.Feature)
		assert.True(t, res.Force)
	})
}

func TestGetLatestVersion_NoRecord(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("no version record exists — returns error", t, func() {
		cs := &base.ClientSet{
			DBClient:    new(db.Database),
			CacheClient: new(cache.Cache),
		}

		mockey.Mock((*versioncache.CacheVersion).GetLatestVersionCache).Return(nil, fmt.Errorf("redis: nil")).Build()
		mockey.Mock((*dbversion.DBVersion).GetLatestVersionByType).Return(nil, nil).Build()

		svc := NewVersionService(context.Background(), cs)
		_, err := svc.GetReleaseVersion()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no release version found")
	})
}

func TestGetLatestVersion_DBError(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("database error — returns error", t, func() {
		cs := &base.ClientSet{
			DBClient:    new(db.Database),
			CacheClient: new(cache.Cache),
		}

		mockey.Mock((*versioncache.CacheVersion).GetLatestVersionCache).Return(nil, fmt.Errorf("redis: nil")).Build()
		mockey.Mock((*dbversion.DBVersion).GetLatestVersionByType).Return(nil, fmt.Errorf("db down")).Build()

		svc := NewVersionService(context.Background(), cs)
		_, err := svc.GetBetaVersion()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "db query error")
	})
}
