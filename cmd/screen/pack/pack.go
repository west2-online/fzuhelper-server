package pack

import (
	"errors"

	"github.com/ozline/tiktok/kitex_gen/screen"
	"github.com/ozline/tiktok/pkg/errno"
)

func BuildBaseResp(err error) *screen.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}

	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceError.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *screen.BaseResp {
	return &screen.BaseResp{
		Code: err.ErrorCode,
		Msg:  err.ErrorMsg,
	}
}
