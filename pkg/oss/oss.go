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
	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/cos"
)

type OSSSet struct {
	Provider string // 供应商

	Cos *CosConfig
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
