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
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// GetHeaderParams 获取请求头的信息，处理 id 和 cookies 并附加到 Context 中
func GetHeaderParams() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := string(c.GetHeader("Id"))
		temp := string(c.GetHeader("Cookies"))
		if id == "" || len(temp) == 0 {
			logger.Errorf("GetHeaderParams: Id or Cookies is empty")
			pack.RespError(c, errno.ParamMissingHeader)
			c.Abort()
			return
		}
		cookies := strings.Split(temp, ",")
		ctx = api.NewContext(ctx, &model.LoginData{
			Id:      id,
			Cookies: cookies,
		})

		c.Next(ctx)
	}
}
