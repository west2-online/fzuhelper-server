package dal

import "github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"

func Init() {
	db.InitMySQL()
}
