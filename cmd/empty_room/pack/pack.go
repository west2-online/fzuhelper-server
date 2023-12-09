package pack

import (
	"errors"

	"github.com/west2-online/fzuhelper-server/kitex_gen/empty_room"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func BuildBaseResp(err error) *empty_room.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}

	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.InternalServiceError.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *empty_room.BaseResp {
	return &empty_room.BaseResp{
		Code: err.ErrorCode,
		Msg:  err.ErrorMsg,
	}
}
