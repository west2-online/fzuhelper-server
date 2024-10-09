package utils

import (
	"errors"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch/errno"
)

func BuildBaseResp(err error) *model.BaseResp {
	if err == nil {
		return ErrToResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return ErrToResp(e)
	}

	_e := errno.ServiceError.WithMessage(err.Error()) //未知错误
	return ErrToResp(_e)
}

func ErrToResp(err errno.ErrNo) *model.BaseResp {
	return &model.BaseResp{
		Code: err.ErrorCode,
		Msg:  err.ErrorMsg,
	}
}
