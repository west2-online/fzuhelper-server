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
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *UserService) DeleteUserFriend(loginData *loginmodel.LoginData, targetStuId string) error {
	stuId := context.ExtractIDFromLoginData(loginData)
	ok, _, err := s.db.User.GetRelationByUserId(s.ctx, stuId, targetStuId)
	if err != nil {
		return fmt.Errorf("service.GetRelationByUserId: %w", err)
	}
	if !ok {
		return fmt.Errorf("service.DeleteUserFriend: RelationShip No Exist")
	}
	if err = s.db.User.DeleteRelation(s.ctx, stuId, targetStuId); err != nil {
		return fmt.Errorf("service.DeleteRelation: %w", err)
	}
	go func() {
		err := s.cache.User.DeleteUserFriendCache(s.ctx, stuId, targetStuId)
		if err != nil {
			logger.Errorf("service. DeleteUserFriendCache: %v", err)
		}
	}()
	return nil
}
