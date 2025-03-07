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

package service

import (
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/internal/version/pack"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func (s *VersionService) DownloadReleaseApk() (string, error) {
	jsonBytes, err := upyun.URlGetFile(upyun.JoinFileName(releaseVersionFileName))
	if err != nil {
		return "", fmt.Errorf("VersionService.DownloadReleaseApk error:%w", err)
	}
	version := new(pack.Version)
	err = sonic.Unmarshal(*jsonBytes, version)
	if err != nil {
		return "", fmt.Errorf("VersionService.DownloadReleaseApk error:%w", err)
	}
	return version.Url, nil
}

func (s *VersionService) DownloadBetaApk() (string, error) {
	jsonBytes, err := upyun.URlGetFile(upyun.JoinFileName(betaVersionFileName))
	if err != nil {
		return "", fmt.Errorf("VersionService.DownloadBetaApk error:%w", err)
	}
	version := new(pack.Version)
	err = sonic.Unmarshal(*jsonBytes, version)
	if err != nil {
		return "", fmt.Errorf("VersionService.DownloadBetaApk error:%w", err)
	}
	return version.Url, nil
}
