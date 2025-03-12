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
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	userDB "github.com/west2-online/fzuhelper-server/pkg/db/user"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestUserService_GetUserInfo(t *testing.T) {
	// 测试案例使用的结构体
	type testCase struct {
		name string
		// 原有字段
		expectedInfo      *dbmodel.Student
		expectedJwch      *jwch.StudentDetail
		expectedExist     bool
		mockError         error
		mockDBCreateError error
		mockJwchError     error
		expectingError    bool
		expectingErrorMsg string

		// 新增字段：用于控制缓存的场景
		cacheExist    bool             // 是否在 Redis 中存在这个 Key
		cacheGetError error            // 获取缓存时是否模拟报错
		cacheStudent  *dbmodel.Student // 如果缓存命中时，要返回的缓存结果
	}
	// 构造一个 dbmodel.Student 作为测试中的“数据库期望返回”
	info := &dbmodel.Student{
		StuId:    "102301000",
		Name:     "testName",
		Sex:      "sex",
		Birthday: "1970-01-01",
		College:  "计算机与大数据学院",
		Grade:    2023,
		Major:    "计算机类",
	}
	// jwch 返回的 StudentDetail
	stuDetail := &jwch.StudentDetail{
		Name:     info.Name,
		Sex:      info.Sex,
		Birthday: info.Birthday,
		College:  info.College,
		Grade:    strconv.FormatInt(info.Grade, 10),
		Major:    info.Major,
	}

	//// 用来表示从 Cache 拿到的数据，这里演示一下不同字段值
	cacheStu := &dbmodel.Student{
		StuId:    info.StuId,
		Name:     "cacheName",
		Sex:      "cacheSex",
		Birthday: "1970-02-02",
		College:  "Another College",
		Grade:    2021,
		Major:    "Cache Major",
	}

	testCases := []testCase{
		{
			name:           "db exist success (no cache)",
			expectedExist:  true,
			expectingError: false,
			expectedInfo:   info,
			mockError:      nil,
		},
		{
			name:           "db not exist success, jwch success, create db success",
			expectedExist:  false,
			expectingError: false,
			expectedInfo:   info,
			expectedJwch:   stuDetail,
			mockError:      nil,
		},
		{
			name:              "jwch error",
			expectedExist:     false,
			expectingError:    true,
			expectedInfo:      info,
			expectedJwch:      stuDetail,
			mockJwchError:     errno.InternalServiceError,
			expectingErrorMsg: errno.InternalServiceError.ErrorMsg,
		},
		{
			name:              "db create error",
			expectedExist:     false,
			expectingError:    true,
			expectedInfo:      info,
			expectedJwch:      stuDetail,
			mockDBCreateError: gorm.ErrInvalidData,
			expectingErrorMsg: "service.GetUserInfo:",
		},
		{
			name:              "db error",
			expectedExist:     false,
			expectingError:    true,
			expectedInfo:      info,
			expectedJwch:      stuDetail,
			mockError:         gorm.ErrInvalidData,
			expectingErrorMsg: "service.GetUserInfo:",
		},
		//// ------------------- 以下为缓存相关测试场景示例 -------------------
		{
			name:          "cache exist success",
			cacheExist:    true, // 缓存里已存在
			cacheGetError: nil,  // 获取缓存不报错
			cacheStudent:  cacheStu,
		},
	}

	defer mockey.UnPatchAll()
	mockey.Mock((*user.CacheUser).SetStuInfoCache).To(
		func(ctx context.Context, key string, stu *dbmodel.Student) error {
			return nil
		},
	).Build()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// 初始化
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			// Mock Cache 方法
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()

			mockey.Mock(time.Time.After).Return(true).Build()

			// 如果缓存存在，则 Mock GetStuInfoCache
			if tc.cacheExist {
				mockey.Mock((*user.CacheUser).GetStuInfoCache).To(
					func(ctx context.Context, key string) (*dbmodel.Student, error) {
						if tc.cacheGetError != nil {
							return nil, tc.cacheGetError
						}
						return tc.cacheStudent, nil
					},
				).Build()
			} else {
				// 如果缓存不存在，一般不会去调 GetStuInfoCache
				// 也可以不 Mock，或 Mock 一个默认返回
				mockey.Mock((*user.CacheUser).GetStuInfoCache).To(
					func(ctx context.Context, key string) (*dbmodel.Student, error) {
						return nil, fmt.Errorf("should not be called if cache doesn't exist")
					},
				).Build()
			}

			// Mock DB 方法
			mockey.Mock((*userDB.DBUser).GetStudentById).To(
				func(ctx context.Context, stuId string) (bool, *dbmodel.Student, error) {
					return tc.expectedExist, tc.expectedInfo, tc.mockError
				},
			).Build()
			mockey.Mock((*userDB.DBUser).CreateStudent).To(
				func(ctx context.Context, userModel *dbmodel.Student) error {
					return tc.mockDBCreateError
				},
			).Build()

			// Mock JWCH
			mockey.Mock((*jwch.Student).WithLoginData).To(
				func(identifier string, cookies []*http.Cookie) *jwch.Student {
					return jwch.NewStudent()
				},
			).Build()
			mockey.Mock((*jwch.Student).GetInfo).To(
				func() (*jwch.StudentDetail, error) {
					return tc.expectedJwch, tc.mockJwchError
				},
			).Build()

			// 开始测试
			stuInfo, err := userService.GetUserInfo(info.StuId)

			// 判断是否期望报错
			if tc.expectingError {
				assert.Nil(t, stuInfo, "stuInfo should be nil on error")
				assert.Error(t, err, "error should not be nil")
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
			} else {
				assert.NoError(t, err, "should be no error")
				assert.NotNil(t, stuInfo, "student info should not be nil on success")
				// 如果缓存存在且拿取缓存成功，期望直接返回缓存数据（看你实际业务需要）
				if tc.cacheExist && tc.cacheGetError == nil {
					assert.Equal(t, tc.cacheStudent.StuId, stuInfo.StuId)
					assert.Equal(t, tc.cacheStudent.Sex, stuInfo.Sex)
					assert.Equal(t, tc.cacheStudent.College, stuInfo.College)
				} else {
					// 否则说明是走 DB 或 JWCH 流程
					assert.Equal(t, tc.expectedInfo.StuId, stuInfo.StuId)
					assert.Equal(t, tc.expectedInfo.Sex, stuInfo.Sex)
				}
			}
		})
	}
}
