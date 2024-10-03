package dal

import (
	"github.com/west2-online/fzuhelper-server/cmd/user/dal/db"
)

func Init() {
	db.InitMySQL()
}
