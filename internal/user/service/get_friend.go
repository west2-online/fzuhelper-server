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

	"github.com/west2-online/fzuhelper-server/internal/user/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	db "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *UserService) GetFriendList(stuId string) ([]*model.UserFriendInfo, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	exist := s.cache.IsKeyExist(s.ctx, userFriendKey)
	var friendRelation []*db.UserFriend
	var err error
	if exist {
		friendRelation, err = s.cache.User.GetUserFriendCache(s.ctx, userFriendKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetUserFriendCache: %w", err)
		}
	} else {
		if friendRelation, err = s.db.User.GetUserFriends(s.ctx, stuId); err != nil {
			return nil, fmt.Errorf("service.GetUserFriendsIdDB: %w", err)
		}
		go func() {
			err := s.cache.User.SetUserFriendListCache(s.ctx, stuId, friendRelation)
			if err != nil {
				logger.Errorf("service. SetUserFriendListCache: %v", err)
			}
		}()
	}
	friendList := make([]*model.UserFriendInfo, 0, len(friendRelation))
	for _, relation := range friendRelation {
		if s.cache.IsKeyExist(s.ctx, relation.FriendId) {
			stuInfo, err := s.cache.User.GetStuInfoCache(s.ctx, relation.FriendId)
			if err != nil {
				return nil, fmt.Errorf("service.GetFriendList: %w", err)
			}
			friendList = append(friendList, pack.BuildFriendInfoResp(stuInfo, relation))
			continue
		}
		// 查询数据库是否存入此学生信息
		stuExist, stuInfo, err := s.db.User.GetStudentById(s.ctx, relation.FriendId)
		if err != nil {
			return nil, fmt.Errorf("service.GetFriendList: %w", err)
		}
		if !stuExist { // 如果数据库也没有该学生信息 则只能模糊返回了
			friendList = append(friendList, &model.UserFriendInfo{
				StuId:     relation.FriendId,
				CreatedAt: relation.UpdatedAt.Unix(),
			})
			continue
		}
		friendList = append(friendList, pack.BuildFriendInfoResp(stuInfo, relation))
	}
	return friendList, nil
}
