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

package main

import (
	"context"

	url "github.com/west2-online/fzuhelper-server/kitex_gen/url"
)

// UrlServiceImpl implements the last service interface defined in the IDL.
type UrlServiceImpl struct{}

// Login implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) Login(ctx context.Context, req *url.LoginRequest) (resp *url.LoginResponse, err error) {
	// TODO: Your code here...
	return
}

// UploadVersion implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) UploadVersion(ctx context.Context, req *url.UploadRequest) (resp *url.UploadResponse, err error) {
	// TODO: Your code here...
	return
}

// UploadParams implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) UploadParams(ctx context.Context, req *url.UploadParamsRequest) (resp *url.UploadParamsResponse, err error) {
	// TODO: Your code here...
	return
}

// DownloadReleaseApk implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) DownloadReleaseApk(ctx context.Context, req *url.DownloadReleaseApkRequest) (resp *url.DownloadReleaseApkResponse, err error) {
	// TODO: Your code here...
	return
}

// DownloadBetaApk implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) DownloadBetaApk(ctx context.Context, req *url.DownloadBetaApkRequest) (resp *url.DownloadBetaApkResponse, err error) {
	// TODO: Your code here...
	return
}

// GetReleaseVersion implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetReleaseVersion(ctx context.Context, req *url.GetReleaseVersionRequest) (resp *url.GetReleaseVersionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetBetaVersion implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetBetaVersion(ctx context.Context, req *url.GetBetaVersionRequest) (resp *url.GetBetaVersionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetSetting implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) GetSetting(ctx context.Context, req *url.GetSettingRequest) (resp *url.GetSettingResponse, err error) {
	// TODO: Your code here...
	return
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
