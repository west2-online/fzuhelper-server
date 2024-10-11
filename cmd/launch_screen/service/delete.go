package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
)

func (s *LaunchScreenService) DeleteImage(id int64, uid int64) (*db.Picture, error) {
	return db.DeleteImage(s.ctx, id, uid)
}
