package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"time"
)

func (s *LaunchScreenService) PutImage(picture *model.Picture) (*db.Picture, error) {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	pictureModel := &db.Picture{
		ID:         picture.Id,
		Uid:        picture.UserId,
		Url:        picture.Url,
		Href:       picture.Href,
		Text:       picture.Text,
		PicType:    picture.PicType,
		ShowTimes:  0,
		PointTimes: 0,
		Duration:   picture.Duration,
		SType:      *picture.SType,
		Frequency:  picture.Frequency,
		StartTime:  picture.StartTime,
		EndTime:    picture.EndTime,
		StartAt:    time.Unix(picture.StartAt, 0).In(Loc),
		EndAt:      time.Unix(picture.EndAt, 0).In(Loc),
	}
	return db.CreateImage(s.ctx, pictureModel)
}
