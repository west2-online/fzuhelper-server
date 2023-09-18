package service

import (
	"github.com/ozline/tiktok/cmd/screen/dal/db"
	"github.com/ozline/tiktok/kitex_gen/screen"
)

func (s *ScreenService) DeletePicture(req *screen.DeletePictureRequest) (*db.Picture, error) {
	picture, err := db.DeletePicture(req.PictureId)
	if err != nil {
		return nil, err
	}
	return picture, nil
}
