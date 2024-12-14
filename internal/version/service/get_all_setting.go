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

	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func (s *VersionService) GetAllCloudSetting() (*[]byte, error) {
	// 获得Json
	settingJson, err := upyun.URlGetFile(upyun.JoinFileName(cloudSettingFileName))
	if err != nil {
		return nil, fmt.Errorf("VersionService.GetAllCloudSetting error:%w", err)
	}
	noCommentSettingJson, err := getJSONWithoutComments(string(*settingJson))
	if err != nil {
		return nil, fmt.Errorf("VersionService.GetAllCloudSetting error:%w", err)
	}

	returnPlan := []byte(noCommentSettingJson)

	return &returnPlan, nil
}
