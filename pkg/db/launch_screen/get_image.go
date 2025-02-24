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
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (c *DBLaunchScreen) GetImageById(ctx context.Context, id int64) (*model.Picture, error) {
	pictureModel := new(model.Picture)
	if err := c.client.WithContext(ctx).Table(constants.LaunchScreenTableName).Where("id = ?", id).First(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.GetImageById error: %w", err)
	}
	return pictureModel, nil
}

func (c *DBLaunchScreen) GetImageBySType(ctx context.Context, sType int64) (*[]model.Picture, int64, error) {
	Loc := utils.LoadCNLocation()
	pictures := new([]model.Picture)
	var count int64 = 0
	now := time.Now().In(Loc)
	hour := now.Hour()
	// 按创建时间降序
	if err := c.client.WithContext(ctx).Table(constants.LaunchScreenTableName).
		Where("s_type = ? AND start_at < ? AND end_at > ? AND start_time <= ? AND end_time >= ?",
			sType, now, now, hour, hour).
		Count(&count).Order("created_at DESC").
		Find(pictures).
		Error; err != nil {
		return nil, -1, err
	}
	return pictures, count, nil
}

func (c *DBLaunchScreen) GetImageByIdList(ctx context.Context, imgIdList *[]int64) (*[]model.Picture, int64, error) {
	Loc := utils.LoadCNLocation()
	pictures := new([]model.Picture)
	var count int64 = 0
	now := time.Now().In(Loc)
	hour := now.Hour()
	err := c.client.WithContext(ctx).Table(constants.LaunchScreenTableName).
		Where("id IN ? AND start_at < ? AND end_at > ? AND start_time <= ? AND end_time >= ?",
			*imgIdList, now, now, hour, hour).Count(&count).Order("created_at DESC").Find(pictures).Error
	if err != nil {
		return nil, -1, fmt.Errorf("dal.GetImageByIdList error: %w", err)
	}
	return pictures, count, nil
}

func (c *DBLaunchScreen) GetLastImageId(ctx context.Context) (int64, error) {
	pictureModel := new(model.Picture)
	if err := c.client.WithContext(ctx).Table(constants.LaunchScreenTableName).Last(pictureModel).Error; err != nil {
		return -1, fmt.Errorf("dal.GetLastImageId error: %w", err)
	}
	return pictureModel.ID, nil
}
