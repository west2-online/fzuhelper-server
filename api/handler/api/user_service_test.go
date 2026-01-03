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
	"mime/multipart"
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
	router.GET("/api/v1/user/invite", func(c context.Context, h *app.RequestContext) {
		GetInvitationCode(c, h)
	})

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
	router.GET("/api/v1/user/friend/info", func(c context.Context, h *app.RequestContext) {
		GetFriendList(c, h)
	})

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetFriendListRPC).To(func(ctx context.Context, req *user.GetFriendListRequest) ([]*model.UserFriendInfo, error) {
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
	router.POST("/api/v1/user/friend/delete", func(c context.Context, h *app.RequestContext) {
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
	router.POST("/api/v1/user/friend/invite/cancel", func(c context.Context, h *app.RequestContext) {
		CancelInvite(c, h)
	})

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

func TestValidateCode(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		method         string
		imageString    string
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/user/validateCode",
			method:         "POST",
			imageString:    "data:image/png;base64,Qk2mCAAAAAAAADYAAAAoAAAASAAAAAoAAAABABgAAAAAAHAIAAASCwAAEgsAAAAAAAAAAAAA+vr/+vr/+vr/lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/AACWAACW+vr/AACW+vr/+vr/+vr/lgD6+vr/lgD6lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACWAACWAACW+vr/+vr/+vr/+vr/+vr/",
			expectContains: `{"code":"200","data":"104","message":"success"}`,
		},
		{
			name:           "invalid_base64",
			url:            "/api/v1/jwch/user/validateCode",
			method:         "POST",
			imageString:    "not-base64",
			expectContains: `"50001"`,
		},
		{
			name:           "malformed_image",
			url:            "/api/v1/jwch/user/validateCode",
			method:         "POST",
			imageString:    "data:image/png;base64,YWJjZA==",
			expectContains: `"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/jwch/user/validateCode", func(c context.Context, h *app.RequestContext) {
		ValidateCode(c, h)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			mw := multipart.NewWriter(&buf)
			_ = mw.WriteField("image", tc.imageString)
			_ = mw.Close()
			res := ut.PerformRequest(router, tc.method, tc.url,
				&ut.Body{Body: bytes.NewBuffer(buf.Bytes()), Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: mw.FormDataContentType()},
			)
			if tc.expectContains != "" {
				assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			}
		})
	}
}

func TestValidateCodeForAndroid(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		method         string
		imageString    string
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/login/validateCode",
			method:         "POST",
			imageString:    "data:image/png;base64,Qk2mCAAAAAAAADYAAAAoAAAASAAAAAoAAAABABgAAAAAAHAIAAASCwAAEgsAAAAAAAAAAAAA+vr/+vr/+vr/lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/AACWAACW+vr/AACW+vr/+vr/+vr/lgD6+vr/lgD6lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACWAACWAACW+vr/+vr/+vr/+vr/+vr/",
			expectContains: `{"code":"200","message":"104"}`,
		},
		{
			name:           "invalid_base64",
			url:            "/api/login/validateCode",
			method:         "POST",
			imageString:    "not-base64",
			expectContains: `"50001"`,
		},
		{
			name:           "malformed_image",
			url:            "/api/login/validateCode",
			method:         "POST",
			imageString:    "data:image/png;base64,YWJjZA==",
			expectContains: `"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/login/validateCode", func(c context.Context, h *app.RequestContext) {
		ValidateCodeForAndroid(c, h)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			mw := multipart.NewWriter(&buf)
			_ = mw.WriteField("validateCode", tc.imageString)
			_ = mw.Close()
			res := ut.PerformRequest(router, tc.method, tc.url,
				&ut.Body{Body: bytes.NewBuffer(buf.Bytes()), Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: mw.FormDataContentType()},
			)
			if tc.expectContains != "" {
				assert.Contains(t, string(res.Result().Body()), tc.expectContains)
			}
		})
	}
}
