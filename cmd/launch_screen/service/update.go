package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"time"
)

func (s *LaunchScreenService) UpdateImageProperty(req *launch_screen.ChangeImagePropertyRequest, uid int64) (*db.Picture, error) {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	pictureModel := &db.Picture{
		ID:        req.PictureId,
		Uid:       uid,
		PicType:   req.PicType,
		Duration:  *req.Duration,
		Href:      *req.Href,
		SType:     req.SType,
		Frequency: req.Frequency,
		Text:      req.Text,
		StartAt:   time.Unix(req.StartAt, 0).In(Loc),
		EndAt:     time.Unix(req.EndAt, 0).In(Loc),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	return db.UpdateImage(s.ctx, pictureModel)
}

func (s *LaunchScreenService) UpdateImagePath(id int64, url string) (*db.Picture, error) {
	pictureModel := &db.Picture{
		ID:  id,
		Url: url,
	}
	return db.UpdateImage(s.ctx, pictureModel)
}
