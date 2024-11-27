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

	"github.com/west2-online/fzuhelper-server/kitex_gen/url"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitUrlRPC() {
	client, err := client.InitUrlRPC()
	if err != nil {
		logger.Fatalf("api.rpc.launch_screen InitLaunchScreenRPC failed, err  %v", err)
	}
	urlClient = *client
}

func LoginRPC(ctx context.Context, req *url.LoginRequest) (err error) {
	resp, err := urlClient.Login(ctx, req)
	if err != nil {
		logger.Errorf("LoginRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func UploadVersionRPC(ctx context.Context, req *url.UploadRequest) (err error) {
	resp, err := urlClient.UploadVersion(ctx, req)
	if err != nil {
		logger.Errorf("UploadVersionRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func UploadParamsRPC(ctx context.Context, req *url.UploadParamsRequest) (*string, *string, error) {
	resp, err := urlClient.UploadParams(ctx, req)
	if err != nil {
		logger.Errorf("UploadParamsRPC: RPC called failed: %v", err.Error())
		return nil, nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Policy, resp.Authorization, nil
}

func DownloadReleaseApkRPC(ctx context.Context, req *url.DownloadReleaseApkRequest) (*string, error) {
	resp, err := urlClient.DownloadReleaseApk(ctx, req)
	if err != nil {
		logger.Errorf("DownloadReleaseApkRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return &resp.RedirectUrl, nil
}

func DownloadBetaApkRPC(ctx context.Context, req *url.DownloadBetaApkRequest) (*string, error) {
	resp, err := urlClient.DownloadBetaApk(ctx, req)
	if err != nil {
		logger.Errorf("DownloadBetaApkRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return &resp.RedirectUrl, nil
}

func GetReleaseVersionRPC(ctx context.Context, req *url.GetReleaseVersionRequest) (*url.GetReleaseVersionResponse, error) {
	resp, err := urlClient.GetReleaseVersion(ctx, req)
	if err != nil {
		logger.Errorf("GetReleaseVersionRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func GetBetaVersionRPC(ctx context.Context, req *url.GetBetaVersionRequest) (*url.GetBetaVersionResponse, error) {
	resp, err := urlClient.GetBetaVersion(ctx, req)
	if err != nil {
		logger.Errorf("GetBetaVersionRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func GetSettingRPC(ctx context.Context, req *url.GetSettingRequest) (*url.GetSettingResponse, error) {
	resp, err := urlClient.GetSetting(ctx, req)
	if err != nil {
		logger.Errorf("GetSettingRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func GetTestRPC(ctx context.Context, req *url.GetSettingRequest) (err error) {
	resp, err := urlClient.GetTest(ctx, req)
	if err != nil {
		logger.Errorf("GetTestRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func GetCloudRPC(ctx context.Context, req *url.GetCloudRequest) (*url.GetCloudResponse, error) {
	resp, err := urlClient.GetCloud(ctx, req)
	if err != nil {
		logger.Errorf("GetCloudRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func SetCloudRPC(ctx context.Context, req *url.SetCloudRequest) (err error) {
	resp, err := urlClient.SetCloud(ctx, req)
	if err != nil {
		logger.Errorf("SetCloudRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func GetDumpRPC(ctx context.Context, req *url.GetDumpRequest) (*url.GetDumpResponse, error) {
	resp, err := urlClient.GetDump(ctx, req)
	if err != nil {
		logger.Errorf("GetDumpRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp, nil
}

func GetCSSRPC(ctx context.Context, req *url.GetCSSRequest) (*[]byte, error) {
	resp, err := urlClient.GetCSS(ctx, req)
	if err != nil {
		logger.Errorf("GetCSSRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return &resp.Css, nil
}

func GetHtmlRPC(ctx context.Context, req *url.GetHtmlRequest) (*string, error) {
	resp, err := urlClient.GetHtml(ctx, req)
	if err != nil {
		logger.Errorf("GetHtmlRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return &resp.Html, nil
}

func GetUserAgreementRPC(ctx context.Context, req *url.GetUserAgreementRequest) (*string, error) {
	resp, err := urlClient.GetUserAgreement(ctx, req)
	if err != nil {
		logger.Errorf("GetUserAgreementRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	return &resp.UserAgreement, nil
}
