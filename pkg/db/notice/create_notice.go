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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// CreateNotice 使用了类似 on conflict 的方式，当数据库中有对应 notice 时，只做更新数据操作
func (d *DBNotice) CreateNotice(ctx context.Context, notice *model.Notice) error {
	id, err := d.sf.NextVal()
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateNotice: NextVal error: %s", err)
	}
	notice.Id = id

	err = d.client.WithContext(ctx).
		Table(constants.NoticeTableName).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "title"},
				{Name: "url"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"published_at": notice.PublishedAt,
				"deleted_at":   notice.DeletedAt,
				"updated_at":   gorm.Expr("CURRENT_TIMESTAMP"),
			}),
		}).
		Create(notice).Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateNotice error: %s", err)
	}

	return nil
}
