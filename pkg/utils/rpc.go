package utils

import (
	"github.com/pkg/errors"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

//该文件负责处理通用的rpc结果处理

func IsSuccess(baseResp *model.BaseResp) error {
	if baseResp.Code != errno.SuccessCode {
		return errors.Wrap(errno.NewErrNo(baseResp.Code, baseResp.Msg), "IsSuccess: the base code is not successful")
	}
	return nil
}
