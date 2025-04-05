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

package oss

import (
	"github.com/upyun/go-sdk/v3/upyun"

	"github.com/west2-online/fzuhelper-server/config"
)

type OSSSet struct {
	Provider string // 供应商

	Upyun *UpYunConfig
}

type UpYunConfig struct {
	upyun          *upyun.UpYun
	TokenSecret    string
	TokenTimeout   int64
	UssDomain      string
	DownloadDomain string
	Path           string
	AvatarPath     string
}

func NewUpYunConfig() *UpYunConfig {
	return &UpYunConfig{
		upyun: upyun.NewUpYun(
			&upyun.UpYunConfig{
				Bucket:   config.UpYun.Bucket,
				Operator: config.UpYun.Operator,
				Password: config.UpYun.Password,
			},
		),
		TokenSecret:    config.UpYun.TokenSecret,
		TokenTimeout:   config.UpYun.TokenTimeout,
		UssDomain:      config.UpYun.UssDomain,
		DownloadDomain: config.UpYun.DownloadDomain,
		Path:           config.UpYun.Path,
		AvatarPath:     config.UpYun.AvatarPath,
	}
}
