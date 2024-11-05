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

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func (c *DBLaunchScreen) AddPointTime(ctx context.Context, id int64) error {
	pictureModel := new(model.Picture)
	if err := c.client.WithContext(ctx).Where("id = ?", id).First(pictureModel).Error; err != nil {
		return fmt.Errorf("dal.AddPointTime error: %v", err)
	}
	pictureModel.PointTimes++
	if err := c.client.WithContext(ctx).Save(pictureModel).Error; err != nil {
		return fmt.Errorf("dal.AddPointTime error: %v", err)
	}
	return nil
}
