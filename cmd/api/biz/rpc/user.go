package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitUserRPC() {
	client, err := client.InitUserRPC()
	if err != nil {
		logger.Fatalf("api.rpc.user InitUserRPC failed, err is %v", err)
	}
	userClient = *client
}

func GetLoginDataRPC(ctx context.Context, req *user.GetLoginDataRequest) (string, []string, error) {
	resp, err := userClient.GetLoginData(ctx, req)
	if err != nil {
		logger.Errorf("GetLoginDataRPC: RPC called failed: %v", err.Error())
		return "", nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", nil, errno.BizError
	}
	return resp.Id, resp.Cookies, nil
}
