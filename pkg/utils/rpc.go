package utils

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

//处理通用rpc请求

func IsSuccess(err error, baseResp *model.BaseResp) (bool, error) {
	if err != nil {
		logger.LoggerObj.Errorf("api.rpc.classroom GetEmptyRoomRPC received rpc error %v", err)
		return false, err
	}
	if baseResp.Code != errno.SuccessCode {
		logger.LoggerObj.Errorf("api.rpc.classroom GetEmptyRoomRPC received failed")
		return false, errno.NewErrNo(baseResp.Code, baseResp.Msg)
	}
	return true, nil
}
