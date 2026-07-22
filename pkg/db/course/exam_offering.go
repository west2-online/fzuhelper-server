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

package course

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBCourse) GetExamOfferingByHash(ctx context.Context, examHash string) (*model.ExamOffering, error) {
	// 查询全局考试变化是否已经被其他刷新任务占用。
	offering := new(model.ExamOffering)
	if err := c.client.WithContext(ctx).
		Table(constants.ExamOfferingsTableName).
		Where("exam_hash = ?", examHash).
		First(offering).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetExamOfferingByHash error: %v", err)
	}
	return offering, nil
}

func (c *DBCourse) CreateExamOffering(ctx context.Context, offering *model.ExamOffering) (*model.ExamOffering, error) {
	// 依靠 exam_hash 唯一索引原子抢占发送资格；重复键表示已经被全局去重。
	if err := c.client.WithContext(ctx).
		Table(constants.ExamOfferingsTableName).
		Create(offering).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, nil
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateExamOffering error: %v", err)
	}
	return offering, nil
}

func (c *DBCourse) DeleteExamOfferingByHash(ctx context.Context, examHash string) error {
	// 释放失败通知的全局去重记录，使用硬删除确保后续可以再次插入相同 hash。
	if err := c.client.WithContext(ctx).
		Table(constants.ExamOfferingsTableName).
		Where("exam_hash = ?", examHash).
		Unscoped().
		Delete(&model.ExamOffering{}).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.DeleteExamOfferingByHash error: %v", err)
	}
	return nil
}
