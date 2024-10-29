package db

import (
	"context"
	"fmt"
)

func CreateImage(ctx context.Context, pictureModel *Picture) (*Picture, error) {
	if err := DB.WithContext(ctx).Create(pictureModel).Error; err != nil {
		return nil, fmt.Errorf("dal.CreateImage error: %v", err)
	}
	return pictureModel, nil
}
