package service

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

type LaunchScreenService struct {
	ctx context.Context
}

func NewLaunchScreenService(ctx context.Context) *LaunchScreenService {
	return &LaunchScreenService{ctx: ctx}
}

func BuildImageResp(dbP *db.Picture) *model.Picture {
	return &model.Picture{
		Id:         dbP.ID,
		Url:        dbP.Url,
		Href:       dbP.Href,
		Text:       dbP.Text,
		PicType:    dbP.PicType,
		ShowTimes:  &dbP.ShowTimes,
		PointTimes: &dbP.PointTimes,
		Duration:   dbP.Duration,
		SType:      &dbP.SType,
		Frequency:  dbP.Frequency,
	}
}
