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

func TestCourseService_GetTermsList(t *testing.T) {
	type testCase struct {
		name            string
		mockTermsReturn *jwch.Term
		mockTermsError  error
		expectResult    []string
		expectError     string
		cacheExist      bool
		cacheGetError   error
		mockSetCacheErr error
		expectSetCache  bool
		expectTaskAdd   bool
		expectGetTerms  bool
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
			mockTermsError:  nil,
			expectSetCache:  true,
			expectTaskAdd:   true,
			expectGetTerms:  true,
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
			expectError:   "redis get error",
		},
		{
			name:           "cache miss get terms error",
			mockTermsError: errors.New("remote failure"),
			expectError:    "Get terms fail",
			expectGetTerms: true,
		},
		{
			name:            "cache miss set cache error",
			mockTermsReturn: successTerm,
			mockSetCacheErr: fmt.Errorf("set cache boom"),
			expectResult:    successTerm.Terms,
			expectSetCache:  true,
			expectTaskAdd:   true,
			expectGetTerms:  true,
		},
	}
	mockLoginData := &model.LoginData{
		Id:      "052106112",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var (
				getTermsCalled int
				taskAdded      int
				setCacheCalled int
				wg             *sync.WaitGroup
			)
			if tc.expectSetCache {
				wg = &sync.WaitGroup{}
				wg.Add(1)
			}

			mockey.Mock((*jwch.Student).GetTerms).To(func(_ *jwch.Student) (*jwch.Term, error) {
				getTermsCalled++
				return tc.mockTermsReturn, tc.mockTermsError
			}).Build()

			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).To(func(_ *taskqueue.BaseTaskQueue, _ string, _ taskqueue.QueueTask) {
				taskAdded++
			}).Build()

			mockey.Mock((*coursecache.CacheCourse).SetTermsCache).To(
				func(ctx context.Context, key string, list []string) error {
					setCacheCalled++
					if wg != nil {
						wg.Done()
					}
					return tc.mockSetCacheErr
				},
			).Build()

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
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))

			result, err := courseService.GetTermsList(&model.LoginData{
				Id:      "123456789",
				Cookies: "magic cookies",
			})

			if wg != nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(100 * time.Millisecond):
					t.Fatalf("SetTermsCache was not called in time")
				}
			} else {
				time.Sleep(10 * time.Millisecond)
			}

			if tc.expectError != "" {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}

			if tc.expectGetTerms {
				assert.Equal(t, 1, getTermsCalled)
			} else {
				assert.Equal(t, 0, getTermsCalled)
			}
			if tc.expectTaskAdd {
				assert.Equal(t, 1, taskAdded)
			} else {
				assert.Equal(t, 0, taskAdded)
			}
			if tc.expectSetCache {
				assert.GreaterOrEqual(t, setCacheCalled, 1)
			} else {
				assert.Equal(t, 0, setCacheCalled)
			}
		})
	}
}

func TestCourseService_GetTermsListYjsy(t *testing.T) {
	type testCase struct {
		name            string
		mockTermsReturn *yjsy.Term
		mockTermsError  error
		expectResult    []string
		expectError     string
		cacheExist      bool
		cacheGetError   error
		mockSetCacheErr error
		expectSetCache  bool
		expectTaskAdd   bool
		expectGetTerms  bool
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
			mockTermsError:  nil,
			expectSetCache:  true,
			expectTaskAdd:   true,
			expectGetTerms:  true,
		},
		{
			name:         "cache exist success",
			cacheExist:   true,
			expectResult: successTerm.Terms,
		},
		{
			name:          "cache exist but get cache error",
			cacheExist:    true,
			cacheGetError: fmt.Errorf("redis get error"),
			expectError:   "redis get error",
		},
		{
			name:           "cache miss get terms error",
			mockTermsError: errors.New("remote failure"),
			expectError:    "Get terms fail",
			expectGetTerms: true,
		},
		{
			name:            "cache miss set cache error",
			mockTermsReturn: successTerm,
			mockSetCacheErr: fmt.Errorf("set cache boom"),
			expectResult:    successTerm.Terms,
			expectSetCache:  true,
			expectTaskAdd:   true,
			expectGetTerms:  true,
		},
	}
	mockLoginData := &model.LoginData{
		Id:      "052106112",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var (
				getTermsCalled int
				taskAdded      int
				setCacheCalled int
				wg             *sync.WaitGroup
			)
			if tc.expectSetCache {
				wg = &sync.WaitGroup{}
				wg.Add(1)
			}

			mockey.Mock((*yjsy.Student).GetTerms).To(func(_ *yjsy.Student) (*yjsy.Term, error) {
				getTermsCalled++
				return tc.mockTermsReturn, tc.mockTermsError
			}).Build()

			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).To(func(_ *taskqueue.BaseTaskQueue, _ string, _ taskqueue.QueueTask) {
				taskAdded++
			}).Build()

			mockey.Mock((*coursecache.CacheCourse).SetTermsCache).To(
				func(ctx context.Context, key string, list []string) error {
					setCacheCalled++
					if wg != nil {
						wg.Done()
					}
					return tc.mockSetCacheErr
				},
			).Build()

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
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))

			result, err := courseService.GetTermsListYjsy(&model.LoginData{
				Id:      "123456789",
				Cookies: "magic cookies",
			})

			if wg != nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(100 * time.Millisecond):
					t.Fatalf("SetTermsCache was not called in time")
				}
			} else {
				time.Sleep(10 * time.Millisecond)
			}

			if tc.expectError != "" {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}

			if tc.expectGetTerms {
				assert.Equal(t, 1, getTermsCalled)
			} else {
				assert.Equal(t, 0, getTermsCalled)
			}
			if tc.expectTaskAdd {
				assert.Equal(t, 1, taskAdded)
			} else {
				assert.Equal(t, 0, taskAdded)
			}
			if tc.expectSetCache {
				assert.GreaterOrEqual(t, setCacheCalled, 1)
			} else {
				assert.Equal(t, 0, setCacheCalled)
			}
		})
	}
}

func TestCourseService_putTermToDatabase(t *testing.T) {
	stuId := "stu-1"
	cases := []struct {
		name          string
		getReturn     *dbmodel.UserTerm
		getErr        error
		nextVal       int64
		nextErr       error
		createErr     error
		updateErr     error
		termList      string
		expectErr     string
		expectCreate  int
		expectUpdate  int
		expectNextVal int
	}{
		{
			name:      "db get error",
			getErr:    errors.New("db boom"),
			termList:  "202401",
			expectErr: "db boom",
		},
		{
			name:          "create success",
			getReturn:     nil,
			nextVal:       10,
			termList:      "202401",
			expectCreate:  1,
			expectNextVal: 1,
		},
		{
			name:          "create nextval error",
			getReturn:     nil,
			nextErr:       errors.New("sf fail"),
			termList:      "202401",
			expectErr:     "sf fail",
			expectNextVal: 1,
		},
		{
			name:          "create insert error",
			getReturn:     nil,
			nextVal:       12,
			createErr:     errors.New("insert fail"),
			termList:      "202401",
			expectErr:     "insert fail",
			expectCreate:  1,
			expectNextVal: 1,
		},
		{
			name:         "update success",
			getReturn:    &dbmodel.UserTerm{Id: 2, StuId: stuId, TermTime: "old"},
			termList:     "new",
			expectUpdate: 1,
		},
		{
			name:         "update error",
			getReturn:    &dbmodel.UserTerm{Id: 3, StuId: stuId, TermTime: "old"},
			termList:     "new",
			updateErr:    errors.New("update fail"),
			expectErr:    "update fail",
			expectUpdate: 1,
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
			var (
				getCalled    int
				nextCalled   int
				createCalled int
				updateCalled int
			)

			mockey.Mock((*dbcourse.DBCourse).GetUserTermByStuId).To(func(_ *dbcourse.DBCourse, _ context.Context, _ string) (*dbmodel.UserTerm, error) {
				getCalled++
				return tc.getReturn, tc.getErr
			}).Build()
			mockey.Mock((*utils.Snowflake).NextVal).To(func(_ *utils.Snowflake) (int64, error) {
				nextCalled++
				return tc.nextVal, tc.nextErr
			}).Build()
			mockey.Mock((*dbcourse.DBCourse).CreateUserTerm).To(func(_ *dbcourse.DBCourse, _ context.Context, m *dbmodel.UserTerm) (*dbmodel.UserTerm, error) {
				createCalled++
				return m, tc.createErr
			}).Build()
			mockey.Mock((*dbcourse.DBCourse).UpdateUserTerm).To(func(_ *dbcourse.DBCourse, _ context.Context, m *dbmodel.UserTerm) (*dbmodel.UserTerm, error) {
				updateCalled++
				return m, tc.updateErr
			}).Build()

			svc := NewCourseService(context.Background(), &base.ClientSet{
				DBClient: &db.Database{Course: new(dbcourse.DBCourse)},
				SFClient: new(utils.Snowflake),
			}, nil)

			err := svc.putTermToDatabase(stuId, tc.termList)

			if tc.expectErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, 1, getCalled)
			assert.Equal(t, tc.expectNextVal, nextCalled)
			assert.Equal(t, tc.expectCreate, createCalled)
			assert.Equal(t, tc.expectUpdate, updateCalled)
		})
	}
}
