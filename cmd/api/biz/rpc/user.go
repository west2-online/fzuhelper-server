package rpc

import (
	"context"
	"fmt"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitUserRPC() {
	client, err := client.InitUserRPC()
	if err != nil {
		logger.LoggerObj.Fatalf("api.rpc.user InitUserRPC failed, err is %v", err)
	}
	userClient = *client
}

func GetLoginDataRPC(ctx context.Context, req *user.GetLoginDataRequest) (string, []string, error) {
	resp, err := userClient.GetLoginData(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("GetLoginDataRPC: received rpc error %w", err)
	}
	if err = utils.IsSuccess(resp.Base); err != nil {
		return "", nil, fmt.Errorf("GetLoginDataRPC: base code is not successful %w", err)
	}
	return resp.Id, resp.Cookies, nil
}
