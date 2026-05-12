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

	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	versioncache "github.com/west2-online/fzuhelper-server/pkg/cache/version"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	dbversion "github.com/west2-online/fzuhelper-server/pkg/db/version"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func mockTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", "2026-04-26 12:00:00")
	return t
}

func TestGetVersionHistoryList(t *testing.T) {
	type testCase struct {
		name            string
		request         *version.GetVersionHistoryListRequest
		mockCheckPwd    bool
		mockReturn      []*model.VersionHistory
		mockNextToken   int64
		mockError       error
		expectLimit     int
		expectPageToken int64
		expectNextToken int64
		expectLen       int
		expectError     string
	}

	now := mockTime()
	limit := int64(1)
	pageToken := int64(3)

	testCases := []testCase{
		{
			name:         "empty list — no versions uploaded yet",
			request:      &version.GetVersionHistoryListRequest{Password: "validpassword"},
			mockCheckPwd: true,
			mockReturn:   []*model.VersionHistory{},
			mockError:    nil,
			expectLimit:  constants.VersionHistoryDefaultPageSize,
			expectLen:    0,
			expectError:  "",
		},
		{
			name:         "single record",
			request:      &version.GetVersionHistoryListRequest{Password: "validpassword"},
			mockCheckPwd: true,
			mockReturn: []*model.VersionHistory{
				{Id: 1, Version: "1.0.0", Code: "100", Url: "http://a.apk", Feature: "init", Force: false, Type: "release", CreatedAt: now},
			},
			mockError:   nil,
			expectLimit: constants.VersionHistoryDefaultPageSize,
			expectLen:   1,
		},
		{
			name:         "multiple records — with page token",
			request:      &version.GetVersionHistoryListRequest{Password: "validpassword", Limit: &limit, PageToken: &pageToken},
			mockCheckPwd: true,
			mockReturn: []*model.VersionHistory{
				{Id: 3, Version: "3.0.0", Code: "300", Type: "release", CreatedAt: now.Add(1)},
				{Id: 2, Version: "2.0.0", Code: "200", Type: "beta", CreatedAt: now},
			},
			mockNextToken:   2,
			expectLimit:     1,
			expectPageToken: pageToken,
			expectNextToken: 2,
			expectLen:       2,
		},
		{
			name:         "database error",
			request:      &version.GetVersionHistoryListRequest{Password: "validpassword"},
			mockCheckPwd: true,
			mockReturn:   nil,
			mockError:    fmt.Errorf("connection refused"),
			expectLimit:  constants.VersionHistoryDefaultPageSize,
			expectLen:    0,
			expectError:  "get version history list error",
		},
		{
			name:         "invalid admin password",
			request:      &version.GetVersionHistoryListRequest{Password: "invalidpassword"},
			mockCheckPwd: false,
			expectError:  "authorization failed",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			cs := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			mockey.Mock(utils.CheckPwd).Return(tc.mockCheckPwd).Build()
			if tc.mockCheckPwd {
				mockey.Mock((*dbversion.DBVersion).GetVersionHistoryList).To(func(_ *dbversion.DBVersion, ctx context.Context, limit int, pageToken int64) ([]*model.VersionHistory, int64, error) {
					assert.Equal(t, tc.expectLimit, limit)
					assert.Equal(t, tc.expectPageToken, pageToken)
					return tc.mockReturn, tc.mockNextToken, tc.mockError
				}).Build()
			}

			svc := NewVersionService(context.Background(), cs)
			res, nextPageToken, err := svc.GetVersionHistoryList(tc.request)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectNextToken, nextPageToken)
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

	mockey.PatchConvey("cache miss — falls back to DB and queues cache population", t, func() {
		cs := &base.ClientSet{
			DBClient:    new(db.Database),
			CacheClient: new(cache.Cache),
		}

		mockey.Mock((*versioncache.CacheVersion).GetLatestVersionCache).Return(nil, fmt.Errorf("redis: nil")).Build()
		mockey.Mock((*dbversion.DBVersion).GetLatestVersionByType).Return(dbRecord, nil).Build()
		mockey.Mock((*versioncache.CacheVersion).SetLatestVersionCache).Return(nil).Build()
		addCalled := false
		mockey.Mock((*taskqueue.BaseTaskQueue).Add).To(func(_ *taskqueue.BaseTaskQueue, key string, task taskqueue.QueueTask) {
			addCalled = true
			assert.Equal(t, "setLatestVersionCache:beta", key)
		}).Build()

		svc := NewVersionService(context.Background(), cs)
		res, err := svc.GetBetaVersion()
		assert.NoError(t, err)
		assert.Equal(t, "3.0.0", res.Version)
		assert.Equal(t, "from db", res.Feature)
		assert.True(t, res.Force)
		assert.True(t, addCalled)
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
