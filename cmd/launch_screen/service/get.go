package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *LaunchScreenService) GetImageById(id int64, uid int64) (*db.Picture, error) {
	img, err := db.GetImageById(s.ctx, id)
	if err != nil {
		return nil, err
	}
	if img.Uid != uid {
		return nil, errno.NoAccessError
	}
	return img, nil
}

func (s *LaunchScreenService) GetImagesByUid(uid int64) (*[]db.Picture, int64, error) {
	imgList, cnt, err := db.ListImageByUid(s.ctx, 1, uid)
	if err != nil {
		return nil, 0, err
	}
	return imgList, cnt, nil
}
