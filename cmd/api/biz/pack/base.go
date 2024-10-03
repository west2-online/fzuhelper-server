package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type Base struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

type BaseResp struct {
	Base Base
}

type RespWithData struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}

type DataList struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
}

func RespError(c *app.RequestContext, err error) {
	Errno := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, BaseResp{
		Base: Base{
			Code: Errno.ErrorCode,
			Msg:  Errno.ErrorMsg,
		},
	})
}

func RespSuccess(c *app.RequestContext) {
	Errno := errno.Success
	c.JSON(consts.StatusOK, BaseResp{
		Base: Base{
			Code: Errno.ErrorCode,
			Msg:  Errno.ErrorMsg,
		},
	})
}

func RespData(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, RespWithData{
		Code: errno.SuccessCode,
		Msg:  "Success",
		Data: data,
	})
}

func RespList(c *app.RequestContext, items any) {
	Errno := errno.Success
	resp := RespWithData{
		Code: Errno.ErrorCode,
		Msg:  Errno.ErrorMsg,
		Data: &DataList{
			Items: items,
		},
	}
	c.JSON(consts.StatusOK, resp)
}

func BuildSuccessResp() *api.BaseResp {
	return &api.BaseResp{
		Code: errno.SuccessCode,
		Msg:  errno.SuccessMsg,
	}
}
