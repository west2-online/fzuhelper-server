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
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// GetVersionHistoryList returns all version history records ordered by created_at descending.
// Returns an empty slice (not nil) when no records exist.
func (c *DBVersion) GetVersionHistoryList(ctx context.Context, limit int, pageToken int64) ([]*model.VersionHistory, int64, error) {
	if limit <= 0 {
		limit = constants.VersionHistoryDefaultPageSize
	}

	tx := c.client.WithContext(ctx).Table(constants.VersionHistoryTableName)
	if pageToken > 0 {
		tx = tx.Where("id < ?", pageToken)
	}

	var history []*model.VersionHistory
	err := tx.
		Order("created_at desc, id desc").
		Limit(limit + 1).
		Find(&history).Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetVersionHistoryList error: %v", err)
	}

	nextPageToken := int64(0)
	if len(history) > limit {
		nextPageToken = history[limit-1].Id
		history = history[:limit]
	}

	return history, nextPageToken, nil
}

// GetLatestVersionByType returns the most recent version record for the given type (release/beta).
// Returns nil, nil when no record exists for the type.
func (c *DBVersion) GetLatestVersionByType(ctx context.Context, versionType string) (*model.VersionHistory, error) {
	vh := new(model.VersionHistory)
	err := c.client.WithContext(ctx).Table(constants.VersionHistoryTableName).
		Where("type = ?", versionType).
		Order("created_at desc, id desc").
		First(vh).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetLatestVersionByType error: %v", err)
	}
	return vh, nil
}
