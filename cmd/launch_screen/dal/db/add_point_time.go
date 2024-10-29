package db

import (
	"context"
	"fmt"
)

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
