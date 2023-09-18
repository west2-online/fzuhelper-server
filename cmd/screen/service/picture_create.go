package service

import (
	"io"
	"time"

	"github.com/west2-online/fzuhelper-server/cmd/screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/screen"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *ScreenService) CreatePicture(req *screen.CreatePictureRequest, img io.Reader) (*db.Picture, error) {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	// 构造->db创建->返回
	pid, err := db.SF.NextVal()
	if err != nil {
		return nil, err
	}
	picture := &db.Picture{
		PictureId:  pid,
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
