package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"strings"
)

//获取请求头的信息

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
		//将id和cookies放入context中
		context.WithValue(ctx, "id", id)
		context.WithValue(ctx, "cookies", cookies)
		c.Next(ctx)
	}
}
