package service

import (
	"io"
	"time"

	"github.com/ozline/tiktok/cmd/screen/dal/db"
	"github.com/ozline/tiktok/kitex_gen/screen"
	"github.com/ozline/tiktok/pkg/utils"
)

func (s *ScreenService) CreatePicture(req *screen.CreatePictureRequest, img io.Reader) (*db.Picture, error) {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	// 构造->db创建->返回
	picture := &db.Picture{
		PictureId:  db.SF.NextVal(),
		Href:       req.Href,
		Text:       req.Text,
		PicType:    req.PicType,
		ShowTimes:  0,
		PointTimes: 0,
		Duration:   req.Duration,
		StartAt:    time.Unix(int64(req.StartAt), 0).In(Loc),
		EndAt:      time.Unix(int64(req.EndAt), 0).In(Loc),
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		SType:      req.SType,
		Frequency:  req.Frequency,
	}
	url, err := utils.UploadImg(img)
	if err != nil {
		return nil, err
	}
	picture.Url = url
	err = db.SavePicture(picture)
	if err != nil {
		return nil, err
	}
	return picture, nil
}
