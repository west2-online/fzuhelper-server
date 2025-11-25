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

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *UserService) BindInvitation(stuId, code string) error {
	mapKey := fmt.Sprintf("code_mapping:%s", code)
	exist := s.cache.IsKeyExist(s.ctx, mapKey)
	if !exist {
		return fmt.Errorf("service.BindInvitation: Invalid InvitationCode")
	}
	friendId, err := s.cache.User.GetCodeStuIdMappingCache(s.ctx, mapKey)
	if err != nil {
		return fmt.Errorf("service.GetCodeStuIdMappingCode: %w", err)
	}
	if friendId == stuId {
		return fmt.Errorf("service.BindInvitation: cannot add yourself as friend")
	}
	// 查找是否关系已经存在
	ok, _, err := s.db.User.GetRelationByUserId(s.ctx, stuId, friendId)
	if err != nil {
		return fmt.Errorf("service.GetRelationByUserId: %w", err)
	}
	if ok {
		return fmt.Errorf("service.BindInvitation: RelationShip Already Exist")
	}
	// 好友列表限制
	confine, err := s.IsFriendNumsConfined(stuId)
	if err != nil {
		return err
	}
	if confine {
		return fmt.Errorf("service.BindInvitation :%v friendList is full", stuId)
	}
	targetConfine, err := s.IsFriendNumsConfined(friendId)
	if err != nil {
		return err
	}
	if targetConfine {
		return fmt.Errorf("service.BindInvitation :%v friendList is full", friendId)
	}

	err = s.db.User.CreateRelation(s.ctx, stuId, friendId)
	if err != nil {
		return fmt.Errorf("service.CreateRelation: %w", err)
	}
	go func() {
		// 目前绑定成功插入双向关系
		err = s.cache.User.SetUserFriendCache(s.ctx, friendId, stuId)
		if err != nil {
			logger.Errorf("service. SetUserFriendCache: %v", err)
		}
		err = s.cache.User.RemoveCodeStuIdMappingCache(s.ctx, mapKey) // 如果邀请码设为一次性
		if err != nil {
			logger.Errorf("service. RemoveCodeStuIdMappingCache: %v", err)
		}
	}()
	return nil
}

func (s *UserService) IsFriendNumsConfined(stuId string) (bool, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	exist := s.cache.IsKeyExist(s.ctx, userFriendKey)
	if exist {
		friends, err := s.cache.User.GetUserFriendCache(s.ctx, userFriendKey)
		if err != nil {
			return false, fmt.Errorf("service.IsFriendNumsConfined get user friend cache: %w", err)
		}
		if int64(len(friends)) >= config.Friend.MaxNum {
			return true, nil
		}
		return false, nil
	} else {
		length, err := s.db.User.GetUserFriendListLength(s.ctx, stuId)
		if err != nil {
			return false, fmt.Errorf("service.IsFriendNumsConfined get user friend length db: %w", err)
		}
		if length >= config.Friend.MaxNum {
			return true, nil
		}
		return false, nil
	}
}
