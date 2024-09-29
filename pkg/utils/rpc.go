package utils

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

//该文件负责处理通用的rpc结果处理

func IsSuccess(baseResp *model.BaseResp) error {
	if baseResp.Code != errno.SuccessCode {
		logger.LoggerObj.Errorf("utils.rpc.IsSuccess base code is not successful %v", baseResp)
		return errno.NewErrNo(baseResp.Code, baseResp.Msg)
	}
	return nil
}
