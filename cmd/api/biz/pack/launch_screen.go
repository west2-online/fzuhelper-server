package pack

import (
	api "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildLaunchScreen(image *model.Picture) *api.Picture {
	return &api.Picture{
		ID:        image.Id,
		UserID:    image.UserId,
		URL:       image.Url,
		PicType:   image.PicType,
		Duration:  image.Duration,
		Href:      image.Href,
		ShowTimes: image.ShowTimes,
		SType:     image.SType,
		Frequency: image.Frequency,
		Text:      image.Text,
		StartAt:   image.StartAt,
		EndAt:     image.EndAt,
		StartTime: image.StartTime,
		EndTime:   image.EndTime,
	}
}
