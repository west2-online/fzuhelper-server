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
	"strconv"

	"github.com/west2-online/fzuhelper-server/config"
	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// GetFriendMaxNum 从 friend_config 表获取好友数量上限
// 支持按 student_id 维度进行精细化匹配（学号精确匹配优先于全局配置）
// 若无匹配记录，回退到 config.yaml 中的 friend.max-nums 配置
func (s *UserService) GetFriendMaxNum(loginData *loginmodel.LoginData) int64 {
	stuId := context.ExtractIDFromLoginData(loginData)
	configs, err := s.db.FriendConfig.GetFriendConfigs(s.ctx)
	if err != nil {
		logger.Errorf("service.GetFriendMaxNum: get friend configs error: %v, fallback to config", err)
		return config.Friend.MaxNum
	}

	var bestMatch *model.FriendConfig

	for _, cfg := range configs {
		if cfg.ConfigKey != constants.FriendConfigKeyMaxNum {
			continue
		}

		// 精确学号匹配优先级最高，直接返回
		if cfg.StudentID != "" && cfg.StudentID == stuId {
			bestMatch = cfg
			break
		}

		// 全局配置（student_id 为空）作为兜底
		if cfg.StudentID == "" && bestMatch == nil {
			bestMatch = cfg
		}
	}

	if bestMatch == nil {
		return config.Friend.MaxNum
	}

	maxNum, err := strconv.ParseInt(bestMatch.Value, 10, 64)
	if err != nil {
		logger.Errorf("service.GetFriendMaxNum: parse value '%s' error: %v, fallback to config", bestMatch.Value, err)
		return config.Friend.MaxNum
	}

	if maxNum <= 0 {
		logger.Errorf("service.GetFriendMaxNum: invalid value %d, fallback to config", maxNum)
		return config.Friend.MaxNum
	}

	return maxNum
}
