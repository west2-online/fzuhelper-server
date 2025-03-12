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
	"errors"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestCourseService_GetTermsList(t *testing.T) {
	type testCase struct {
		name            string
		mockTermsReturn *jwch.Term
		mockTermsError  error
		expectResult    []string
		expectError     error
		cacheExist      bool
		cacheGetError   error
	}
	successTerm := &jwch.Term{
		Terms:           []string{"202401"},
		ViewState:       "viewstate123",
		EventValidation: "eventvalidation123",
	}
	testCases := []testCase{
		{
			name:            "Success",
			expectResult:    successTerm.Terms,
			expectError:     nil,
			mockTermsReturn: successTerm,
			mockTermsError:  nil,
		},
		{
			name:          "cache exist success",
			cacheExist:    true, // 缓存里已存在
			cacheGetError: nil,  // 获取缓存不报错
			expectResult:  successTerm.Terms,
		},
		{
			name:          "cache exist but get cache error",
			cacheExist:    true,
			cacheGetError: fmt.Errorf("redis get error"),
			expectError:   errors.New("redis get error"),
		},
	}
	mockLoginData := &model.LoginData{
		Id:      "052106112",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()
	mockey.Mock((*coursecache.CacheCourse).SetTermsCache).To(
		func(ctx context.Context, key string, list []string) error {
			return nil
		},
	).Build()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*jwch.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()
			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						if tc.cacheGetError != nil {
							return nil, tc.cacheGetError
						}
						return successTerm.Terms, nil
					},
				).Build()
			} else {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						return nil, fmt.Errorf("should not be called if cache doesn't exist")
					},
				).Build()
			}

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, nil)

			result, err := courseService.GetTermsList(&model.LoginData{
				Id:      "123456789",
				Cookies: "magic cookies",
			})

			if tc.expectError != nil {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
