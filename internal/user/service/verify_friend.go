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

	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *UserService) VerifyUserFriend(stuId string, friendId string) (bool, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	exist := s.cache.IsKeyExist(s.ctx, userFriendKey)
	ok := false
	var err error
	// 验证好友
	if exist {
		ok, err = s.cache.User.IsFriendCache(s.ctx, stuId, friendId)
		if err != nil {
			return false, errno.Errorf(errno.InternalDatabaseErrorCode, "User.VerifyUserFriend: Get friend cache fail: %v", err)
		}
	} else {
		ok, _, err = s.db.User.GetRelationByUserId(s.ctx, stuId, friendId)
		if err != nil {
			return false, errno.Errorf(errno.InternalDatabaseErrorCode, "User.VerifyUserFriend: Get friend db fail: %v", err)
		}
	}
	return ok, nil
}
