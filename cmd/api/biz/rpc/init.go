package rpc

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
)

var (
	classroomClient classroomservice.Client
	userClient      userservice.Client
)

func Init() {
	InitClassroomRPC()
	InitUserRPC()
}
