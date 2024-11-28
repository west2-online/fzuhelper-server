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
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/url"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// UploadParams 实际上是获取上传参数给前端使用
func (s *UrlService) UploadParams(req *url.UploadParamsRequest) (string, string, error) {
	if !utils.CheckPwd(req.Password) {
		return "", "", buildAuthFailedError()
	}
	policy := upyun.GetPolicy(config.UrlService.Bucket, config.UrlService.Path, int(config.UrlService.TokenTimeout))
	authorization := upyun.SignStr(config.UrlService.Operator, config.UrlService.Pass, config.UrlService.Bucket, policy)
	return policy, authorization, nil
}
