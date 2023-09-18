package pack

import (
	"time"

	"github.com/west2-online/fzuhelper-server/cmd/screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/screen"
)

func BuildPicturesResp(pics []*db.Picture) []*screen.Picture {
	pictures := make([]*screen.Picture, 0)
	for _, pic := range pics {
		pictures = append(pictures, ConvertPicture(pic))
	}
	return pictures
}

func ConvertPicture(db_pic *db.Picture) *screen.Picture {
	return &screen.Picture{
		PictureId:  db_pic.PictureId,
		CreateAt:   db_pic.CreatedAt.Format(time.RFC3339),
		UpdateAt:   db_pic.UpdatedAt.Format(time.RFC3339),
		Url:        db_pic.Url,
		Herf:       db_pic.Href,
		Text:       db_pic.Text,
		PicType:    db_pic.PicType,
		SType:      db_pic.SType,
		ShowTimes:  db_pic.ShowTimes,
		PointTimes: db_pic.PointTimes,
		Duration:   db_pic.Duration,
		StartAt:    db_pic.StartAt.Format(time.RFC3339),
		EndAt:      db_pic.EndAt.Format(time.RFC3339),
		Frequency:  db_pic.Frequency,
	}
}
