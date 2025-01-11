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

package mw

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/base/login_data"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// GetHeaderParams 获取请求头的信息，处理 id 和 cookies 并附加到 Context 中
func GetHeaderParams() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := string(c.GetHeader("Id"))
		cookies := string(c.GetHeader("Cookies"))
		if id == "" || cookies == "" {
			pack.RespError(c, errno.ParamMissingHeader)
			c.Abort()
			return
		}
		ctx = login_data.NewContext(ctx, &model.LoginData{
			Id:      id,
			Cookies: cookies,
		})

		// deliver to RPC server
		ctx = metainfoContext.WithLoginData(ctx, &model.LoginData{
			Id:      id,
			Cookies: cookies,
		})
		c.Next(ctx)
	}
}
