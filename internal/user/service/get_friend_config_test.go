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
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	friendConfigDB "github.com/west2-online/fzuhelper-server/pkg/db/friend_config"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUserService_GetFriendMaxNum(t *testing.T) {
	// 初始化配置，读取 config.example.yaml 中的 friend.max-nums 作为回退值
	err := config.InitForTest("user")
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	type testCase struct {
		name         string
		mockConfigs  []*dbmodel.FriendConfig
		mockError    error
		expectResult int64
	}

	// config.Friend.MaxNum 的默认值来自 config.example.yaml (max-nums: 3)
	fallbackMaxNum := config.Friend.MaxNum
	// mock login data
	loginData := &model.LoginData{
		Id:      "102300217",
		Cookies: "cookie1=value1;cookie2=value2",
	}

	testCases := []testCase{
		{
			name: "global config only",
			mockConfigs: []*dbmodel.FriendConfig{
				{
					ConfigKey: "max_num",
					Value:     "5",
					StudentID: "",
				},
			},
			mockError:    nil,
			expectResult: 5,
		},
		{
			name: "student specific config",
			mockConfigs: []*dbmodel.FriendConfig{
				{
					ConfigKey: "max_num",
					Value:     "5",
					StudentID: "",
				},
				{
					ConfigKey: "max_num",
					Value:     "10",
					StudentID: "102300217",
				},
			},
			mockError:    nil,
			expectResult: 10,
		},
		{
			name: "student specific config for different student",
			mockConfigs: []*dbmodel.FriendConfig{
				{
					ConfigKey: "max_num",
					Value:     "5",
					StudentID: "",
				},
				{
					ConfigKey: "max_num",
					Value:     "10",
					StudentID: "102300218",
				},
			},
			mockError:    nil,
			expectResult: 5,
		},
		{
			name:         "no config found fallback to yaml",
			mockConfigs:  []*dbmodel.FriendConfig{},
			mockError:    nil,
			expectResult: fallbackMaxNum,
		},
		{
			name:         "db error fallback to yaml",
			mockConfigs:  nil,
			mockError:    fmt.Errorf("database error"),
			expectResult: fallbackMaxNum,
		},
		{
			name: "invalid value string fallback to yaml",
			mockConfigs: []*dbmodel.FriendConfig{
				{
					ConfigKey: "max_num",
					Value:     "invalid",
					StudentID: "",
				},
			},
			mockError:    nil,
			expectResult: fallbackMaxNum,
		},
		{
			name: "negative value fallback to yaml",
			mockConfigs: []*dbmodel.FriendConfig{
				{
					ConfigKey: "max_num",
					Value:     "-1",
					StudentID: "",
				},
			},
			mockError:    nil,
			expectResult: fallbackMaxNum,
		},
		{
			name: "zero value fallback to yaml",
			mockConfigs: []*dbmodel.FriendConfig{
				{
					ConfigKey: "max_num",
					Value:     "0",
					StudentID: "",
				},
			},
			mockError:    nil,
			expectResult: fallbackMaxNum,
		},
		{
			name: "irrelevant config key ignored",
			mockConfigs: []*dbmodel.FriendConfig{
				{
					ConfigKey: "other_key",
					Value:     "100",
					StudentID: "",
				},
			},
			mockError:    nil,
			expectResult: fallbackMaxNum,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: &db.Database{
					FriendConfig: &friendConfigDB.DBFriendConfig{},
				},
			}
			userService := NewUserService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))

			mockey.Mock((*friendConfigDB.DBFriendConfig).GetFriendConfigs).Return(tc.mockConfigs, tc.mockError).Build()

			result := userService.GetFriendMaxNum(loginData)
			assert.Equal(t, tc.expectResult, result)
		})
	}
}
