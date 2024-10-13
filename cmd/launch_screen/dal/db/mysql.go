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
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type Picture struct {
	ID         int64
	Uid        int64
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
	StudentId  int64
	DeviceType int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `sql:"index"`
}

func CreateImage(ctx context.Context, pictureModel *Picture) (*Picture, error) {
	if err := DB.WithContext(ctx).Create(pictureModel).Error; err != nil {
		return nil, err
	}
	return pictureModel, nil
}

func GetImageById(ctx context.Context, id int64) (*Picture, error) {
	pictureModel := new(Picture)
	if err := DB.WithContext(ctx).Where("id = ?", id).First(pictureModel).Error; err != nil {
		return nil, err
	}
	return pictureModel, nil
}

func ListImageByUid(ctx context.Context, pageNum int, uid int64) (*[]Picture, int64, error) {
	pictures := new([]Picture)
	var count int64 = 0
	// 按创建时间降序
	// Offset((1 - 1) * constants.PageSize).
	if err := DB.WithContext(ctx).Where("uid = ?", uid).Count(&count).Order("created_at DESC").
		Limit(constants.PageSize).
		Find(pictures).
		Error; err != nil {
		return nil, 114514, err
	}
	return pictures, count, nil
}

func DeleteImage(ctx context.Context, id int64, uid int64) (*Picture, error) {
	pictureModel := &Picture{
		ID: id,
	}
	if err := DB.WithContext(ctx).Take(pictureModel).Error; err != nil {
		return nil, err
	}
	if pictureModel.Uid != uid {
		return nil, errno.NoAccessError
	}

	if err := DB.WithContext(ctx).Delete(pictureModel).Error; err != nil {
		return nil, err
	}
	return pictureModel, nil
}

func AddPointTime(ctx context.Context, id int64) error {
	pictureModel := new(Picture)
	if err := DB.WithContext(ctx).Where("id = ?", id).First(pictureModel).Error; err != nil {
		return err
	}
	pictureModel.ShowTimes++
	if err := DB.WithContext(ctx).Save(pictureModel).Error; err != nil {
		return err
	}
	return nil
}

func UpdateImage(ctx context.Context, pictureModel *Picture) (*Picture, error) {
	if err := DB.WithContext(ctx).Save(pictureModel).Take(pictureModel).Error; err != nil {
		return nil, err
	}
	return pictureModel, nil
}

func GetImageByStuId(ctx context.Context, pictureModel *Picture) (*[]Picture, int64, error) {
	pictures := new([]Picture)
	var count int64 = 0
	now := time.Now().Add(time.Hour * 8)
	hour := strings.Split(strings.Split(now.Format("2006-01-02 15:04:05"), " ")[1], ":")[0]
	// 按创建时间降序
	// Offset((1 - 1) * constants.PageSize).
	if err := DB.WithContext(ctx).
		Where("student_id = ? AND device_type = ? AND s_type = ? AND start_at < ? AND end_at > ? AND start_time <= ? AND end_time >= ?",
			pictureModel.StudentId, pictureModel.DeviceType, pictureModel.SType, now, now, hour, hour).
		Count(&count).Order("created_at DESC").
		Limit(constants.PageSize).
		Find(pictures).
		Error; err != nil {
		return nil, 114514, err
	}
	return pictures, count, nil
}
