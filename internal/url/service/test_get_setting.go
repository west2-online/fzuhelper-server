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

	"github.com/west2-online/fzuhelper-server/internal/url/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/url"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *UrlService) TestSetting(req *url.GetTestRequest) (*[]byte, error) {
	// 获得Json
	settingJson, err := utils.GetJSON(constants.StatisticPath + cloudSettingFileName)
	if err != nil {
		return nil, fmt.Errorf("UrlService.TestSetting error:%w", err)
	}
	noCommentSettingJson, err := getJSONWithoutComments(string(settingJson))
	if err != nil {
		return nil, fmt.Errorf("UrlService.TestSetting error:%w", err)
	}

	// 绑定结构体
	cloudSettings := new(pack.CloudSetting)
	err = json.Unmarshal([]byte(noCommentSettingJson), cloudSettings)
	if err != nil {
		return nil, fmt.Errorf("UrlService.TestSetting error:%w", err)
	}

	criteria := &pack.Plan{
		Account:   req.Account,
		Version:   req.Version,
		Beta:      req.Beta,
		Phone:     req.Phone,
		IsLogin:   req.IsLogin,
		LoginType: req.LoginType,
	}
	plan, err := findMatchingPlan(&cloudSettings.Plans, criteria)
	if err != nil {
		return nil, fmt.Errorf("UrlService.TestSetting error:%w", err)
	}
	returnPlan := []byte(plan.Plan)
	return &returnPlan, nil
}
