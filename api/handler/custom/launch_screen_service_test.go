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

package custom

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestMobileGetImage(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       []*model.Picture
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name: "success",
			url:  "/launch_screen/api/screen?type=1&student_id=102300001&college=computer&device=ios",
			mockResp: []*model.Picture{
				{
					Id:   1,
					Url:  "https://example.com/image.png",
					Href: "https://example.com",
				},
			},
			expectContains: `"code":200,"message":"ok","data":`,
		},
		{
			name:           "rpc error",
			url:            "/launch_screen/api/screen?type=1&student_id=102300001&college=computer&device=ios",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error - missing type",
			url:            "/launch_screen/api/screen?student_id=102300001&college=computer&device=ios",
			expectContains: `"code":"20001","message":"参数错误,`,
		},
		{
			name:           "bind error - missing student_id",
			url:            "/launch_screen/api/screen?type=1&college=computer&device=ios",
			expectContains: `"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/launch_screen/api/screen", MobileGetImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.MobileGetImageRPC).To(func(ctx context.Context, req *launch_screen.MobileGetImageRequest) ([]*model.Picture, *int64, error) {
				return tc.mockResp, nil, tc.mockRPCErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestAddImagePointTime(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/launch_screen/api/image/point?picture_id=1",
			expectContains: `"code":200,"message":"ok"`,
		},
		{
			name:           "rpc error",
			url:            "/launch_screen/api/image/point?picture_id=1",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error - missing picture_id",
			url:            "/launch_screen/api/image/point",
			expectContains: `"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/launch_screen/api/image/point", AddImagePointTime)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.AddImagePointTimeRPC).To(func(ctx context.Context, req *launch_screen.AddImagePointTimeRequest) (*model.Picture, error) {
				return nil, tc.mockRPCErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
