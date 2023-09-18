package service

import (
	"errors"
	"strconv"

	"github.com/ozline/tiktok/cmd/screen/dal/db"
	"github.com/ozline/tiktok/kitex_gen/screen"
)

func (s *ScreenService) RetPicture(req *screen.RetPictureRequest) ([]*db.Picture, error) {
	s_type, err := strconv.Atoi(req.Type)
	if err != nil {
		return nil, errors.New("type error")
	}
	imgs, err := db.GetRetPicture(int8(s_type))
	if err != nil {
		return nil, err
	}
	return imgs, nil
}
