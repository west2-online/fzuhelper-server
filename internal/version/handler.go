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

package version

import (
	"context"

	"github.com/west2-online/fzuhelper-server/internal/version/pack"
	"github.com/west2-online/fzuhelper-server/internal/version/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/singleflight"
)

type androidVersionResult struct {
	release *pack.Version
	beta    *pack.Version
}

// VersionServiceImpl implements the last service interface defined in the IDL.
type VersionServiceImpl struct {
	ClientSet    *base.ClientSet
	singleflight singleflight.Group
}

func NewVersionService(clientSet *base.ClientSet) *VersionServiceImpl {
	return &VersionServiceImpl{
		ClientSet: clientSet,
	}
}

// Login implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) Login(ctx context.Context, req *version.LoginRequest) (resp *version.LoginResponse, err error) {
	resp = new(version.LoginResponse)
	err = service.NewVersionService(ctx, s.ClientSet).Login(req)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

// UploadVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) UploadVersion(ctx context.Context, req *version.UploadRequest) (resp *version.UploadResponse, err error) {
	resp = new(version.UploadResponse)
	err = service.NewVersionService(ctx, s.ClientSet).UploadVersion(req)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

// UploadParams implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) UploadParams(ctx context.Context, req *version.UploadParamsRequest) (resp *version.UploadParamsResponse, err error) {
	resp = new(version.UploadParamsResponse)
	policy, auth, err := service.NewVersionService(ctx, s.ClientSet).UploadParams(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.UploadParams: %v", err)
		return resp, nil
	}
	resp.Policy = &policy
	resp.Authorization = &auth
	return resp, nil
}

// DownloadReleaseApk implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) DownloadReleaseApk(ctx context.Context, req *version.DownloadReleaseApkRequest) (
	resp *version.DownloadReleaseApkResponse, err error,
) {
	resp = new(version.DownloadReleaseApkResponse)
	// 下载地址按发布渠道拆 key，避免 release/beta 并发时复用到另一个渠道的地址。
	v, err := s.singleflight.Do("download_release", func() (any, error) {
		return service.NewVersionService(ctx, s.ClientSet).DownloadReleaseApk()
	})
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.DownloadReleaseApk: %v", err)
		return resp, nil
	}
	redirectUrl, ok := v.(string)
	if !ok {
		resp.Base = base.BuildBaseResp(singleflight.ErrInvalidType)
		return resp, nil
	}
	resp.RedirectUrl = redirectUrl
	return resp, nil
}

// DownloadBetaApk implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) DownloadBetaApk(ctx context.Context, req *version.DownloadBetaApkRequest) (
	resp *version.DownloadBetaApkResponse, err error,
) {
	resp = new(version.DownloadBetaApkResponse)
	v, err := s.singleflight.Do("download_beta", func() (any, error) {
		return service.NewVersionService(ctx, s.ClientSet).DownloadBetaApk()
	})
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.DownloadBetaApk: %v", err)
		return resp, nil
	}
	redirectUrl, ok := v.(string)
	if !ok {
		resp.Base = base.BuildBaseResp(singleflight.ErrInvalidType)
		return resp, nil
	}
	resp.RedirectUrl = redirectUrl
	return resp, nil
}

// GetReleaseVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetReleaseVersion(ctx context.Context, req *version.GetReleaseVersionRequest) (
	resp *version.GetReleaseVersionResponse, err error,
) {
	resp = new(version.GetReleaseVersionResponse)
	// 版本信息按发布渠道拆 key，避免 release/beta 版本元数据互相复用。
	v, err := s.singleflight.Do("release_version", func() (any, error) {
		return service.NewVersionService(ctx, s.ClientSet).GetReleaseVersion()
	})
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.GetReleaseVersion: %v", err)
		return resp, nil
	}
	res, ok := v.(*pack.Version)
	if !ok {
		resp.Base = base.BuildBaseResp(singleflight.ErrInvalidType)
		return resp, nil
	}
	resp.Version = &res.Version
	resp.Url = &res.Url
	resp.Feature = &res.Feature
	resp.Code = &res.Code
	resp.Force = &res.Force
	return resp, nil
}

// GetBetaVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetBetaVersion(ctx context.Context, req *version.GetBetaVersionRequest) (resp *version.GetBetaVersionResponse, err error) {
	resp = new(version.GetBetaVersionResponse)
	v, err := s.singleflight.Do("beta_version", func() (any, error) {
		return service.NewVersionService(ctx, s.ClientSet).GetBetaVersion()
	})
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.GetBetaVersion: %v", err)
		return resp, nil
	}
	res, ok := v.(*pack.Version)
	if !ok {
		resp.Base = base.BuildBaseResp(singleflight.ErrInvalidType)
		return resp, nil
	}
	resp.Version = &res.Version
	resp.Url = &res.Url
	resp.Feature = &res.Feature
	resp.Code = &res.Code
	resp.Force = &res.Force
	return resp, nil
}

// GetSetting implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetSetting(ctx context.Context, req *version.GetSettingRequest) (resp *version.GetSettingResponse, err error) {
	resp = new(version.GetSettingResponse)
	setting, err := service.NewVersionService(ctx, s.ClientSet).GetCloudSetting(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.GetSetting: %v", err)
	}
	resp.CloudSetting = *setting
	return resp, nil
}

// GetTest implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetTest(ctx context.Context, req *version.GetTestRequest) (resp *version.GetTestResponse, err error) {
	resp = new(version.GetTestResponse)
	setting, err := service.NewVersionService(ctx, s.ClientSet).TestSetting(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.GetTest: %v", err)
	}
	resp.CloudSetting = *setting
	return resp, nil
}

// GetCloud implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetCloud(ctx context.Context, req *version.GetCloudRequest) (resp *version.GetCloudResponse, err error) {
	resp = new(version.GetCloudResponse)
	v, err := s.singleflight.Do("cloud", func() (any, error) {
		return service.NewVersionService(ctx, s.ClientSet).GetAllCloudSetting()
	})
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.GetCloud: %v", err)
		return resp, nil
	}
	setting, ok := v.(*[]byte)
	if !ok {
		resp.Base = base.BuildBaseResp(singleflight.ErrInvalidType)
		return resp, nil
	}
	resp.CloudSetting = *setting
	return resp, nil
}

// SetCloud implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) SetCloud(ctx context.Context, req *version.SetCloudRequest) (resp *version.SetCloudResponse, err error) {
	resp = new(version.SetCloudResponse)
	err = service.NewVersionService(ctx, s.ClientSet).SetSetting(req)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

// GetDump implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) GetDump(ctx context.Context, req *version.GetDumpRequest) (resp *version.GetDumpResponse, err error) {
	resp = new(version.GetDumpResponse)
	dump, err := service.NewVersionService(ctx, s.ClientSet).GetDump()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.GetDump: %v", err)
		return resp, nil
	}
	resp.Data = dump
	return resp, nil
}

// AndroidGetVersion implements the VersionServiceImpl interface.
func (s *VersionServiceImpl) AndroidGetVersion(ctx context.Context, req *version.AndroidGetVersioneRequest) (
	resp *version.AndroidGetVersionResponse, err error,
) {
	resp = new(version.AndroidGetVersionResponse)
	// Android 接口一次返回 release 和 beta 两份数据，使用独立 key 避免和单独版本查询混用。
	v, err := s.singleflight.Do("android_version", func() (any, error) {
		r, b, err := service.NewVersionService(ctx, s.ClientSet).AndroidGetVersion()
		if err != nil {
			return androidVersionResult{}, err
		}
		return androidVersionResult{release: r, beta: b}, nil
	})
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Infof("Version.AndroidGetVersion: %v", err)
		return resp, nil
	}
	result, ok := v.(androidVersionResult)
	if !ok {
		resp.Base = base.BuildBaseResp(singleflight.ErrInvalidType)
		return resp, nil
	}
	resp.Release = pack.BuildVersion(result.release)
	resp.Beta = pack.BuildVersion(result.beta)
	return resp, err
}
