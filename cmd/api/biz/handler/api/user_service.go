// Code generated by hertz generator.

package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	api "github.com/west2-online/fzuhelper-server/biz/model/api"
)

// GetUserInfo .
// @router /user/info [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UserInfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.UserInfoResp)

	c.JSON(consts.StatusOK, resp)
}

// ValidateCode .
// @router user/validateCode [POST]
func ValidateCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ValidateCodeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.ValidateCodeResp)

	c.JSON(consts.StatusOK, resp)
}

// ChangePassword .
// @router /user/info [PUT]
func ChangePassword(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ChangePasswordRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.ChangePasswordResp)

	c.JSON(consts.StatusOK, resp)
}

// GetSchoolCalendar .
// @router /user/schoolCalendar [GET]
func GetSchoolCalendar(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.SchoolCalendarRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.SchoolCalendarResponse)

	c.JSON(consts.StatusOK, resp)
}
