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

package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (c *DBUser) GetRelationByUserId(ctx context.Context, followerId, followedId string) (bool, *model.FollowRelation, error) {
	relationModel := new(model.FollowRelation)
	if err := c.client.WithContext(ctx).Table(constants.UserRelationTableName).
		Where("follower_id = ? and followed_id = ?", followerId, followedId).
		Where("status = ?", constants.RelationOKStatus).First(relationModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		logger.Errorf("dal.GetRelationByUserId error:%v", err)
		return false, nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetRelationByUserId error:%v", err)
	}
	return true, relationModel, nil
}

func (c *DBUser) GetUserFriendsId(ctx context.Context, stuId string) (friendsId []string, err error) {
	err = c.client.WithContext(ctx).
		Table(constants.UserRelationTableName).
		Where("follower_id = ? and status = ?", stuId, constants.RelationOKStatus).
		Pluck("followed_id", &friendsId).
		Error
	if err != nil {
		logger.Errorf("dal.GetUserFriendsId error: %v", err)
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetUserFriendsId error: %v", err)
	}
	return friendsId, err
}

func (c *DBUser) GetUserFriendListLength(ctx context.Context, stuId string) (length int64, err error) {
	err = c.client.WithContext(ctx).
		Table(constants.UserRelationTableName).
		Where("follower_id = ? and status = ?", stuId, constants.RelationOKStatus).
		Count(&length).
		Error
	if err != nil {
		logger.Errorf("dal.GetUserFriendListLength error: %v", err)
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetUserFriendListLength error: %v", err)
	}
	return length, nil
}
