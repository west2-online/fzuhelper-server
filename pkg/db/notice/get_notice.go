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
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (d *DBNotice) GetNoticeByPage(ctx context.Context, pageNum int) (list []model.Notice, err error) {
	// 不使用[]*的原因：Find 返回多个结果时，只能使用[]
	offset := (pageNum - 1) * constants.NoticePageSize
	if err := d.client.WithContext(ctx).
		Table(constants.NoticeTableName).
		Order("published_at DESC").
		Limit(constants.NoticePageSize).Offset(offset).
		Find(&list).
		Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetNoticeByPage error: %s", err)
	}
	return list, nil
}
