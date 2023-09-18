package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Picture struct {
	PictureId  int64     `gorm:"primarykey"`
	Url        string    `json:"url"`
	Href       string    `json:"href"`
	Text       string    `json:"text"`
	PicType    int8      `json:"pic_type" gorm:"default:1"`
	ShowTimes  int64     `json:"show_times" gorm:"default:0"`
	PointTimes int64     `json:"point_times" gorm:"default:0"`
	Duration   int64     `json:"duration" gorm:"default:3"`
	StartAt    time.Time `json:"start_at"`                    // 开始时间
	EndAt      time.Time `json:"end_at"`                      // 结束时间
	StartTime  int64     `json:"start_time" gorm:"default:0"` // 开始时段 0~24
	EndTime    int64     `json:"end_time" gorm:"default:24"`  // 结束时段 0~24
	SType      int8      `json:"s_type"`                      // 类型
	Frequency  int64     `json:"frequency"`                   // 一天展示次数
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func SavePicture(picture *Picture) error {
	return DB.Save(&picture).Error
}

func GetPictureById(picture_id int64) (*Picture, error) {
	picture := new(Picture)
	err := DB.Where("picture_id=?", picture_id).First(&picture).Error
	if err != nil {
		return nil, err
	}
	return picture, nil
}
func GetPictures() ([]*Picture, error) {
	pics := make([]*Picture, 0)
	err := DB.Find(&pics).Error
	return pics, err
}

func GetRetPicture(stype int8) ([]*Picture, error) {
	now := time.Now()
	hour := strings.Split(strings.Split(now.Format(time.DateTime), " ")[1], ":")[0]
	pics := make([]*Picture, 0)
	err := DB.Where("start_at < ?", now).
		Where("end_at > ?", now).
		Where("s_type = ?", stype).
		Where("start_time <= ?", hour).
		Where("end_time >= ?", hour).Find(&pics).Error
	if err != nil {
		return nil, err
	}
	return pics, nil
}

func DeletePicture(picture_id int64) (*Picture, error) {
	pic := new(Picture)
	err := DB.Where("picture_id=?", picture_id).Delete(&pic).Error
	return pic, err
}
