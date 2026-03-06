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

package academic

import (
	"context"
	"fmt"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// GetScoresBatchByStuId 按 stu_id 升序分批获取 scores 表数据
func (c *DBAcademic) GetScoresBatchByStuId(ctx context.Context, lastStuId string, batchSize int) ([]*model.Score, error) {
	var scores []*model.Score
	if err := c.client.WithContext(ctx).
		Table(constants.ScoreTableName).
		Where("stu_id > ?", lastStuId).
		Order("stu_id ASC").
		Limit(batchSize).
		Find(&scores).Error; err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.GetScoresBatchByStuId error: %v", err))
	}
	return scores, nil
}
