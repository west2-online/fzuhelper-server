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
	"sync"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestGetInvitationCode(t *testing.T) {
	type testCase struct {
		name           string
		expectExist    bool
		mockError      error
		expectError    string
		IsRefresh      bool
		cacheExist     bool
		cacheGetError  error
		cacheCode      string
		cacheCreatedAt int64
	}

	stuId := "102300217"

	testCases := []testCase{
		{
			name:           "cache get error",
			IsRefresh:      false,
			expectExist:    true,
			expectError:    "service.GetInvitationCode:",
			mockError:      errno.InternalServiceError,
			cacheExist:     true,
			cacheGetError:  errno.InternalServiceError,
			cacheCode:      "",
			cacheCreatedAt: -1,
		},
		{
			name:           "IsRefresh true - force regenerate",
			IsRefresh:      true,
			expectExist:    true,
			cacheExist:     true,
			cacheGetError:  nil,
			cacheCode:      "123456",
			cacheCreatedAt: 1321045012,
		},
		{
			name:           "cache code exist and no refresh",
			expectExist:    true,
			IsRefresh:      false,
			cacheExist:     true,
			cacheGetError:  nil,
			cacheCode:      "123456",
			cacheCreatedAt: 1321045012,
		},
		{
			name:        "cache not exist and refresh true",
			IsRefresh:   true,
			expectExist: false,
			cacheExist:  false,
		},
		{
			name:        "cache not exist and refresh false",
			IsRefresh:   false,
			expectExist: false,
			cacheExist:  false,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			shouldWait := !tc.cacheExist || tc.IsRefresh
			var wg sync.WaitGroup
			if shouldWait {
				wg.Add(3)
			}

			mockClientSet := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			setInvitationGuard := mockey.Mock((*user.CacheUser).SetInvitationCodeCache).To(func(ctx context.Context, key string, code string) error {
				if shouldWait {
					wg.Done()
				}
				return nil
			}).Build()
			defer setInvitationGuard.UnPatch()

			setMappingGuard := mockey.Mock((*user.CacheUser).SetCodeStuIdMappingCache).To(func(ctx context.Context, key string, stuId string) error {
				if shouldWait {
					wg.Done()
				}
				return nil
			}).Build()
			defer setMappingGuard.UnPatch()

			removeMappingGuard := mockey.Mock((*user.CacheUser).RemoveCodeStuIdMappingCache).To(func(ctx context.Context, key string) error {
				if shouldWait {
					wg.Done()
				}
				return nil
			}).Build()
			defer removeMappingGuard.UnPatch()
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			mockey.Mock((*user.CacheUser).GetInvitationCodeCache).Return(tc.cacheCode, tc.cacheCreatedAt, tc.cacheGetError).Build()

			if !tc.cacheExist || tc.IsRefresh {
				mockey.Mock(utils.GenerateRandomCode).Return("ABCDEF").Build()
			}

			code, expireAt, err := userService.GetInvitationCode(stuId, tc.IsRefresh)
			if shouldWait && err == nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Fatalf("async cache update did not finish in time")
				}
			}
			if tc.expectError != "" {
				assert.Equal(t, "", code)
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				if tc.cacheExist && !tc.IsRefresh && tc.cacheGetError == nil {
					assert.Equal(t, tc.cacheCode, code)
					assert.Equal(t, tc.cacheCreatedAt+int64(constants.UserInvitationCodeKeyExpire/time.Second),
						expireAt)
				}
				if !tc.cacheExist || tc.IsRefresh {
					assert.Equal(t, 6, len(code))
					assert.Equal(t, "ABCDEF", code)
				}
			}
		})
	}
}
