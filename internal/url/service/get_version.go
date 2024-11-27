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

	"github.com/west2-online/fzuhelper-server/internal/url/pack"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *UrlService) GetReleaseVersion() (*pack.Version, error) {
	jsonBytes, err := utils.GetJSON(constants.StatisticPath + releaseVersionFileName)
	if err != nil {
		return nil, fmt.Errorf("UrlService.GetReleaseVersion error:%w", err)
	}
	version := new(pack.Version)
	err = sonic.Unmarshal(jsonBytes, version)
	if err != nil {
		return nil, fmt.Errorf("UrlService.GetReleaseVersion error:%w", err)
	}
	return version, nil
}

func (s *UrlService) GetBetaVersion() (*pack.Version, error) {
	jsonBytes, err := utils.GetJSON(constants.StatisticPath + betaVersionFileName)
	if err != nil {
		return nil, fmt.Errorf("UrlService.GetBetaVersion error:%w", err)
	}
	version := new(pack.Version)
	err = sonic.Unmarshal(jsonBytes, version)
	if err != nil {
		return nil, fmt.Errorf("UrlService.GetBetaVersion error:%w", err)
	}
	return version, nil
}
