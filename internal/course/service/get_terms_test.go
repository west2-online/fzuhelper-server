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

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbcourse "github.com/west2-online/fzuhelper-server/pkg/db/course"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func TestGetTermsList(t *testing.T) {
	type testCase struct {
		name            string
		mockTermsReturn *jwch.Term
		mockTermsError  error
		expectResult    []string
		expectError     string
		cacheExist      bool
		cacheGetError   error
		mockSetCacheErr error
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
			mockTermsReturn: successTerm,
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
			cacheGetError: assert.AnError,
			expectError:   "assert.AnError",
		},
		{
			name:           "cache miss get terms error",
			mockTermsError: assert.AnError,
			expectError:    "Get terms fail",
		},
		{
			name:            "cache miss set cache error",
			mockTermsReturn: successTerm,
			mockSetCacheErr: assert.AnError,
			expectResult:    successTerm.Terms,
		},
	}

	mockLoginData := &model.LoginData{
		Id:      "052106112",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			shouldWait := !tc.cacheExist && tc.mockTermsError == nil
			var wg sync.WaitGroup
			if shouldWait {
				wg.Add(1)
			}

			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*jwch.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()
			setTermsGuard := mockey.Mock((*coursecache.CacheCourse).SetTermsCache).To(func(ctx context.Context, key string, terms []string) error {
				if shouldWait {
					wg.Done()
				}
				return tc.mockSetCacheErr
			}).Build()
			defer setTermsGuard.UnPatch()
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).Return(successTerm.Terms, tc.cacheGetError).Build()
			} else {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).Return(nil, assert.AnError).Build()
			}

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetTermsList(mockLoginData)
			if shouldWait && err == nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Fatalf("async cache set did not finish in time")
				}
			}
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestGetTermsListYjsy(t *testing.T) {
	type testCase struct {
		name            string
		mockTermsReturn *yjsy.Term
		mockTermsError  error
		expectResult    []string
		expectError     string
		cacheExist      bool
		cacheGetError   error
		mockSetCacheErr error
	}

	successTerm := &yjsy.Term{
		Terms:           []string{"202401"},
		ViewState:       "viewstate123",
		EventValidation: "eventvalidation123",
	}

	testCases := []testCase{
		{
			name:            "Success",
			expectResult:    successTerm.Terms,
			mockTermsReturn: successTerm,
		},
		{
			name:         "cache exist success",
			cacheExist:   true,
			expectResult: successTerm.Terms,
		},
		{
			name:          "cache exist but get cache error",
			cacheExist:    true,
			cacheGetError: assert.AnError,
			expectError:   "assert.AnError",
		},
		{
			name:           "cache miss get terms error",
			mockTermsError: assert.AnError,
			expectError:    "Get terms fail",
		},
		{
			name:            "cache miss set cache error",
			mockTermsReturn: successTerm,
			mockSetCacheErr: assert.AnError,
			expectResult:    successTerm.Terms,
		},
	}

	mockLoginData := &model.LoginData{
		Id:      "052106112",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			shouldWait := !tc.cacheExist && tc.mockTermsError == nil
			var wg sync.WaitGroup
			if shouldWait {
				wg.Add(1)
			}

			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*yjsy.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()
			setTermsGuard := mockey.Mock((*coursecache.CacheCourse).SetTermsCache).To(func(ctx context.Context, key string, terms []string) error {
				if shouldWait {
					wg.Done()
				}
				return tc.mockSetCacheErr
			}).Build()
			defer setTermsGuard.UnPatch()
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).Return(successTerm.Terms, tc.cacheGetError).Build()
			} else {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).Return(nil, assert.AnError).Build()
			}

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetTermsListYjsy(mockLoginData)
			if shouldWait && err == nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Fatalf("async cache set did not finish in time")
				}
			}
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestPutTermToDatabase(t *testing.T) {
	type testCase struct {
		name      string
		getReturn *dbmodel.UserTerm
		getErr    error
		nextVal   int64
		nextErr   error
		createErr error
		updateErr error
		termList  string
		expectErr string
	}

	stuId := "stu-1"

	cases := []testCase{
		{
			name:      "db get error",
			getErr:    assert.AnError,
			termList:  "202401",
			expectErr: "assert.AnError",
		},
		{
			name:     "create success",
			nextVal:  10,
			termList: "202401",
		},
		{
			name:      "create nextval error",
			nextErr:   assert.AnError,
			termList:  "202401",
			expectErr: "assert.AnError",
		},
		{
			name:      "create insert error",
			nextVal:   12,
			createErr: assert.AnError,
			termList:  "202401",
			expectErr: "assert.AnError",
		},
		{
			name:      "update success",
			getReturn: &dbmodel.UserTerm{Id: 2, StuId: stuId, TermTime: "old"},
			termList:  "new",
		},
		{
			name:      "update error",
			getReturn: &dbmodel.UserTerm{Id: 3, StuId: stuId, TermTime: "old"},
			termList:  "new",
			updateErr: assert.AnError,
			expectErr: "assert.AnError",
		},
		{
			name:      "no change",
			getReturn: &dbmodel.UserTerm{Id: 4, StuId: stuId, TermTime: "same"},
			termList:  "same",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range cases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			mockey.Mock((*dbcourse.DBCourse).GetUserTermByStuId).Return(tc.getReturn, tc.getErr).Build()
			mockey.Mock((*utils.Snowflake).NextVal).Return(tc.nextVal, tc.nextErr).Build()
			mockey.Mock((*dbcourse.DBCourse).CreateUserTerm).Return(nil, tc.createErr).Build()
			mockey.Mock((*dbcourse.DBCourse).UpdateUserTerm).Return(nil, tc.updateErr).Build()

			svc := NewCourseService(context.Background(), mockClientSet, nil)
			err := svc.putTermToDatabase(stuId, tc.termList)
			if tc.expectErr != "" {
				assert.ErrorContains(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
