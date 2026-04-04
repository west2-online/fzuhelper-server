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

	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *UserService) ReorderFriendList(loginData *loginmodel.LoginData, friendIds []string) error {
	stuId := context.ExtractIDFromLoginData(loginData)
	if err := s.db.User.ReorderFriendList(s.ctx, stuId, friendIds); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "User.ReorderFriendList: %v", err)
	}

	// 删除好友列表缓存
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	if s.cache.IsKeyExist(s.ctx, userFriendKey) {
		if err := s.cache.User.InvalidateFriendListCache(s.ctx, stuId); err != nil {
			logger.Errorf("service.ReorderFriendList: delete cache failed: %v", err)
		}
	}

	return nil
}
