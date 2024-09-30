package middleware

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
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
		fmt.Println(id)
		fmt.Println(cookies)
		ctx = api.NewContext(ctx, &model.LoginData{
			Id:      id,
			Cookies: cookies,
		})
		c.Next(ctx)
	}
}
