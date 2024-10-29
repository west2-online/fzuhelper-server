package db

import (
	"context"
	"fmt"
	"time"
)

func GetImageById(ctx context.Context, id int64) (*Picture, error) {
	pictureModel := new(Picture)
	if err := DB.WithContext(ctx).Where("id = ?", id).First(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.GetImageById error: %v", err)
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

func GetLastImageId(ctx context.Context) (int64, error) {
	pictureModel := new(Picture)
	if err := DB.WithContext(ctx).Last(pictureModel).Error; err != nil {
		return -1, fmt.Errorf("dal.GetLastImageId error: %v", err)
	}
	return pictureModel.ID, nil
}
