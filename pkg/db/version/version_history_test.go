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
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func makeVersionHistory(id int64, version, versionType string, createdAt time.Time) *model.VersionHistory {
	return &model.VersionHistory{
		Id:        id,
		Version:   version,
		Code:      fmt.Sprintf("%d", id),
		Url:       "https://example.com/app.apk",
		Feature:   "feature",
		Force:     false,
		Type:      versionType,
		CreatedAt: createdAt,
	}
}

func TestDBVersion_GetVersionHistoryList(t *testing.T) {
	defer mockey.UnPatchAll()

	now := time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name            string
		limit           int
		pageToken       int64
		mockRows        []*model.VersionHistory
		mockError       error
		expectWhere     bool
		expectLimit     int
		expectLen       int
		expectNextToken int64
		expectError     string
		expectFirstID   int64
		expectPageToken int64
	}{
		{
			name:        "empty list with default limit",
			limit:       0,
			mockRows:    []*model.VersionHistory{},
			expectLimit: constants.VersionHistoryDefaultPageSize + 1,
			expectLen:   0,
		},
		{
			name:      "returns page and next token",
			limit:     2,
			pageToken: 10,
			mockRows: []*model.VersionHistory{
				makeVersionHistory(9, "3.0.0", "release", now.Add(2*time.Hour)),
				makeVersionHistory(8, "2.0.0", "beta", now.Add(time.Hour)),
				makeVersionHistory(7, "1.0.0", "release", now),
			},
			expectWhere:     true,
			expectLimit:     3,
			expectLen:       2,
			expectNextToken: 8,
			expectFirstID:   9,
			expectPageToken: 10,
		},
		{
			name:        "database error",
			limit:       3,
			mockError:   gorm.ErrInvalidDB,
			expectLimit: 4,
			expectError: "dal.GetVersionHistoryList error",
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			dbVersion := NewDBVersion(mockGormDB, nil)
			whereCalled := false

			mockey.Mock((*gorm.DB).WithContext).To(func(_ *gorm.DB, ctx context.Context) *gorm.DB {
				assert.NotNil(t, ctx)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(_ *gorm.DB, name string, args ...interface{}) *gorm.DB {
				assert.Equal(t, constants.VersionHistoryTableName, name)
				assert.Empty(t, args)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(_ *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
				whereCalled = true
				assert.Equal(t, "id < ?", query)
				assert.Equal(t, []interface{}{tc.expectPageToken}, args)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Order).To(func(_ *gorm.DB, value interface{}) *gorm.DB {
				assert.Equal(t, "created_at desc, id desc", value)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Limit).To(func(_ *gorm.DB, limit int) *gorm.DB {
				assert.Equal(t, tc.expectLimit, limit)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Find).To(func(_ *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				assert.Empty(t, conds)
				if tc.mockError != nil {
					return &gorm.DB{Error: tc.mockError}
				}
				history, ok := dest.(*[]*model.VersionHistory)
				assert.True(t, ok)
				*history = tc.mockRows
				return &gorm.DB{}
			}).Build()

			res, nextToken, err := dbVersion.GetVersionHistoryList(context.Background(), tc.limit, tc.pageToken)
			assert.Equal(t, tc.expectWhere, whereCalled)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, res)
				assert.Zero(t, nextToken)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, res, tc.expectLen)
			assert.Equal(t, tc.expectNextToken, nextToken)
			if tc.expectFirstID > 0 {
				assert.Equal(t, tc.expectFirstID, res[0].Id)
			}
		})
	}
}

func TestDBVersion_GetLatestVersionByType(t *testing.T) {
	defer mockey.UnPatchAll()

	now := time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)
	latest := makeVersionHistory(2, "2.0.0", "release", now)

	testCases := []struct {
		name        string
		mockRecord  *model.VersionHistory
		mockError   error
		expectNil   bool
		expectError string
	}{
		{
			name:       "success",
			mockRecord: latest,
		},
		{
			name:      "record not found",
			mockError: gorm.ErrRecordNotFound,
			expectNil: true,
		},
		{
			name:        "database error",
			mockError:   gorm.ErrInvalidDB,
			expectNil:   true,
			expectError: "dal.GetLatestVersionByType error",
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			dbVersion := NewDBVersion(mockGormDB, nil)

			mockey.Mock((*gorm.DB).WithContext).To(func(_ *gorm.DB, ctx context.Context) *gorm.DB {
				assert.NotNil(t, ctx)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(_ *gorm.DB, name string, args ...interface{}) *gorm.DB {
				assert.Equal(t, constants.VersionHistoryTableName, name)
				assert.Empty(t, args)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(_ *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
				assert.Equal(t, "type = ?", query)
				assert.Equal(t, []interface{}{"release"}, args)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Order).To(func(_ *gorm.DB, value interface{}) *gorm.DB {
				assert.Equal(t, "created_at desc, id desc", value)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(_ *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				assert.Empty(t, conds)
				if tc.mockError != nil {
					return &gorm.DB{Error: tc.mockError}
				}
				versionHistory, ok := dest.(*model.VersionHistory)
				assert.True(t, ok)
				*versionHistory = *tc.mockRecord
				return &gorm.DB{}
			}).Build()

			res, err := dbVersion.GetLatestVersionByType(context.Background(), "release")
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			if tc.expectNil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, latest.Version, res.Version)
				assert.Equal(t, latest.Type, res.Type)
			}
		})
	}
}

func TestDBVersion_CreateVersionHistory(t *testing.T) {
	defer mockey.UnPatchAll()

	versionHistory := makeVersionHistory(1, "1.0.0", "release", time.Now())

	testCases := []struct {
		name        string
		mockError   error
		expectError string
	}{
		{
			name: "success",
		},
		{
			name:        "database error",
			mockError:   gorm.ErrInvalidDB,
			expectError: "dal.CreateVersionHistory error",
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			dbVersion := NewDBVersion(mockGormDB, nil)

			mockey.Mock((*gorm.DB).WithContext).To(func(_ *gorm.DB, ctx context.Context) *gorm.DB {
				assert.NotNil(t, ctx)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(_ *gorm.DB, name string, args ...interface{}) *gorm.DB {
				assert.Equal(t, constants.VersionHistoryTableName, name)
				assert.Empty(t, args)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Create).To(func(_ *gorm.DB, value interface{}) *gorm.DB {
				assert.Equal(t, versionHistory, value)
				return &gorm.DB{Error: tc.mockError}
			}).Build()

			err := dbVersion.CreateVersionHistory(context.Background(), versionHistory)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
