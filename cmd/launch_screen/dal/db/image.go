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

package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Picture struct {
	ID         int64
	Url        string
	Href       string
	Text       string
	PicType    int64
	ShowTimes  int64
	PointTimes int64
	Duration   int64
	StartAt    time.Time // 开始时间
	EndAt      time.Time // 结束时间
	StartTime  int64     // 开始时段 0~24
	EndTime    int64     // 结束时段 0~24
	SType      int64     // 类型
	Frequency  int64     // 一天展示次数
	Regex      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `sql:"index"`
}

func CreateImage(ctx context.Context, pictureModel *Picture) (*Picture, error) {
	if err := DB.WithContext(ctx).Create(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.CreateImage error: %v", err)
	}
	return pictureModel, nil
}

func GetImageById(ctx context.Context, id int64) (*Picture, error) {
	pictureModel := new(Picture)
	if err := DB.WithContext(ctx).Where("id = ?", id).First(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.GetImageById error: %v", err)
	}
	return pictureModel, nil
}

func DeleteImage(ctx context.Context, id int64) (*Picture, error) {
	pictureModel := &Picture{
		ID: id,
	}
	if err := DB.WithContext(ctx).Take(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.DeleteImage error: %v", err)
	}

	if err := DB.WithContext(ctx).Delete(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.DeleteImage error: %v", err)
	}
	return pictureModel, nil
}

func AddPointTime(ctx context.Context, id int64) error {
	pictureModel := new(Picture)
	if err := DB.WithContext(ctx).Where("id = ?", id).First(pictureModel).Error; err != nil {
		return fmt.Errorf("dal.AddPointTime error: %v", err)
	}
	pictureModel.PointTimes++
	if err := DB.WithContext(ctx).Save(pictureModel).Error; err != nil {
		return fmt.Errorf("dal.AddPointTime error: %v", err)
	}
	return nil
}

func UpdateImage(ctx context.Context, pictureModel *Picture) (*Picture, error) {
	if err := DB.WithContext(ctx).Save(pictureModel).Take(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.UpdateImage error: %v", err)
	}
	return pictureModel, nil
}

func GetImageBySType(ctx context.Context, sType int64) (*[]Picture, int64, error) {
	pictures := new([]Picture)
	var count int64 = 0
	now := time.Now().Add(time.Hour * 8)
	hour := now.Hour() + 8
	if hour > 24 {
		hour -= 24
	}
	// 按创建时间降序
	if err := DB.WithContext(ctx).
		Where("s_type = ? AND start_at < ? AND end_at > ? AND start_time <= ? AND end_time >= ?",
			sType, now, now, hour, hour).
		Count(&count).Order("created_at DESC").
		Find(pictures).
		Error; err != nil {
		return nil, -1, err
	}
	return pictures, count, nil
}

func GetLastImageId(ctx context.Context) (int64, error) {
	pictureModel := new(Picture)
	if err := DB.WithContext(ctx).Last(pictureModel).Error; err != nil {
		return -1, fmt.Errorf("dal.GetLastImageId error: %v", err)
	}
	return pictureModel.ID, nil
}

func GetImageByIdList(ctx context.Context, imgIdList *[]int64) (*[]Picture, int64, error) {
	pictures := new([]Picture)
	var count int64 = 0
	now := time.Now().Add(time.Hour * 8)
	hour := now.Hour() + 8
	if hour > 24 {
		hour -= 24
	}
	err := DB.WithContext(ctx).
		Where("id IN ? AND start_at < ? AND end_at > ? AND start_time <= ? AND end_time >= ?",
			*imgIdList, now, now, hour, hour).Count(&count).Order("created_at DESC").Find(pictures).Error
	if err != nil {
		return nil, -1, fmt.Errorf("dal.GetImageByIdList error: %w", err)
	}
	return pictures, count, nil
}

func AddImageListShowTime(ctx context.Context, pictureList *[]Picture) error {
	for i := range *pictureList {
		(*pictureList)[i].ShowTimes++
	}
	if err := DB.WithContext(ctx).Save(pictureList).Error; err != nil {
		return fmt.Errorf("dal.AddImageListShowTime error: %v", err)
	}
	return nil
}
