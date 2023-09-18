package service

import (
	"errors"
	"io"
	"time"

	"github.com/ozline/tiktok/cmd/screen/dal/db"
	"github.com/ozline/tiktok/kitex_gen/screen"
	"github.com/ozline/tiktok/pkg/utils"
)

func (s *ScreenService) UpdatePicture(req *screen.PutPictureRequset) (*db.Picture, error) {
	// update
	if req.StartTime > req.EndTime || req.StartAt > req.EndAt {
		return nil, errors.New("time error")
	}
	picture, err := db.GetPictureById(req.PictureId)
	if err != nil {
		return nil, err
	}
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	picture.Href = req.Href
	picture.PicType = req.PicType
	picture.Duration = req.Duration
	picture.StartAt = time.Unix(int64(req.StartAt), 0).In(Loc)
	picture.EndAt = time.Unix(int64(req.EndAt), 0).In(Loc)
	picture.StartTime = req.StartTime
	picture.EndTime = req.EndTime
	picture.SType = req.SType
	picture.Text = req.Text
	err = db.SavePicture(picture)
	if err != nil {
		return nil, err
	}
	return picture, nil
}

func (s *ScreenService) UpdatePictureImg(req *screen.PutPictureImgRequset, img io.Reader) (*db.Picture, error) {
	// update img

	picture, err := db.GetPictureById(req.PictureId)
	if err != nil {
		return nil, err
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
