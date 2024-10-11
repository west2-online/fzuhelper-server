package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
)

func (s *LaunchScreenService) AddPointTime(id int64) error {
	return db.AddPointTime(s.ctx, id)
}
