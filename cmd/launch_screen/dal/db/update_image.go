package db

import (
	"context"
	"fmt"
)

func UpdateImage(ctx context.Context, pictureModel *Picture) (*Picture, error) {
	if err := DB.WithContext(ctx).Save(pictureModel).Take(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.UpdateImage error: %v", err)
	}
	return pictureModel, nil
}
