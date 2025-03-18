/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by hertz generator.

package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/fzuhelper-server/api/mw"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// GetCourseList .
// @router /api/v1/jwch/course/list [GET]
func GetCourseList(ctx context.Context, c *app.RequestContext) {
	var req api.CourseListRequest
	var err error

	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	res, err := rpc.GetCourseListRPC(ctx, &course.CourseListRequest{
		Term:      req.Term,
		IsRefresh: req.IsRefresh,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	resp := new(api.CourseListResponse)
	resp.Data = pack.BuildCourseList(res)
	pack.RespList(c, resp.Data)
}

// GetTermList .
// @router /api/v1/jwch/term/list [GET]
func GetTermList(ctx context.Context, c *app.RequestContext) {
	var req api.TermListRequest
	var err error

	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	res, err := rpc.GetCourseTermsListRPC(ctx, &course.TermListRequest{})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	resp := new(api.CourseTermListResponse)
	resp.Data = res.Data
	pack.RespList(c, resp.Data)
}

// GetLocateDate .
// @router /api/v1/course/date [GET]
func GetLocateDate(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetLocateDateRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.GetLocateDateResponse)
	var date *model.LocateDate
	date, err = rpc.GetLocateDateRPC(ctx, &course.GetLocateDateRequest{})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.LocateDate = pack.BuildLocateDate(date)
	pack.RespData(c, resp.LocateDate)
}

// SubscribeCalendar .
// @router /api/v1/course/calendar/subscribe [GET]
func SubscribeCalendar(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.SubscribeCalendarRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.SubscribeCalendarResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetCalendar .
// @router /api/v1/jwch/course/calendar/token [GET]
func GetCalendar(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetCalendarTokenRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	// 生成 calendar token
	var token string
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		pack.RespError(c, errno.AuthError.WithError(err))
		return
	}
	token, err = mw.CreateToken(constants.TypeCalendarToken, utils.ParseJwchStuId(loginData.Id))
	if err != nil {
		pack.RespError(c, errno.AuthError.WithError(err))
		return
	}
	pack.RespData(c, token)
}
