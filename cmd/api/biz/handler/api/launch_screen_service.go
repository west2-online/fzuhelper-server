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

	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"

	api "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
)

// CreateImage .
// @Summary CreateImage
// @Description create launch_screen image
// @Accept json/form
// @Produce json
// @Param type query int true "1为空，2为页面跳转，3为app跳转"
// @Param duration query int false "展示时间"
// @Param href query string false "链接"
// @Param image formData file true "图片"
// @Param start_at query int true "开始time(时间戳)"
// @Param end_at query int true "结束time(时间戳)"
// @Param s_type query int true "类型"
// @Param frequency query int false "一天展示次数"
// @Param start_time query int true "每日起始hour"
// @Param end_time query int true "每日结束hour"
// @Param text query string true "描述"
// @Param regex query int true "regex"
// @router /launch_screen/api/image [POST]
func CreateImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.CreateImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.CreateImage: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}
	imageFile, err := c.FormFile("image")
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.CreateImageResponse)

	if !utils.IsAllowImageFile(imageFile) {
		pack.RespError(c, errno.SuffixError)
		return
	}

	imageByte, err := utils.FileToByteArray(imageFile)
	if err != nil {
		pack.RespError(c, errno.BizError)
		return
	}

	respImage, err := rpc.CreateImageRPC(ctx, &launch_screen.CreateImageRequest{
		PicType:     req.PicType,
		Duration:    req.Duration,
		Href:        req.Href,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		SType:       req.SType,
		Frequency:   req.Frequency,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Text:        req.Text,
		Regex:       req.Regex,
		BufferCount: int64(len(imageByte)),
	}, imageByte)
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Picture = pack.BuildLaunchScreen(respImage)
	pack.RespData(c, resp.Picture)
}

// GetImage .
// @Summary GetImage
// @Description get image
// @Accept json/form
// @Produce json
// @Param picture_id query int true "图片id"
// @router /launch_screen/api/image [GET]
func GetImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.GetImage: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.GetImageResponse)

	respImage, err := rpc.GetImageRPC(ctx, &launch_screen.GetImageRequest{
		PictureId: req.PictureID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Picture = pack.BuildLaunchScreen(respImage)
	pack.RespData(c, resp.Picture)
}

// ChangeImageProperty .
// @Summary ChangeImageProperty
// @Description change image's properties
// @Accept json/form
// @Produce json
// @Param picture_id query int true "图片id"
// @Param type query int true "1为空，2为页面跳转，3为app跳转"
// @Param duration query int false "展示时间"
// @Param href query string false "链接"
// @Param start_at query int true "开始time(时间戳)"
// @Param end_at query int true "结束time(时间戳)"
// @Param s_type query int true "类型"
// @Param frequency query int false "一天展示次数"
// @Param start_time query int true "每日起始hour"
// @Param end_time query int true "每日结束hour"
// @Param text query string true "描述"
// @Param regex query int true "regex"
// @router /launch_screen/api/image [PUT]
func ChangeImageProperty(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ChangeImagePropertyRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.ChangeImageProperty: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.ChangeImagePropertyResponse)

	respImage, err := rpc.ChangeImagePropertyRPC(ctx, &launch_screen.ChangeImagePropertyRequest{
		PictureId: req.PictureID,
		PicType:   req.PicType,
		Duration:  req.Duration,
		Href:      req.Href,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
		SType:     req.SType,
		Frequency: req.Frequency,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Text:      req.Text,
		Regex:     req.Regex,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Picture = pack.BuildLaunchScreen(respImage)
	pack.RespData(c, resp.Picture)
}

// ChangeImage .
// @Summary ChangeImage
// @Description change image
// @Accept json/form
// @Produce json
// @Param picture_id query int true "图片id"
// @Param image formData file true "图片"
// @Param stu_id query int true "学生id"
// @router /launch_screen/api/image/img [PUT]
func ChangeImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ChangeImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.ChangeImage: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}
	imageFile, err := c.FormFile("image")
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}
	resp := new(api.ChangeImageResponse)

	if !utils.IsAllowImageFile(imageFile) {
		pack.RespError(c, errno.SuffixError)
		return
	}

	imageByte, err := utils.FileToByteArray(imageFile)
	if err != nil {
		pack.RespError(c, errno.BizError)
		return
	}

	respImage, err := rpc.ChangeImageRPC(ctx, &launch_screen.ChangeImageRequest{
		PictureId:   req.PictureID,
		BufferCount: int64(len(imageByte)),
	}, imageByte)
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Picture = pack.BuildLaunchScreen(respImage)
	pack.RespData(c, resp.Picture)
}

// DeleteImage .
// @Summary DeleteImage
// @Description delete image
// @Accept json/form
// @Produce json
// @Param picture_id query int true "图片id"
// @router /launch_screen/api/image [DELETE]
func DeleteImage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DeleteImageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.DeleteImage: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	_, err = rpc.DeleteImageRPC(ctx, &launch_screen.DeleteImageRequest{
		PictureId: req.PictureID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespSuccess(c)
}

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

	pack.RespList(c, resp.PictureList)
}

// AddImagePointTime .
// @Summary AddImagePointTime
// @Description add image point time(for frontend)
// @Accept json/form
// @Produce json
// @Param picture_id query int true "图片id"
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
	pack.RespSuccess(c)
}
