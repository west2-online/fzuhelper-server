package rpc

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
)

var (
	userClient         userservice.Client
	launchScreenClient launchscreenservice.Client
)

func Init() {
	InitUserRPC()
	InitLaunchScreenRPC()
}
