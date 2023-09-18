package service

import "github.com/west2-online/fzuhelper-server/cmd/screen/dal/db"

func (s *ScreenService) AddPoint(picture_id int64) error {
	picture, err := db.GetPictureById(picture_id)
	if err != nil {
		return err
	}

	err = db.SavePicture(picture)
	if err != nil {
		return err
	}
	return nil
}
