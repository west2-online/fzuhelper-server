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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestGetInvitationCode(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockCode       string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/invite",
			mockCode:       "ABCD1234",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "with refresh",
			url:            "/api/v1/user/invite?is_refresh=true",
			mockCode:       "EFGH5678",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/invite",
			mockCode:       "",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/user/invite", func(c context.Context, h *app.RequestContext) {
		GetInvitationCode(c, h)
	})

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetInvitationCodeRPC).To(func(ctx context.Context, req *user.GetInvitationCodeRequest) (string, error) {
				if tc.mockRPCError != nil {
					return "", tc.mockRPCError
				}
				return tc.mockCode, nil
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
	router.GET("/api/v1/user/friend/bind", func(c context.Context, h *app.RequestContext) {
		BindInvitation(c, h)
	})

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.BindInvitationRPC).To(func(ctx context.Context, req *user.BindInvitationRequest) error {
				return tc.mockRPCError
			}).Build()

			res := ut.PerformRequest(router, "GET", tc.url, nil)
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
		mockInfo       []*model.UserInfo
		mockRPCError   error
		expectContains string
	}

	okInfo := []*model.UserInfo{
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
	router.GET("/api/v1/user/friend/info", func(c context.Context, h *app.RequestContext) {
		GetFriendList(c, h)
	})

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetFriendListRPC).To(func(ctx context.Context, req *user.GetFriendListRequest) ([]*model.UserInfo, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return tc.mockInfo, nil
			}).Build()

			res := ut.PerformRequest(router, "GET", tc.url, nil)
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
			url:            "/api/v1/user/friend/delete?id=1",
			method:         "DELETE",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/user/friend/delete",
			method:         "DELETE",
			mockRPCError:   nil,
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/user/friend/delete?id=1",
			method:         "DELETE",
			mockRPCError:   errno.InternalServiceError,
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.DELETE("/api/v1/user/friend/delete", func(c context.Context, h *app.RequestContext) {
		DeleteFriend(c, h)
	})

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
