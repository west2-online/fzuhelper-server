package rpc

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
)

var (
	classroomClient classroomservice.Client
)

func Init() {
	InitClassroomRPC()
}
