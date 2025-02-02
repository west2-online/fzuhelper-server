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
	"net/http"
	"strconv"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	meta "github.com/west2-online/fzuhelper-server/pkg/base/context"
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
	type testCase struct {
		name              string
		expectedInfo      *dbmodel.Student
		expectedJwch      *jwch.StudentDetail
		expectedExist     bool
		mockError         error
		mockDBCreateError error
		mockJwchError     error
		expectingError    bool
		expectingErrorMsg string
	}
	info := &dbmodel.Student{
		StuId:    "102301000",
		Sex:      "sex",
		Birthday: "1970-01-01",
		College:  "计算机与大数据学院",
		Grade:    2023,
		Major:    "计算机类",
	}
	stuDetail := &jwch.StudentDetail{
		Sex:      info.Sex,
		Birthday: info.Birthday,
		College:  info.College,
		Grade:    strconv.FormatInt(info.Grade, 10),
		Major:    info.Major,
	}
	testCases := []testCase{
		{
			name:           "exist success",
			expectedExist:  true,
			expectingError: false,
			expectedInfo:   info,
			mockError:      nil,
		},
		{
			name:           "not exist success",
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
	}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			mockey.Mock((*cache.Cache).IsKeyExist).Return(false).Build()
			mockey.Mock((*user.CacheUser).SetStuInfoCache).Return(nil).Build()
			mockey.Mock((*user.CacheUser).GetStuInfoCache).Return(nil).Build()
			mockey.Mock((*userDB.DBUser).GetStudentById).To(func(ctx context.Context, stuId string) (bool, *dbmodel.Student, error) {
				return tc.expectedExist, tc.expectedInfo, tc.mockError
			}).Build()
			mockey.Mock((*jwch.Student).WithLoginData).To(func(identifier string, cookies []*http.Cookie) *jwch.Student {
				return jwch.NewStudent()
			}).Build()
			mockey.Mock(meta.GetLoginData).To(func(ctx context.Context) (*model.LoginData, error) {
				return &model.LoginData{
					Id:      "1111111111111111111111111111111111",
					Cookies: "",
				}, nil
			}).Build()
			mockey.Mock((*jwch.Student).GetInfo).To(func() (resp *jwch.StudentDetail, err error) {
				return tc.expectedJwch, tc.mockJwchError
			}).Build()
			mockey.Mock((*userDB.DBUser).CreateStudent).To(func(ctx context.Context, userModel *dbmodel.Student) error {
				return tc.mockDBCreateError
			}).Build()

			stuInfo, err := userService.GetUserInfo(info.StuId)
			if tc.expectingError {
				assert.Nil(t, stuInfo)
				assert.Contains(t, err.Error(), tc.expectingErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInfo, stuInfo)
			}
		})
	}
}
