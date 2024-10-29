package db

import (
	"context"
	"fmt"
)

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
