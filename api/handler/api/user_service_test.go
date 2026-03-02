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
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "bind error - missing params",
			url:            "/api/v1/user/login",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/login?id=user123&password=wrongpass",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/login", GetLoginData)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetLoginDataRPC).To(func(ctx context.Context, req *user.GetLoginDataRequest) (string, string, error) {
				return tc.mockID, tc.mockCookies, tc.mockRPCError
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
		url            string
		authHeader     string
		mockTokenType  int64
		mockCheckErr   error
		mockCreateErr  error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/login/refreshToken",
			authHeader:     "valid_refresh_token",
			mockTokenType:  1, // TypeRefreshToken = 1
			expectContains: `"code":"10000","message":"ok"`,
		},
		{
			name:           "auth missing - no token",
			url:            "/api/v1/login/refreshToken",
			expectContains: `{"code":"30002","message":"缺失合法鉴权数据"`,
		},
		{
			name:           "check token failed",
			url:            "/api/v1/login/refreshToken",
			authHeader:     "invalid_token",
			mockCheckErr:   errno.AuthError,
			expectContains: `"code":"30001","message":"鉴权失败"`,
		},
		{
			name:           "token type is access token, not refresh token",
			url:            "/api/v1/login/refreshToken",
			authHeader:     "valid_access_token",
			mockTokenType:  0, // TypeAccessToken = 0
			expectContains: `"code":"30002","message":"token type is access token, need refresh token"`,
		},
		{
			name:           "create token failed",
			url:            "/api/v1/login/refreshToken",
			authHeader:     "valid_refresh_token",
			mockTokenType:  1, // TypeRefreshToken = 1
			mockCreateErr:  errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/login/refreshToken", RefreshToken)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(mw.CheckToken).To(func(token string) (int64, string, error) {
				return tc.mockTokenType, "", tc.mockCheckErr
			}).Build()

			mockey.Mock(mw.CreateAllToken).To(func() (string, string, error) {
				return "", "", tc.mockCreateErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil,
				ut.Header{Key: "Authorization", Value: tc.authHeader})
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
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
	}

	testCases := []testCase{
		{
			name:           "success - undergraduate (jwch)",
			url:            "/api/v1/login/access-token",
			expectContains: `"code":"10000","message":"ok"`,
		},
		{
			name:           "success - graduate (yjsy)",
			url:            "/api/v1/login/access-token",
			isGraduate:     true,
			expectContains: `"code":"10000","message":"ok"`,
		},
		{
			name:           "jwch check session failed",
			url:            "/api/v1/login/access-token",
			mockCheckError: errno.AuthError,
			expectContains: `"code":"30001","message":"(jwch) check id and session failed`,
		},
		{
			name:           "yjsy check session failed",
			url:            "/api/v1/login/access-token",
			isGraduate:     true,
			mockCheckError: errno.AuthError,
			expectContains: `"code":"30001","message":"(yjsy) check id and session failed`,
		},
		{
			name:           "create token failed",
			url:            "/api/v1/login/access-token",
			mockTokenError: errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/login/access-token", GetToken)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(utils.IsGraduate).To(func(identifier string) bool {
				return tc.isGraduate
			}).Build()

			mockey.Mock((*jwch.Student).CheckSession).To(func(s *jwch.Student) error {
				return tc.mockCheckError
			}).Build()

			mockey.Mock((*yjsy.Student).CheckSession).To(func(s *yjsy.Student) error {
				return tc.mockCheckError
			}).Build()

			mockey.Mock(mw.CreateAllToken).To(func() (string, string, error) {
				return "", "", tc.mockTokenError
			}).Build()

			mockey.Mock(utils.ParseCookies).To(func(cookiesStr string) []*http.Cookie {
				return []*http.Cookie{
					{Name: "ASP.NET_SessionId", Value: "test_session"},
				}
			}).Build()

			mockey.Mock(metainfocontext.ExtractIDFromIdentifier).To(func(identifier string) string {
				return "052106112"
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil,
				ut.Header{Key: "id", Value: "202412615623052106112"},
				ut.Header{Key: "cookies", Value: "test_cookies"})
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestTestAuth(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/login/ping",
			expectContains: `"code":"10000","message":"ok"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/login/ping", TestAuth)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
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
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/user/info",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/user/info", GetUserInfo)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetUserInfoRPC).To(func(ctx context.Context, req *user.GetUserInfoRequest) (*model.UserInfo, error) {
				return tc.mockInfo, tc.mockRPCError
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
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "bind error - missing params",
			url:            "/api/v1/internal/yjsy/user/login",
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/internal/yjsy/user/login?id=user123&password=wrongpass",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/internal/yjsy/user/login", GetGetLoginDataForYJSY)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetLoginDataForYJSYRPC).To(func(ctx context.Context, req *user.GetLoginDataForYJSYRequest) (string, string, error) {
				return tc.mockID, tc.mockCookies, tc.mockRPCError
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
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "with refresh",
			url:            "/api/v1/user/invite?is_refresh=true",
			mockCode:       "EFGH5678",
			mockCreatedAt:  1690000000,
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/invite",
			mockCreatedAt:  -1,
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/invite", GetInvitationCode)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetInvitationCodeRPC).To(func(ctx context.Context, req *user.GetInvitationCodeRequest) (string, int64, error) {
				return tc.mockCode, tc.mockCreatedAt, tc.mockRPCError
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
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
			expectContains: `{"code":"10000","message":"ok"`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/user/friend/bind",
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/friend/bind?invitation_code=ABCD1234",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
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
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
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
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/friend/info",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/friend/info", GetFriendList)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetFriendListRPC).To(func(ctx context.Context, req *user.GetFriendListRequest) ([]*model.UserFriendInfo, error) {
				return tc.mockInfo, tc.mockRPCError
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestDeleteFriend(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/friend/delete?student_id=1",
			expectContains: `{"code":"10000","message":"ok"`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/user/friend/delete",
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "DELETE",
			url:            "/api/v1/user/friend/delete?student_id=1",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
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

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestCancelInvite(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/friend/invite/cancel",
			expectContains: `{"code":"10000","message":"ok"`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/friend/invite/cancel",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
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

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
