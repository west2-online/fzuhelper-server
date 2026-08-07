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
	"bytes"
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestSSOLogin(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		authorizeUrl   string
		rpcError       error
		expectStatus   int
		expectLocation string
		expectContains string
		expectReturnTo string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/admin/auth/login?returnTo=http%3A%2F%2Flocalhost%3A5173%2Fauth%2Fcallback",
			authorizeUrl:   "https://accounts.feishu.cn/open-apis/authen/v1/authorize?state=test-state",
			expectStatus:   consts.StatusFound,
			expectLocation: "https://accounts.feishu.cn/open-apis/authen/v1/authorize?state=test-state",
			expectReturnTo: "http://localhost:5173/auth/callback",
		},
		{
			name:           "bind error",
			url:            "/api/v1/admin/auth/login",
			expectStatus:   consts.StatusOK,
			expectContains: `"code":"20001"`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/admin/auth/login?returnTo=http%3A%2F%2Flocalhost%3A5173%2Fauth%2Fcallback",
			rpcError:       errno.RedisError,
			expectStatus:   consts.StatusOK,
			expectContains: `"code":"50003"`,
			expectReturnTo: "http://localhost:5173/auth/callback",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/admin/auth/login", SSOLogin)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var actualReturnTo string
			mockey.Mock(rpc.SSOLoginRPC).To(func(_ context.Context, req *admin.LoginRequest) (string, error) {
				actualReturnTo = req.ReturnTo
				return tc.authorizeUrl, tc.rpcError
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)

			assert.Equal(t, tc.expectStatus, res.Result().StatusCode())
			assert.Equal(t, tc.expectLocation, string(res.Result().Header.Peek("Location")))
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			assert.Equal(t, tc.expectReturnTo, actualReturnTo)
		})
	}
}

func TestCallback(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		redirectUrl    string
		rpcError       error
		expectStatus   int
		expectLocation string
		expectContains string
		expectState    string
		expectCode     string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/admin/auth/callback?state=test-state&code=test-code",
			redirectUrl:    "http://localhost:5173/auth/callback?ticket=test-ticket",
			expectStatus:   consts.StatusFound,
			expectLocation: "http://localhost:5173/auth/callback?ticket=test-ticket",
			expectState:    "test-state",
			expectCode:     "test-code",
		},
		{
			name:           "bind error",
			url:            "/api/v1/admin/auth/callback?state=test-state",
			expectStatus:   consts.StatusOK,
			expectContains: `"code":"20001"`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/admin/auth/callback?state=expired&code=test-code",
			rpcError:       errno.AuthError,
			expectStatus:   consts.StatusOK,
			expectContains: `"code":"30001"`,
			expectState:    "expired",
			expectCode:     "test-code",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/admin/auth/callback", Callback)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var actualState, actualCode string
			mockey.Mock(rpc.AdminCallbackRPC).To(func(_ context.Context, req *admin.CallbackRequest) (string, error) {
				actualState = req.State
				actualCode = req.Code
				return tc.redirectUrl, tc.rpcError
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)

			assert.Equal(t, tc.expectStatus, res.Result().StatusCode())
			assert.Equal(t, tc.expectLocation, string(res.Result().Header.Peek("Location")))
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			assert.Equal(t, tc.expectState, actualState)
			assert.Equal(t, tc.expectCode, actualCode)
		})
	}
}

func TestExchange(t *testing.T) {
	type testCase struct {
		name           string
		body           string
		accessToken    string
		rpcError       error
		expectContains string
		expectTicket   string
	}

	testCases := []testCase{
		{
			name:           "success",
			body:           `{"ticket":"login-ticket"}`,
			accessToken:    "admin-access-token",
			expectContains: `"code":"10000","message":"Success","data":{"accessToken":"admin-access-token"}`,
			expectTicket:   "login-ticket",
		},
		{
			name:           "bind error",
			body:           `{"ticket":`,
			expectContains: `"code":"20001"`,
		},
		{
			name:           "rpc error",
			body:           `{"ticket":"expired-ticket"}`,
			rpcError:       errno.AuthError,
			expectContains: `"code":"30001"`,
			expectTicket:   "expired-ticket",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/admin/auth/exchange", Exchange)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var actualTicket string
			mockey.Mock(rpc.AdminExchangeTicketRPC).To(func(_ context.Context, req *admin.ExchangeTicketRequest) (string, error) {
				actualTicket = req.Ticket
				return tc.accessToken, tc.rpcError
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/admin/auth/exchange",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"})

			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			assert.Equal(t, tc.expectTicket, actualTicket)
		})
	}
}
