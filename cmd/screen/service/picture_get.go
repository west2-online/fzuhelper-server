package service

import "github.com/ozline/tiktok/cmd/screen/dal/db"

func (s *ScreenService) GetPicture(picture_id int64) ([]*db.Picture, error) {
	pics := make([]*db.Picture, 0)
	if picture_id != 0 {
		picture, err := db.GetPictureById(picture_id)
		if err != nil {
			return nil, err
		}
		pics = append(pics, picture)
	} else {
		picture, err := db.GetPictures()
		if err != nil {
			return nil, err
		}
		pics = append(pics, picture...)
	}
	return pics, nil
}
