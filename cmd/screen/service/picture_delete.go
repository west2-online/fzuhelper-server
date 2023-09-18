package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/screen"
)

func (s *ScreenService) DeletePicture(req *screen.DeletePictureRequest) (*db.Picture, error) {
	picture, err := db.DeletePicture(req.PictureId)
	if err != nil {
		return nil, err
	}
	return picture, nil
}
