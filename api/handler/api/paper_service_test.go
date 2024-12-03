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
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
)

func TestGetDownloadUrl(t *testing.T) {
	type TestCase struct {
		ExpectedError  bool
		ExpectedResult interface{}
		FilePath       string
	}

	testCases := []TestCase{
		{
			ExpectedError:  false,
			ExpectedResult: `{"code":"10000","message":"Success","data":{"url":"file://url"}}`,
			FilePath:       "url",
		},
		{
			ExpectedError:  true,
			ExpectedResult: `{"code":"50001","message":"GetDownloadUrlRPC: RPC called failed: wrong filepath"}`,
			FilePath:       "",
		},
	}

	router := route.NewEngine(&config.Options{})

	router.GET("/api/v1/paper/download", func(ctx context.Context, c *app.RequestContext) {
		filepath := c.DefaultQuery("filepath", "/")
		mockey.Mock(rpc.GetDownloadUrlRPC).To(func(ctx context.Context, req *paper.GetDownloadUrlRequest) (url string, err error) {
			// 因为handler不能重复注册，无法添加一个TestError(bool)来作为错误的判断，只能先用filepath为空暂代了
			if filepath == "" {
				return "", errors.New("GetDownloadUrlRPC: RPC called failed: wrong filepath")
			}
			return "file://" + req.Filepath, nil
		}).Build()
		defer mockey.UnPatchAll()

		url, err := rpc.GetDownloadUrlRPC(ctx, &paper.GetDownloadUrlRequest{
			Filepath: filepath,
		})
		if err != nil {
			pack.RespError(c, err)
			return
		}

		resp := new(api.GetDownloadUrlResponse)
		resp.URL = url

		pack.RespData(c, resp)
	})
	for _, tc := range testCases {
		url := "/api/v1/paper/download" + "?filepath=" + tc.FilePath

		resp := ut.PerformRequest(router, consts.MethodGet, url, nil)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, tc.ExpectedResult, string(resp.Result().Body()))
	}
}
