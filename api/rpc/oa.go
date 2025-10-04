package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitOaRPC() {
	c, err := client.InitOaRPC()
	if err != nil {
		logger.Fatalf("api.rpc.oa InitOaRPC failed, err is %v", err)
	}
	logger.Infof("InitOaRPC: etcd=%s service=%s", config.Etcd.Addr, constants.OaServiceName)
	oaClient = *c
}

func CreateFeedbackRPC(ctx context.Context, req *oa.CreateFeedbackRequest) error {
	resp, err := oaClient.CreateFeedback(ctx, req)
	if err != nil {
		logger.Errorf("CreateFeedbackRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.BizError.WithMessage("创建反馈表单失败：" + resp.Base.Msg)
	}
	return nil
}

func GetFeedbackRPC(ctx context.Context, req *oa.GetFeedbackRequest) (*model.Feedback, error) {
	resp, err := oaClient.GetFeedback(ctx, req)
	if err != nil {
		logger.Errorf("GetFeedbackRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.BizError.WithMessage("获取反馈表单失败：" + resp.Base.Msg)
	}
	return resp.Data, nil
}
