package main

import (
	"context"
	version "github.com/west2-online/fzuhelper-server/kitex_gen/version"
)

// VersionServiceImpl implements the last service interface defined in the IDL.
type VersionServiceImpl struct{}

// Login implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) Login(ctx context.Context, req *version.LoginRequest) (resp *version.LoginResponse, err error) {
	// TODO: Your code here...
	return
}

// UploadVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) UploadVersion(ctx context.Context, req *version.UploadRequest) (resp *version.UploadResponse, err error) {
	// TODO: Your code here...
	return
}

// UploadParams implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) UploadParams(ctx context.Context, req *version.UploadParamsRequest) (resp *version.UploadParamsResponse, err error) {
	// TODO: Your code here...
	return
}

// DownloadReleaseApk implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) DownloadReleaseApk(ctx context.Context, req *version.DownloadReleaseApkRequest) (resp *version.DownloadReleaseApkResponse, err error) {
	// TODO: Your code here...
	return
}

// DownloadBetaApk implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) DownloadBetaApk(ctx context.Context, req *version.DownloadBetaApkRequest) (resp *version.DownloadBetaApkResponse, err error) {
	// TODO: Your code here...
	return
}

// GetReleaseVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetReleaseVersion(ctx context.Context, req *version.GetReleaseVersionRequest) (resp *version.GetReleaseVersionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetBetaVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetBetaVersion(ctx context.Context, req *version.GetBetaVersionRequest) (resp *version.GetBetaVersionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetSetting implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetSetting(ctx context.Context, req *version.GetSettingRequest) (resp *version.GetSettingResponse, err error) {
	// TODO: Your code here...
	return
}

// GetTest implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetTest(ctx context.Context, req *version.GetTestRequest) (resp *version.GetTestResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCloud implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetCloud(ctx context.Context, req *version.GetCloudRequest) (resp *version.GetCloudResponse, err error) {
	// TODO: Your code here...
	return
}

// SetCloud implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) SetCloud(ctx context.Context, req *version.SetCloudRequest) (resp *version.SetCloudResponse, err error) {
	// TODO: Your code here...
	return
}

// GetDump implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetDump(ctx context.Context, req *version.GetDumpRequest) (resp *version.GetDumpResponse, err error) {
	// TODO: Your code here...
	return
}

// AndroidGetVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) AndroidGetVersion(ctx context.Context, req *version.AndroidGetVersioneRequest) (resp *version.AndroidGetVersionResponse, err error) {
	// TODO: Your code here...
	return
}
