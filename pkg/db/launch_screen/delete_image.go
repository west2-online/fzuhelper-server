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

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func (c *DBLaunchScreen) DeleteImage(ctx context.Context, id int64) (*model.Picture, error) {
	pictureModel := &model.Picture{
		ID: id,
	}
	if err := c.client.WithContext(ctx).Table(constants.LaunchScreenTableName).Take(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.DeleteImage error: %w", err)
	}

	if err := c.client.WithContext(ctx).Table(constants.LaunchScreenTableName).Delete(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.DeleteImage error: %w", err)
	}
	return pictureModel, nil
}
