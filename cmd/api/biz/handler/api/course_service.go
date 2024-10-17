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

	"github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// GetCourseList .
// @router /api/v1/jwch/course/list [GET]
func GetCourseList(ctx context.Context, c *app.RequestContext) {
	user, err := api.GetLoginData(ctx)
	if err != nil {
		logger.Errorf("api.GetCourseList: GetLoginData error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	var req api.CourseListRequest

	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.GetCourseList: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	res, err := rpc.GetCourseListRPC(ctx, &course.CourseListRequest{
		LoginData: &model.LoginData{
			Id:      user.Id,
			Cookies: user.Cookies,
		},
		Term: req.Term,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	resp := new(api.CourseListResponse)
	resp.Data = pack.BuildCourseList(res)
	pack.RespList(c, resp.Data)
}
