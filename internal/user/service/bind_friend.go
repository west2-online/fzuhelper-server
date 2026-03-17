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

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *UserService) BindInvitation(stuId, code string) error {
	mapKey := fmt.Sprintf("code_mapping:%s", code)
	exist := s.cache.IsKeyExist(s.ctx, mapKey)
	if !exist {
		return fmt.Errorf("无效邀请码")
	}
	friendId, err := s.cache.User.GetCodeStuIdMappingCache(s.ctx, mapKey)
	if err != nil {
		return fmt.Errorf("service.GetCodeStuIdMappingCode: %w", err)
	}
	if friendId == stuId {
		return fmt.Errorf("无法添加自己为好友")
	}
	// 查找是否关系已经存在
	ok, _, err := s.db.User.GetRelationByUserId(s.ctx, stuId, friendId)
	if err != nil {
		return fmt.Errorf("service.GetRelationByUserId: %w", err)
	}
	if ok {
		return fmt.Errorf("好友关系已存在")
	}
	// 好友列表限制
	maxNum := s.GetFriendMaxNum(stuId)
	confine, err := s.IsFriendNumsConfined(stuId, maxNum)
	if err != nil {
		return err
	}
	if confine {
		return fmt.Errorf("您的好友列表已满，最多拥有 %v 名好友",
			maxNum)
	}
	targetMaxNum := s.GetFriendMaxNum(friendId)
	targetConfine, err := s.IsFriendNumsConfined(friendId, targetMaxNum)
	if err != nil {
		return err
	}
	if targetConfine {
		return fmt.Errorf("对方好友列表已满，最多拥有 %v 名好友",
			targetMaxNum)
	}

	err = s.writeRelationToDB(stuId, friendId)
	if err != nil {
		return fmt.Errorf("service.CreateRelation: %w", err)
	}

	// 同步清除缓存，避免 DB 写入后客户端立即查询仍读到旧数据
	codeKey := fmt.Sprintf("codes:%s", friendId)
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	targetFriendKey := fmt.Sprintf("user_friends:%v", friendId)

	if s.cache.IsKeyExist(s.ctx, userFriendKey) {
		if err = s.cache.User.InvalidateFriendListCache(s.ctx, stuId); err != nil {
			logger.Errorf("service.InvalidateFriendListCache: %v", err)
		}
	}
	if s.cache.IsKeyExist(s.ctx, targetFriendKey) {
		if err = s.cache.User.InvalidateFriendListCache(s.ctx, friendId); err != nil {
			logger.Errorf("service.InvalidateFriendListCache: %v", err)
		}
	}
	if err = s.cache.User.RemoveCodeStuIdMappingCache(s.ctx, mapKey); err != nil {
		logger.Errorf("service.RemoveCodeStuIdMappingCache: %v", err)
	}
	if err = s.cache.User.RemoveInvitationCodeCache(s.ctx, codeKey); err != nil {
		logger.Errorf("service.RemoveInvitationCodeCache: %v", err)
	}
	return nil
}

func (s *UserService) IsFriendNumsConfined(stuId string, maxNum int64) (bool, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	exist := s.cache.IsKeyExist(s.ctx, userFriendKey)
	if exist {
		friends, err := s.cache.User.GetUserFriendCache(s.ctx, userFriendKey)
		if err != nil {
			return false, fmt.Errorf("service.IsFriendNumsConfined get user friend cache: %w", err)
		}
		if int64(len(friends)) >= maxNum {
			return true, nil
		}
		return false, nil
	} else {
		length, err := s.db.User.GetUserFriendListLength(s.ctx, stuId)
		if err != nil {
			return false, fmt.Errorf("service.IsFriendNumsConfined get user friend length db: %w", err)
		}
		if length >= maxNum {
			return true, nil
		}
		return false, nil
	}
}

func (s *UserService) writeRelationToDB(followedId, followerId string) error {
	var relation []*model.FollowRelation
	dbId, err := s.sf.NextVal()
	if err != nil {
		return err
	}
	relation = append(relation, &model.FollowRelation{
		Id:         dbId,
		FollowedId: followedId,
		FollowerId: followerId,
	})
	dbId, err = s.sf.NextVal()
	if err != nil {
		return err
	}
	relation = append(relation, &model.FollowRelation{
		Id:         dbId,
		FollowedId: followerId,
		FollowerId: followedId,
	})
	return s.db.User.CreateRelation(s.ctx, relation)
}
