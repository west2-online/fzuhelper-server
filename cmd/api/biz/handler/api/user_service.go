// Code generated by hertz generator.

package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// GetLoginData .
// @router /api/v1/user/login [GET]
func GetLoginData(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetLoginDataRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.GetLoginData: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamEmpty)
		return
	}
	resp := new(api.GetLoginDataResponse)
	id, cookies, err := rpc.GetLoginDataRPC(ctx, &user.GetLoginDataRequest{
		Id:       req.ID,
		Password: req.Password,
	})
	if err != nil {
		logger.Errorf("api.GetLoginData: GetEmptyRoomRPC  error %v", err)
		pack.RespError(c, errno.InternalServiceError)
		return
	}
	resp.ID = id
	resp.Cookies = cookies
	pack.RespData(c, resp)
}
