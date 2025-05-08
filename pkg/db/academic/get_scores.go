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
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBAcademic) GetScoreByStuId(ctx context.Context, stuId string) (*model.Score, error) {
	scoreModel := new(model.Score)
	if err := c.client.WithContext(ctx).Table(constants.ScoreTableName).Where("stu_id = ?", stuId).First(scoreModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.GetUserScoreByStuId error: %v", err))
	}
	return scoreModel, nil
}

func (c *DBAcademic) GetScoreSha256ByStuId(ctx context.Context, stuId string) (string, error) {
	scoreModel := new(model.Score)
	if err := c.client.WithContext(ctx).
		Table(constants.ScoreTableName).
		Select("stu_id", "scores_info_sha256").
		Where("stu_id = ?", stuId).
		First(scoreModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.GetUserScoreSha256ByStuId error: %v", err))
	}
	return scoreModel.ScoresInfoSHA256, nil
}
