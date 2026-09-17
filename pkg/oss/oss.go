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
	"fmt"

	tencentyun "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/upyun/go-sdk/v3/upyun"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/cos"
)

type OSSSet struct {
	Provider string // 供应商

	Cos   *CosConfig
	Upyun *UpYunConfig
}

type CosConfig struct {
	client         *tencentyun.Client
	TokenSecret    string
	TokenTimeout   int64
	DownloadDomain string
	Path           string
	AvatarPath     string
}

func NewCosConfig() *CosConfig {
	return &CosConfig{
		client:         cos.NewCos(),
		TokenSecret:    config.Cos.TokenSecret,
		TokenTimeout:   config.Cos.TokenTimeout,
		DownloadDomain: config.Cos.DownloadDomain,
		Path:           config.Cos.Path,
		AvatarPath:     config.Cos.AvatarPath,
	}
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

func NewUpYunConfig() (*UpYunConfig, error) {
	if config.UpYun == nil {
		return nil, fmt.Errorf("upyun config for current service is missing")
	}
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
	}, nil
}
