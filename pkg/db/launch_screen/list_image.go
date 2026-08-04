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

package launch_screen

import (
	"context"
	"fmt"
	"math"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// ListImage returns one page of launch screen pictures ordered by id descending.
func (c *DBLaunchScreen) ListImage(ctx context.Context, pageNum, pageSize int) (*[]model.Picture, int64, error) {
	if pageNum <= 0 || pageSize <= 0 {
		return nil, 0, errno.NewErrNo(errno.ParamErrorCode, "page_num and page_size must be positive")
	}
	if pageNum-1 > math.MaxInt/pageSize {
		return nil, 0, errno.NewErrNo(errno.ParamErrorCode, "page offset is too large")
	}

	var total int64
	if err := c.client.WithContext(ctx).
		Table(constants.LaunchScreenTableName).
		Count(&total).Error; err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.ListImage count error: %v", err))
	}

	pictures := new([]model.Picture)
	offset := (pageNum - 1) * pageSize
	if err := c.client.WithContext(ctx).
		Table(constants.LaunchScreenTableName).
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(pictures).Error; err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.ListImage error: %v", err))
	}

	return pictures, total, nil
}
