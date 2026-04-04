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

	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *UserService) CancelInvitationCode(loginData *kitexModel.LoginData) error {
	stuId := context.ExtractIDFromLoginData(loginData)
	codeKey := fmt.Sprintf("codes:%s", stuId)
	if !s.cache.IsKeyExist(s.ctx, codeKey) {
		return errno.Errorf(errno.InternalRedisErrorCode, "当前账号暂无处于生效状态的邀请码")
	}
	code, _, err := s.cache.User.GetInvitationCodeCache(s.ctx, codeKey)
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "User.GetInvitationCodeCache: %v", err)
	}
	mapKey := fmt.Sprintf("code_mapping:%s", code)
	Id, err := s.cache.User.GetCodeStuIdMappingCache(s.ctx, mapKey)
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "User.GetCodeStuIdMappingCodeCache: %v", err)
	}
	// 删除cache中存储的邀请码及映射关系
	if Id == stuId {
		err = s.cache.User.RemoveCodeStuIdMappingCache(s.ctx, mapKey)
		if err != nil {
			return errno.Errorf(errno.InternalRedisErrorCode, "User.RemoveCodeStuIdMappingCache: %v", err)
		}
		err = s.cache.User.RemoveInvitationCodeCache(s.ctx, codeKey)
		if err != nil {
			return errno.Errorf(errno.InternalRedisErrorCode, "User.RemoveInvitationCodeCache: %v", err)
		}
	}
	return nil
}
