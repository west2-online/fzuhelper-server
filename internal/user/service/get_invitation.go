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

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *UserService) GetInvitationCode(stuId string, isRefresh bool) (code string, err error) {
	codeKey := fmt.Sprintf("codes:%s", stuId)
	exist := s.cache.IsKeyExist(s.ctx, codeKey)
	// 存在返回cache内已经生成的邀请码
	if exist {
		code, err = s.cache.User.GetInvitationCodeCache(s.ctx, codeKey)
		if err != nil {
			return "", fmt.Errorf("service.GetInvitationCode: %w", err)
		}
		if !isRefresh {
			return code, nil
		}
	}
	newCode := utils.GenerateRandomCode(constants.CommonInvitationCodeLength)
	mapKey := fmt.Sprintf("code_mapping:%s", code)
	newMapKey := fmt.Sprintf("code_mapping:%s", newCode)
	go func() {
		err := s.cache.User.RemoveCodeStuIdMappingCache(s.ctx, mapKey)
		if err != nil {
			logger.Errorf("service. RemoveCodeStuIdMappingCache: %v", err)
		}
		err = s.cache.User.SetInvitationCodeCache(s.ctx, codeKey, newCode)
		if err != nil {
			logger.Errorf("service. SetInvitationCode: %v", err)
		}
		err = s.cache.User.SetCodeStuIdMappingCache(s.ctx, newMapKey, stuId)
		if err != nil {
			logger.Errorf("service. SetCodeStuIdMappingCache: %v", err)
		}
	}()
	return newCode, nil
}
