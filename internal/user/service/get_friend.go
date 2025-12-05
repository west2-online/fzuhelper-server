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

	db "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *UserService) GetFriendList(stuId string) ([]*db.Student, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	exist := s.cache.IsKeyExist(s.ctx, userFriendKey)
	var friendId []string
	var err error
	if exist {
		friendId, err = s.cache.User.GetUserFriendCache(s.ctx, userFriendKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetUserFriendCache: %w", err)
		}
	} else {
		if friendId, err = s.db.User.GetUserFriendsId(s.ctx, stuId); err != nil {
			return nil, fmt.Errorf("service.GetUserFriendsIdDB: %w", err)
		}
		go func() {
			err := s.cache.User.SetUserFriendListCache(s.ctx, stuId, friendId)
			if err != nil {
				logger.Errorf("service. SetUserFriendListCache: %v", err)
			}
		}()
	}
	// 考虑到我们现在没有传入StuId返回StuInfo的接口。这边将好友信息查完后返回
	friendList := make([]*db.Student, 0, len(friendId))
	for _, id := range friendId {
		existCache := s.cache.IsKeyExist(s.ctx, id)
		if existCache {
			stuInfo, err := s.cache.User.GetStuInfoCache(s.ctx, id)
			if err != nil {
				return nil, fmt.Errorf("service.GetFriendList: %w", err)
			}
			friendList = append(friendList, stuInfo)
			continue
		}
		// 查询数据库是否存入此学生信息
		stuExist, stuInfo, err := s.db.User.GetStudentById(s.ctx, id)
		if err != nil {
			return nil, fmt.Errorf("service.GetFriendList: %w", err)
		}
		if !stuExist { // 如果数据库也没有该学生信息 则只能模糊返回了
			friendList = append(friendList, &db.Student{StuId: id})
			continue
		}
		friendList = append(friendList, stuInfo)
	}
	return friendList, nil
}
