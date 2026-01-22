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

package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/mw"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	metainfocontext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func TestGetLoginData(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockID         string
		mockCookies    string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/login?id=user123&password=pass123",
			mockID:         "202400001",
			mockCookies:    "session_cookie_value",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "bind error - missing params",
			url:            "/api/v1/user/login",
			mockID:         "",
			mockCookies:    "",
			mockRPCError:   nil,
			expectContains: `{"code":"20001","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/login?id=user123&password=wrongpass",
			mockID:         "",
			mockCookies:    "",
			mockRPCError:   errors.New("invalid credentials"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/login", GetLoginData)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetLoginDataRPC).To(func(ctx context.Context, req *user.GetLoginDataRequest) (string, string, error) {
				if tc.mockRPCError != nil {
					return "", "", tc.mockRPCError
				}
				return tc.mockID, tc.mockCookies, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type testCase struct {
		name           string
		method         string
		url            string
		authHeader     string
		mockTokenType  int64
		mockCheckErr   error
		mockCreateErr  error
		expectContains string
		expectStatus   int
	}

	testCases := []testCase{
		{
			name:           "success",
			method:         consts.MethodPost,
			url:            "/api/v1/login/refreshToken",
			authHeader:     "valid_refresh_token",
			mockTokenType:  1, // TypeRefreshToken = 1
			mockCheckErr:   nil,
			mockCreateErr:  nil,
			expectContains: `"code":"10000"`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "auth missing - no token",
			method:         consts.MethodPost,
			url:            "/api/v1/login/refreshToken",
			authHeader:     "",
			expectContains: `{"code":"30002","message":`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "check token failed",
			method:         consts.MethodPost,
			url:            "/api/v1/login/refreshToken",
			authHeader:     "invalid_token",
			mockCheckErr:   errno.AuthError,
			expectContains: `"code":"30001"`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "token type is access token, not refresh token",
			method:         consts.MethodPost,
			url:            "/api/v1/login/refreshToken",
			authHeader:     "valid_access_token",
			mockTokenType:  0, // TypeAccessToken = 0
			mockCheckErr:   nil,
			expectContains: `"code":"30002"`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "create token failed",
			method:         consts.MethodPost,
			url:            "/api/v1/login/refreshToken",
			authHeader:     "valid_refresh_token",
			mockTokenType:  1, // TypeRefreshToken = 1
			mockCheckErr:   nil,
			mockCreateErr:  errors.New("token generation failed"),
			expectContains: `"code":"50001"`,
			expectStatus:   consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/login/refreshToken", RefreshToken)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(mw.CheckToken).To(func(token string) (int64, string, error) {
				if tc.mockCheckErr != nil {
					return 0, "", tc.mockCheckErr
				}
				return tc.mockTokenType, "test_stu_id", nil
			}).Build()

			mockey.Mock(mw.CreateAllToken).To(func() (string, string, error) {
				if tc.mockCreateErr != nil {
					return "", "", tc.mockCreateErr
				}
				return "access_token", "refresh_token", nil
			}).Build()

			var headers []ut.Header
			if tc.authHeader != "" {
				headers = []ut.Header{
					{Key: "Authorization", Value: tc.authHeader},
				}
			}
			res := ut.PerformRequest(router, tc.method, tc.url, nil, headers...)
			assert.Equal(t, tc.expectStatus, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetToken(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		isGraduate     bool
		mockCheckError error
		mockTokenError error
		expectContains string
		expectStatus   int
	}

	testCases := []testCase{
		{
			name:           "success - undergraduate (jwch)",
			url:            "/api/v1/login/access-token",
			isGraduate:     false,
			mockCheckError: nil,
			mockTokenError: nil,
			expectContains: `"code":"10000"`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "success - graduate (yjsy)",
			url:            "/api/v1/login/access-token",
			isGraduate:     true,
			mockCheckError: nil,
			mockTokenError: nil,
			expectContains: `"code":"10000"`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "jwch check session failed",
			url:            "/api/v1/login/access-token",
			isGraduate:     false,
			mockCheckError: errors.New("invalid session"),
			mockTokenError: nil,
			expectContains: `"code":"30001"`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "yjsy check session failed",
			url:            "/api/v1/login/access-token",
			isGraduate:     true,
			mockCheckError: errors.New("invalid session"),
			mockTokenError: nil,
			expectContains: `"code":"30001"`,
			expectStatus:   consts.StatusOK,
		},
		{
			name:           "create token failed",
			url:            "/api/v1/login/access-token",
			isGraduate:     false,
			mockCheckError: nil,
			mockTokenError: errors.New("token creation failed"),
			expectContains: `"code":"50001"`,
			expectStatus:   consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/login/access-token", GetToken)
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock utils.IsGraduate
			mockey.Mock(utils.IsGraduate).To(func(identifier string) bool {
				return tc.isGraduate
			}).Build()

			// Mock jwch.NewStudent().CheckSession
			mockey.Mock((*jwch.Student).CheckSession).To(func(s *jwch.Student) error {
				return tc.mockCheckError
			}).Build()

			// Mock yjsy.NewStudent().CheckSession
			mockey.Mock((*yjsy.Student).CheckSession).To(func(s *yjsy.Student) error {
				return tc.mockCheckError
			}).Build()

			// Mock mw.CreateAllToken
			mockey.Mock(mw.CreateAllToken).To(func() (string, string, error) {
				if tc.mockTokenError != nil {
					return "", "", tc.mockTokenError
				}
				return "access_token_test", "refresh_token_test", nil
			}).Build()

			// Mock utils.ParseCookies to return a non-nil slice
			mockey.Mock(utils.ParseCookies).To(func(cookiesStr string) []*http.Cookie {
				return []*http.Cookie{
					{Name: "ASP.NET_SessionId", Value: "test_session"},
				}
			}).Build()

			// Mock metainfocontext.ExtractIDFromIdentifier
			mockey.Mock(metainfocontext.ExtractIDFromIdentifier).To(func(identifier string) string {
				return "052106112"
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil,
				ut.Header{Key: "id", Value: "202412615623052106112"},
				ut.Header{Key: "cookies", Value: "test_cookies"})
			assert.Equal(t, tc.expectStatus, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestTestAuth(t *testing.T) {
	type testCase struct {
		name         string
		url          string
		expectStatus int
		expectCode   string
	}

	testCases := []testCase{
		{
			name:         "success",
			url:          "/api/v1/login/ping",
			expectStatus: consts.StatusOK,
			expectCode:   `"code":"10000"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/login/ping", TestAuth)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, tc.expectStatus, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectCode)
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockInfo       *model.UserInfo
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/user/info",
			mockInfo:       &model.UserInfo{Name: "Test User"},
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/user/info",
			mockInfo:       nil,
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/user/info", GetUserInfo)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetUserInfoRPC).To(func(ctx context.Context, req *user.GetUserInfoRequest) (*model.UserInfo, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return tc.mockInfo, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetGetLoginDataForYJSY(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockID         string
		mockCookies    string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/internal/yjsy/user/login?id=user123&password=pass123",
			mockID:         "202400001",
			mockCookies:    "session_cookie_value",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "bind error - missing params",
			url:            "/api/v1/internal/yjsy/user/login",
			mockID:         "",
			mockCookies:    "",
			mockRPCError:   nil,
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/internal/yjsy/user/login?id=user123&password=wrongpass",
			mockID:         "",
			mockCookies:    "",
			mockRPCError:   errors.New("authentication failed"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/internal/yjsy/user/login", GetGetLoginDataForYJSY)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetLoginDataForYJSYRPC).To(func(ctx context.Context, req *user.GetLoginDataForYJSYRequest) (string, string, error) {
				if tc.mockRPCError != nil {
					return "", "", tc.mockRPCError
				}
				return tc.mockID, tc.mockCookies, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetInvitationCode(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockCode       string
		mockCreatedAt  int64
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/invite",
			mockCode:       "ABCD1234",
			mockCreatedAt:  1690000000,
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "with refresh",
			url:            "/api/v1/user/invite?is_refresh=true",
			mockCode:       "EFGH5678",
			mockCreatedAt:  1690000000,
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/invite",
			mockCode:       "",
			mockCreatedAt:  -1,
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/invite", GetInvitationCode)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetInvitationCodeRPC).To(func(ctx context.Context, req *user.GetInvitationCodeRequest) (string, int64, error) {
				if tc.mockRPCError != nil {
					return "", -1, tc.mockRPCError
				}
				return tc.mockCode, tc.mockCreatedAt, nil
			}).Build()

			res := ut.PerformRequest(router, "GET", tc.url, nil)
			if tc.expectContains != "" {
				assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			}
		})
	}
}

func TestBindInvitation(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/friend/bind?invitation_code=ABCD1234",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/user/friend/bind",
			mockRPCError:   nil,
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/friend/bind?invitation_code=ABCD1234",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/friend/bind", BindInvitation)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.BindInvitationRPC).To(func(ctx context.Context, req *user.BindInvitationRequest) error {
				return tc.mockRPCError
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			if tc.expectContains != "" {
				assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			}
		})
	}
}

func TestGetFriendList(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockInfo       []*model.UserFriendInfo
		mockRPCError   error
		expectContains string
	}

	okInfo := []*model.UserFriendInfo{
		{
			StuId: "102300217",
		},
		{
			StuId: "102300218",
		},
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/friend/info",
			mockInfo:       okInfo,
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/friend/info",
			mockInfo:       nil,
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/friend/info", GetFriendList)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetFriendListRPC).To(func(ctx context.Context, req *user.GetFriendListRequest) ([]*model.UserFriendInfo, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return tc.mockInfo, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			if tc.expectContains != "" {
				assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			}
		})
	}
}

func TestDeleteFriend(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		method         string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/friend/delete?student_id=1",
			method:         "POST",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/user/friend/delete",
			method:         "POST",
			mockRPCError:   nil,
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "DELETE",
			url:            "/api/v1/user/friend/delete?student_id=1",
			method:         "POST",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/user/friend/delete", DeleteFriend)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.DeleteFriendRPC).To(func(ctx context.Context, req *user.DeleteFriendRequest) error {
				return tc.mockRPCError
			}).Build()

			res := ut.PerformRequest(router, tc.method, tc.url, nil)
			if tc.expectContains != "" {
				assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			}
		})
	}
}

func TestCancelInvite(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		method         string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/friend/invite/cancel",
			method:         "POST",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/friend/invite/cancel",
			method:         "POST",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/user/friend/invite/cancel", CancelInvite)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.CancelInviteRPC).To(func(ctx context.Context, req *user.CancelInviteRequest) error {
				return tc.mockRPCError
			}).Build()

			res := ut.PerformRequest(router, tc.method, tc.url, nil)
			if tc.expectContains != "" {
				assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			}
		})
	}
}
