package utils

import (
	"errors"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// IsSuccess 通用的rpc结果处理
func IsSuccess(baseResp *model.BaseResp) error {
	if baseResp.Code != errno.SuccessCode {
		return errors.New("utils.IsSuccess: the base code is not successful")
	}
	return nil
}
