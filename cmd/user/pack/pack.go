package pack

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func BuildBaseResp(err error) *user.BaseResp {
	if err == nil {
		return &user.BaseResp{
			Code: errno.SuccessCode,
			Msg:  errno.Success.ErrorMsg,
		}
	}
	Errno := errno.ConvertErr(err)
	return &user.BaseResp{
		Code: Errno.ErrorCode,
		Msg:  Errno.ErrorMsg,
	}
}
