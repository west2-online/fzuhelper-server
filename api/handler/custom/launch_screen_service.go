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

package custom

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// MobileGetImage .
// @router /launch_screen/api/screen [GET]
func MobileGetImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.MobileGetImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.MobileGetImage: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.MobileGetImageResponse)

	respImageList, _, err := rpc.MobileGetImageRPC(ctx, &launch_screen.MobileGetImageRequest{
		SType:     req.Type,
		StudentId: req.StudentID,
		College:   req.College,
		Device:    req.Device,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.PictureList = pack.BuildLaunchScreenList(respImageList)

	pack.CustomLaunchScreenRespList(c, resp.PictureList)
}

// AddImagePointTime .
// @router /launch_screen/api/image/point [GET]
func AddImagePointTime(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.AddImagePointTimeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.AddImagePointTime: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	_, err = rpc.AddImagePointTimeRPC(ctx, &launch_screen.AddImagePointTimeRequest{
		PictureId: req.PictureID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.CustomLaunchScreenRespSuccess(c)
}
