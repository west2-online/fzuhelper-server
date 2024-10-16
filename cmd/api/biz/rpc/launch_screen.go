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

	"github.com/west2-online/fzuhelper-server/pkg/client"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitLaunchScreenRPC() {
	client, err := client.InitLaunchScreenRPC()
	if err != nil {
		logger.Fatalf("api.rpc.launch_screen InitLaunchScreenRPC failed, err  %v", err)
	}
	launchScreenClient = *client
}

func CreateImageRPC(ctx context.Context, req *launch_screen.CreateImageRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.CreateImage(ctx, req)
	if err != nil {
		logger.Errorf("CreateImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}

func GetImageRPC(ctx context.Context, req *launch_screen.GetImageRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.GetImage(ctx, req)
	if err != nil {
		logger.Errorf("GetImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}

func GetImagesByIdRPC(ctx context.Context, req *launch_screen.GetImagesByUserIdRequest) (image []*model.Picture, cnt *int64, err error) {
	resp, err := launchScreenClient.GetImagesByUserId(ctx, req)
	if err != nil {
		logger.Errorf("GetImagesByIdRPC: RPC called failed: %v", err.Error())
		return nil, nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.PictureList, resp.Count, nil
}

func ChangeImagePropertyRPC(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.ChangeImageProperty(ctx, req)
	if err != nil {
		logger.Errorf("ChangeImagePropertyRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}

func ChangeImageRPC(ctx context.Context, req *launch_screen.ChangeImageRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.ChangeImage(ctx, req)
	if err != nil {
		logger.Errorf("ChangeImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}

func DeleteImageRPC(ctx context.Context, req *launch_screen.DeleteImageRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.DeleteImage(ctx, req)
	if err != nil {
		logger.Errorf("DeleteImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}

func MobileGetImageRPC(ctx context.Context, req *launch_screen.MobileGetImageRequest) (image []*model.Picture, cnt *int64, err error) {
	resp, err := launchScreenClient.MobileGetImage(ctx, req)
	if err != nil {
		logger.Errorf("MobileGetImageRPC: RPC called failed: %v", err.Error())
		return nil, nil, errno.InternalServiceError.WithMessage(err.Error())
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
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Picture, nil
}