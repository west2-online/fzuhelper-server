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

package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// 获取请求头的信息

func GetHeaderParams() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := string(c.GetHeader("id"))
		temp := string(c.GetHeader("cookies"))
		if id == "" || len(temp) == 0 {
			pack.RespError(c, errno.ParamMissingHeader)
			c.Abort()
			return
		}
		cookies := strings.Split(temp, ",")
		// 将id和cookies放入context中
		fmt.Println(id)
		fmt.Println(cookies)
		ctx = api.NewContext(ctx, &model.LoginData{
			Id:      id,
			Cookies: cookies,
		})
		c.Next(ctx)
	}
}
