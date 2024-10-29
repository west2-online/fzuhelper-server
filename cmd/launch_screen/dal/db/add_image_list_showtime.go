package db

import (
	"context"
	"fmt"
)

func AddImageListShowTime(ctx context.Context, pictureList *[]Picture) error {
	for i := range *pictureList {
		(*pictureList)[i].ShowTimes++
	}
	if err := DB.WithContext(ctx).Save(pictureList).Error; err != nil {
		return fmt.Errorf("dal.AddImageListShowTime error: %v", err)
	}
	return nil
}
