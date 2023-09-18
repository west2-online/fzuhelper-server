package dal

import (
	"github.com/west2-online/fzuhelper-server/cmd/template/dal/cache"
	"github.com/west2-online/fzuhelper-server/cmd/template/dal/db"
	"github.com/west2-online/fzuhelper-server/cmd/template/dal/mq"
)

func Init() {
	db.Init()
	mq.Init()
	cache.Init()
}
