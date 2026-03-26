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

package friend_config

import (
	"context"
	"reflect"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBFriendConfig_GetFriendConfigs(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		mockConfigs    []*model.FriendConfig
		expectingError bool
		errorMsg       string
	}

	testCases := []testCase{
		{
			name:      "success_with_configs",
			mockError: nil,
			mockConfigs: []*model.FriendConfig{
				{
					Id:        1,
					ConfigKey: "max_num",
					Value:     "3",
					StudentID: "",
				},
				{
					Id:        2,
					ConfigKey: "max_num",
					Value:     "10",
					StudentID: "102300217",
				},
			},
			expectingError: false,
		},
		{
			name:           "success_with_empty_configs",
			mockError:      nil,
			mockConfigs:    []*model.FriendConfig{},
			expectingError: false,
		},
		{
			name:           "database_error",
			mockError:      gorm.ErrInvalidDB,
			mockConfigs:    nil,
			expectingError: true,
			errorMsg:       "dal.GetFriendConfigs error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBFriendConfig := NewDBFriendConfig(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}

				if tc.mockConfigs != nil {
					destValue := reflect.ValueOf(dest)
					if destValue.Kind() == reflect.Ptr {
						elem := destValue.Elem()
						mockValue := reflect.ValueOf(tc.mockConfigs)
						elem.Set(mockValue)
					}
				}

				return mockGormDB
			}).Build()

			configs, err := mockDBFriendConfig.GetFriendConfigs(context.Background())

			if tc.expectingError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
				assert.Nil(t, configs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.mockConfigs), len(configs))
				for i, cfg := range configs {
					assert.Equal(t, tc.mockConfigs[i].Id, cfg.Id)
					assert.Equal(t, tc.mockConfigs[i].ConfigKey, cfg.ConfigKey)
					assert.Equal(t, tc.mockConfigs[i].Value, cfg.Value)
					assert.Equal(t, tc.mockConfigs[i].StudentID, cfg.StudentID)
				}
			}
		})
	}
}
