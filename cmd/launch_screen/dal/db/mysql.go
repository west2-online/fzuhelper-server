package db

import (
	"context"
	"gorm.io/gorm"
	"time"
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
