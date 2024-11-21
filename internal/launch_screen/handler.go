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

package launch_screen

import (
	"bytes"
	"context"

	"github.com/west2-online/fzuhelper-server/internal/launch_screen/pack"
	"github.com/west2-online/fzuhelper-server/internal/launch_screen/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// LaunchScreenServiceImpl implements the last service interface defined in the IDL.
type LaunchScreenServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewLaunchScreenService(clientSet *base.ClientSet) *LaunchScreenServiceImpl {
	return &LaunchScreenServiceImpl{
		ClientSet: clientSet,
	}
}

// CreateImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) CreateImage(stream launch_screen.LaunchScreenService_CreateImageServer) (err error) {
	resp := new(launch_screen.CreateImageResponse)
	// 首先取得除文件外的其他字段
	req, err := stream.Recv()
	if err != nil {
		logger.Errorf("LaunchScreen.CreateImage recv request: %+v", err)
		resp.Base = base.BuildBaseResp(err)
		return stream.SendAndClose(resp)
	}

	// 通过第一次获得的count来流式读取
	for i := 0; i < int(req.BufferCount); i++ {
		fileReq, err := stream.Recv()
		if err != nil {
			logger.Errorf("LaunchScreen.CreateImage recv file: %+v", err)
			resp.Base = base.BuildBaseResp(err)
			return stream.SendAndClose(resp)
		}
		req.Image = bytes.Join([][]byte{req.Image, fileReq.Image}, []byte(""))
	}

	pic, err := service.NewLaunchScreenService(stream.Context(), s.ClientSet).CreateImage(req)
	if err != nil {
		logger.Errorf("LaunchScreen.CreateImage: %+v", err)
		resp.Base = base.BuildBaseResp(err)
		return stream.SendAndClose(resp)
	}
	resp.Base = base.BuildSuccessResp()
	resp.Picture = pack.BuildImageResp(pic)
	return stream.SendAndClose(resp)
}

// GetImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) GetImage(ctx context.Context, req *launch_screen.GetImageRequest) (
	resp *launch_screen.GetImageResponse, err error,
) {
	resp = new(launch_screen.GetImageResponse)

	pic, err := service.NewLaunchScreenService(ctx, s.ClientSet).GetImageById(req.PictureId)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Errorf("LaunchScreen.GetImage: %+v", err)
		return resp, nil
	}
	resp.Picture = pack.BuildImageResp(pic)
	return resp, nil
}

// ChangeImageProperty implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) ChangeImageProperty(ctx context.Context,
	req *launch_screen.ChangeImagePropertyRequest,
) (resp *launch_screen.ChangeImagePropertyResponse, err error) {
	resp = new(launch_screen.ChangeImagePropertyResponse)
	pic, err := service.NewLaunchScreenService(ctx, s.ClientSet).UpdateImageProperty(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Errorf("LaunchScreen.ChangeImageProperty: %+v", err)
		return resp, nil
	}
	resp.Picture = pack.BuildImageResp(pic)
	return resp, nil
}

// ChangeImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) ChangeImage(stream launch_screen.LaunchScreenService_ChangeImageServer) (err error) {
	resp := new(launch_screen.ChangeImageResponse)
	// 首先取得除文件外的其他字段
	req, err := stream.Recv()
	if err != nil {
		logger.Errorf("LaunchScreen.ChangeImage recv request: %+v", err)
		resp.Base = base.BuildBaseResp(err)
		return stream.SendAndClose(resp)
	}

	// 通过第一次获得的count来流式读取
	for i := 0; i < int(req.BufferCount); i++ {
		fileReq, err := stream.Recv()
		if err != nil {
			logger.Errorf("LaunchScreen.ChangeImage recv file: %+v", err)
			resp.Base = base.BuildBaseResp(err)
			return stream.SendAndClose(resp)
		}
		req.Image = bytes.Join([][]byte{req.Image, fileReq.Image}, []byte(""))
	}
	pic, err := service.NewLaunchScreenService(stream.Context(), s.ClientSet).UpdateImagePath(req)
	if err != nil {
		logger.Errorf("LaunchScreen.ChangeImage: %+v", err)
		resp.Base = base.BuildBaseResp(err)
		return stream.SendAndClose(resp)
	}
	resp.Base = base.BuildSuccessResp()
	resp.Picture = pack.BuildImageResp(pic)
	return stream.SendAndClose(resp)
}

// DeleteImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) DeleteImage(ctx context.Context, req *launch_screen.DeleteImageRequest) (resp *launch_screen.DeleteImageResponse, err error) {
	resp = new(launch_screen.DeleteImageResponse)

	err = service.NewLaunchScreenService(ctx, s.ClientSet).DeleteImage(req.PictureId)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		logger.Errorf("LaunchScreen.DeleteImage: %+v", err)
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	// resp.Picture = pack.BuildImageResp(pic)
	return resp, nil
}

// MobileGetImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) MobileGetImage(ctx context.Context, req *launch_screen.MobileGetImageRequest) (
	resp *launch_screen.MobileGetImageResponse, err error,
) {
	resp = new(launch_screen.MobileGetImageResponse)
	pictureList, cnt, err := service.NewLaunchScreenService(ctx, s.ClientSet).MobileGetImage(req)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		logger.Errorf("LaunchScreen.MobileGetImage: %+v", err)
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	resp.Count = &cnt
	resp.PictureList = pack.BuildImagesResp(pictureList)
	return resp, nil
}

// AddImagePointTime implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) AddImagePointTime(ctx context.Context, req *launch_screen.AddImagePointTimeRequest) (
	resp *launch_screen.AddImagePointTimeResponse, err error,
) {
	resp = new(launch_screen.AddImagePointTimeResponse)
	err = service.NewLaunchScreenService(ctx, s.ClientSet).AddPointTime(req.PictureId)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Errorf("LaunchScreen.AddImagePointTime: %+v", err)
	}
	return resp, nil
}
