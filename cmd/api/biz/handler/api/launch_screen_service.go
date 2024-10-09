// Code generated by hertz generator.

package api

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	api "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
)

// CreateImage .
// @Summary CreateImage
// @Description create launch_screen image
// @Accept json/form
// @Produce json
// @Param pic_type query int true "1为空，2为页面跳转，3为app跳转"
// @Param duration query int false "展示时间"
// @Param href query string false "链接"
// @Param image form_data file true "图片"
// @Param start_at query int true "开始time(时间戳)"
// @Param end_at query int true "结束time(时间戳)"
// @Param s_type query int true "类型"
// @Param frequency query int false "一天展示次数"
// @Param start_time query int true "每日起始hour"
// @Param end_time query int true "每日结束hour"
// @Param authorization header string false "token"
// @Param text query string true "描述"
// @Param regex query string false "json类型，正则匹配项"
// @router /launch_screen/api/image [POST]
func CreateImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.CreateImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	token := c.GetHeader("authorization")
	imageFile, err := c.FormFile("image")
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.CreateImageResponse)

	if !pack.IsAllowImageExt(imageFile.Filename) {
		pack.RespError(c, errno.SuffixError)
		return
	}

	imageByte, err := pack.FileToByte(imageFile)
	if err != nil {
		pack.RespError(c, errno.BizError)
		return
	}

	respImage, err := rpc.CreateImageRPC(ctx, &launch_screen.CreateImageRequest{
		PicType:   req.PicType,
		Duration:  req.Duration,
		Href:      req.Href,
		Image:     imageByte,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
		SType:     req.SType,
		Frequency: req.Frequency,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Text:      req.Text,
		Regex:     req.Regex,
		Token:     string(token),
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Base = pack.BuildSuccessResp()
	resp.Picture = pack.BuildLaunchScreen(respImage)
	c.JSON(consts.StatusOK, resp)
}

// GetImage .
// @router /launch_screen/api/image [GET]
func GetImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetImagesByUserIdResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetImagesByUserId .
// @router /launch_screen/api/images [GET]
func GetImagesByUserId(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetImagesByUserIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetImagesByUserIdResponse)

	c.JSON(consts.StatusOK, resp)
}

// ChangeImageProperty .
// @router /launch_screen/api/image [PUT]
func ChangeImageProperty(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ChangeImagePropertyRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.ChangeImagePropertyResponse)

	c.JSON(consts.StatusOK, resp)
}

// ChangeImage .
// @router /launch_screen/api/image/img [PUT]
func ChangeImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ChangeImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.ChangeImageResponse)

	c.JSON(consts.StatusOK, resp)
}

// DeleteImage .
// @router /launch_screen/api/image [DELETE]
func DeleteImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DeleteImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.DeleteImageResponse)

	c.JSON(consts.StatusOK, resp)
}

// MobileGetImage .
// @router /launch_screen/api/screen [GET]
func MobileGetImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.MobileGetImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.MobileGetImageResponse)

	c.JSON(consts.StatusOK, resp)
}

// AddImagePointTime .
// @router /launch_screen/api/image/point [GET]
func AddImagePointTime(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.AddImagePointTimeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.AddImagePointTimeResponse)

	c.JSON(consts.StatusOK, resp)
}
