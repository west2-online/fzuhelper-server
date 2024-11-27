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

package url

import (
	"context"

	"github.com/west2-online/fzuhelper-server/internal/url/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/url"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// UrlServiceImpl implements the last service interface defined in the IDL.
type UrlServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewUrlService(clientSet *base.ClientSet) *UrlServiceImpl {
	return &UrlServiceImpl{
		ClientSet: clientSet,
	}
}

// Login implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) Login(ctx context.Context, req *url.LoginRequest) (resp *url.LoginResponse, err error) {
	resp = new(url.LoginResponse)
	err = service.NewUrlService(ctx, s.ClientSet).Login(req)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

// UploadVersion implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) UploadVersion(ctx context.Context, req *url.UploadRequest) (resp *url.UploadResponse, err error) {
	resp = new(url.UploadResponse)
	err = service.NewUrlService(ctx, s.ClientSet).UploadVersion(req)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

// UploadParams implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) UploadParams(ctx context.Context, req *url.UploadParamsRequest) (resp *url.UploadParamsResponse, err error) {
	resp = new(url.UploadParamsResponse)
	policy, auth, err := service.NewUrlService(ctx, s.ClientSet).UploadParams(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Infof("Url.UploadParams: %v", err)
		return resp, nil
	}
	resp.Policy = &policy
	resp.Authorization = &auth
	return resp, nil
}

// DownloadReleaseApk implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) DownloadReleaseApk(ctx context.Context, req *url.DownloadReleaseApkRequest) (resp *url.DownloadReleaseApkResponse, err error) {
	resp = new(url.DownloadReleaseApkResponse)
	redirectUrl, err := service.NewUrlService(ctx, s.ClientSet).DownloadReleaseApk()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Infof("Url.DownloadReleaseApk: %v", err)
		return resp, nil
	}
	resp.RedirectUrl = redirectUrl
	return resp, nil
}

// DownloadBetaApk implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) DownloadBetaApk(ctx context.Context, req *url.DownloadBetaApkRequest) (resp *url.DownloadBetaApkResponse, err error) {
	resp = new(url.DownloadBetaApkResponse)
	redirectUrl, err := service.NewUrlService(ctx, s.ClientSet).DownloadBetaApk()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Infof("Url.DownloadReleaseApk: %v", err)
	}
	resp.RedirectUrl = redirectUrl
	return resp, nil
}

// GetReleaseVersion implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetReleaseVersion(ctx context.Context, req *url.GetReleaseVersionRequest) (resp *url.GetReleaseVersionResponse, err error) {
	resp = new(url.GetReleaseVersionResponse)
	version, err := service.NewUrlService(ctx, s.ClientSet).GetReleaseVersion()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Infof("Url.GetReleaseVersion: %v", err)
	}
	resp.Version = &version.Version
	resp.Url = &version.Url
	resp.Feature = &version.Feature
	resp.Code = &version.Code
	return resp, nil
}

// GetBetaVersion implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetBetaVersion(ctx context.Context, req *url.GetBetaVersionRequest) (resp *url.GetBetaVersionResponse, err error) {
	resp = new(url.GetBetaVersionResponse)
	version, err := service.NewUrlService(ctx, s.ClientSet).GetBetaVersion()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Infof("Url.GetBetaVersion: %v", err)
	}
	resp.Version = &version.Version
	resp.Url = &version.Url
	resp.Feature = &version.Feature
	resp.Code = &version.Code
	return resp, nil
}

// GetSetting implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetSetting(ctx context.Context, req *url.GetSettingRequest) (resp *url.GetSettingResponse, err error) {
	resp = new(url.GetSettingResponse)
	setting, err := service.NewUrlService(ctx, s.ClientSet).GetCloudSetting(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.Infof("Url.GetSetting: %v", err)
	}
	resp.Data = setting
	return resp, nil
}

// GetTest implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetTest(ctx context.Context, req *url.GetSettingRequest) (resp *url.GetTestResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCloud implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetCloud(ctx context.Context, req *url.GetCloudRequest) (resp *url.GetCloudResponse, err error) {
	// TODO: Your code here...
	return
}

// SetCloud implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) SetCloud(ctx context.Context, req *url.SetCloudRequest) (resp *url.SetCloudResponse, err error) {
	// TODO: Your code here...
	return
}

// GetDump implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetDump(ctx context.Context, req *url.GetDumpRequest) (resp *url.GetDumpResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCSS implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetCSS(ctx context.Context, req *url.GetCSSRequest) (resp *url.GetCSSResponse, err error) {
	// TODO: Your code here...
	return
}

// GetHtml implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetHtml(ctx context.Context, req *url.GetHtmlRequest) (resp *url.GetHtmlResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserAgreement implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetUserAgreement(ctx context.Context, req *url.GetUserAgreementRequest) (resp *url.GetUserAgreementResponse, err error) {
	// TODO: Your code here...
	return
}
