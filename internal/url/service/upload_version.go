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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/west2-online/fzuhelper-server/internal/url/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/url"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var mu = &sync.Mutex{}

func (s *UrlService) UploadVersion(req *url.UploadRequest) error {
	if !utils.CheckPwd(req.Password) {
		return buildAuthFailedError()
	}
	version := &pack.Version{
		Version: req.Version,
		Code:    req.Code,
		Url:     req.Url,
		Feature: req.Feature,
	}
	jsonBytes, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("UrlService.UploadVersion json marshal err: %w", err)
	}

	switch req.Type {
	case apkTypeRelease:
		mu.Lock()
		defer mu.Unlock()

		return utils.SaveJSON(constants.StatisticPath+releaseVersionFileName, jsonBytes)
	case apkTypeBeta:
		mu.Lock()
		defer mu.Unlock()

		return utils.SaveJSON(constants.StatisticPath+betaVersionFileName, jsonBytes)
	default:
		return errno.ParamError
	}
}
