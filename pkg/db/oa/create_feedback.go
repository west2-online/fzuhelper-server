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

package oa

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (c *DBOA) CreateFeedback(ctx context.Context, fbModel *model.Feedback) error {
	if err := c.client.WithContext(ctx).Table(constants.FeedbackTableName).Create(fbModel).Error; err != nil {
		logger.Errorf("dal.CreateFeedback error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateFeedback error: %v", err)
	}
	return nil
}
