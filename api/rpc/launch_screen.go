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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitLaunchScreenRPC() {
	c, err := client.InitLaunchScreenRPC()
	if err != nil {
		logger.Fatalf("api.rpc.launch_screen InitLaunchScreenRPC failed, err  %v", err)
	}
	launchScreenClient = *c
}

func InitLaunchScreenStreamRPC() {
	c, err := client.InitLaunchScreenStreamRPC()
	if err != nil {
		logger.Fatalf("api.rpc.launch_screen InitLaunchScreenStreamRPC failed, err  %v", err)
	}
	launchScreenStreamClient = *c
}

func CreateImageRPC(ctx context.Context, req *launch_screen.CreateImageRequest, file [][]byte) (image *model.Picture, err error) {
	stream, err := launchScreenStreamClient.CreateImage(ctx)
	if err != nil {
		logger.Errorf("CreateImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	// 第一次先发送字段
	err = stream.Send(req)
	if err != nil {
		logger.Errorf("CreateImageRPC: RPC stream failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	// 之后发送文件
	for _, fileBlock := range file {
		err = stream.Send(&launch_screen.CreateImageRequest{Image: fileBlock})
		if err != nil {
			logger.Errorf("CreateImageRPC: RPC stream failed: %v", err.Error())
			return nil, errno.InternalServiceError.WithError(err)
		}
	}
	// 终止传输
	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("CreateImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	return resp.Picture, nil
}

func GetImageRPC(ctx context.Context, req *launch_screen.GetImageRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.GetImage(ctx, req)
	if err != nil {
		logger.Errorf("GetImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}

func ChangeImagePropertyRPC(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.ChangeImageProperty(ctx, req)
	if err != nil {
		logger.Errorf("ChangeImagePropertyRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}

func ChangeImageRPC(ctx context.Context, req *launch_screen.ChangeImageRequest, file [][]byte) (image *model.Picture, err error) {
	stream, err := launchScreenStreamClient.ChangeImage(ctx)
	if err != nil {
		logger.Errorf("ChangeImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	// 第一次先发送字段
	err = stream.Send(req)
	if err != nil {
		logger.Errorf("ChangeImageRPC: RPC stream failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	// 之后发送文件
	for _, fileBlock := range file {
		err = stream.Send(&launch_screen.ChangeImageRequest{Image: fileBlock})
		if err != nil {
			logger.Errorf("ChangeImageRPC: RPC stream failed: %v", err.Error())
			return nil, errno.InternalServiceError.WithError(err)
		}
	}
	// 终止传输
	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("ChangeImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	return resp.Picture, nil
}

func DeleteImageRPC(ctx context.Context, req *launch_screen.DeleteImageRequest) (err error) {
	resp, err := launchScreenClient.DeleteImage(ctx, req)
	if err != nil {
		logger.Errorf("DeleteImageRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func MobileGetImageRPC(ctx context.Context, req *launch_screen.MobileGetImageRequest) (image []*model.Picture, cnt *int64, err error) {
	resp, err := launchScreenClient.MobileGetImage(ctx, req)
	if err != nil {
		logger.Errorf("MobileGetImageRPC: RPC called failed: %v", err.Error())
		return nil, nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.PictureList, resp.Count, nil
}

func AddImagePointTimeRPC(ctx context.Context, req *launch_screen.AddImagePointTimeRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.AddImagePointTime(ctx, req)
	if err != nil {
		logger.Errorf("AddImagePointTimeRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}
