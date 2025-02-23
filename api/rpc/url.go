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

	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitVersionRPC() {
	client, err := client.InitVersionRPC()
	if err != nil {
		logger.Fatalf("api.rpc.version InitVersionRPC failed, err  %v", err)
	}
	versionClient = *client
}

func LoginRPC(ctx context.Context, req *version.LoginRequest) (err error) {
	resp, err := versionClient.Login(ctx, req)
	if err != nil {
		logger.Errorf("LoginRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func UploadVersionRPC(ctx context.Context, req *version.UploadRequest) (err error) {
	resp, err := versionClient.UploadVersion(ctx, req)
	if err != nil {
		logger.Errorf("UploadVersionRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func UploadParamsRPC(ctx context.Context, req *version.UploadParamsRequest) (*string, *string, error) {
	resp, err := versionClient.UploadParams(ctx, req)
	if err != nil {
		logger.Errorf("UploadParamsRPC: RPC called failed: %v", err.Error())
		return nil, nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Policy, resp.Authorization, nil
}

func DownloadReleaseApkRPC(ctx context.Context, req *version.DownloadReleaseApkRequest) (*string, error) {
	resp, err := versionClient.DownloadReleaseApk(ctx, req)
	if err != nil {
		logger.Errorf("DownloadReleaseApkRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return &resp.RedirectUrl, nil
}

func DownloadBetaApkRPC(ctx context.Context, req *version.DownloadBetaApkRequest) (*string, error) {
	resp, err := versionClient.DownloadBetaApk(ctx, req)
	if err != nil {
		logger.Errorf("DownloadBetaApkRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return &resp.RedirectUrl, nil
}

func GetReleaseVersionRPC(ctx context.Context, req *version.GetReleaseVersionRequest) (*version.GetReleaseVersionResponse, error) {
	resp, err := versionClient.GetReleaseVersion(ctx, req)
	if err != nil {
		logger.Errorf("GetReleaseVersionRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func GetBetaVersionRPC(ctx context.Context, req *version.GetBetaVersionRequest) (*version.GetBetaVersionResponse, error) {
	resp, err := versionClient.GetBetaVersion(ctx, req)
	if err != nil {
		logger.Errorf("GetBetaVersionRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func GetSettingRPC(ctx context.Context, req *version.GetSettingRequest) (*version.GetSettingResponse, error) {
	resp, err := versionClient.GetSetting(ctx, req)
	if err != nil {
		logger.Errorf("GetSettingRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func GetTestRPC(ctx context.Context, req *version.GetTestRequest) (*version.GetTestResponse, error) {
	resp, err := versionClient.GetTest(ctx, req)
	if err != nil {
		logger.Errorf("GetTestRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp, nil
}

func GetCloudRPC(ctx context.Context, req *version.GetCloudRequest) (*version.GetCloudResponse, error) {
	resp, err := versionClient.GetCloud(ctx, req)
	if err != nil {
		logger.Errorf("GetCloudRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func SetCloudRPC(ctx context.Context, req *version.SetCloudRequest) (err error) {
	resp, err := versionClient.SetCloud(ctx, req)
	if err != nil {
		logger.Errorf("SetCloudRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func GetDumpRPC(ctx context.Context, req *version.GetDumpRequest) (*version.GetDumpResponse, error) {
	resp, err := versionClient.GetDump(ctx, req)
	if err != nil {
		logger.Errorf("GetDumpRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func AndroidVersionRPC(ctx context.Context, req *version.AndroidGetVersioneRequest) (*version.AndroidGetVersionResponse, error) {
	resp, err := versionClient.AndroidGetVersion(ctx, req)
	if err != nil {
		logger.Errorf("AndroidVersionRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}
