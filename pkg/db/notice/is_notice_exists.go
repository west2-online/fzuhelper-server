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

package notice

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// IsNoticeExists 根据 title 和 url 做为唯一索引
func (d *DBNotice) IsNoticeExists(ctx context.Context, title string, url string) (ok bool, err error) {
	var count int64
	err = d.client.WithContext(ctx).Table(constants.NoticeTableName).Where("title = ? AND url = ?", title).Count(&count).Error
	if err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.IsTitleExists error: %s", err)
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}
