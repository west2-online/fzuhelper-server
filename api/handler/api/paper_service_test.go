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
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
)

func TestGetDownloadUrl(t *testing.T) {
	type TestCase struct {
		Name           string
		ExpectedError  bool
		ExpectedResult interface{}
		Url            string
	}

	testCases := []TestCase{
		{
			Name:           "GetUrlSuccessfully",
			ExpectedError:  false,
			ExpectedResult: `{"code":"10000","message":"Success","data":{"url":"file://url"}}`,
			Url:            "/api/v1/paper/download?filepath=url",
		},
		{
			Name:           "GetUrlFailed",
			ExpectedError:  true,
			ExpectedResult: `{"code":"50001","message":"GetDownloadUrlRPC: RPC called failed: wrong filepath"}`,
			Url:            "/api/v1/paper/download?filepath=",
		},
		{
			Name:          "BindAndValidateError",
			ExpectedError: false,
			ExpectedResult: `{"code":"20001","message":"参数错误, 'filepath' field is a 'required' parameter` +
				`, but the request body does not have this parameter 'filepath'"}`,
			Url: "/api/v1/paper/download",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/paper/download", GetDownloadUrl)

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(rpc.GetDownloadUrlRPC).To(func(ctx context.Context, req *paper.GetDownloadUrlRequest) (url string, err error) {
				if tc.ExpectedError {
					return "", errors.New("GetDownloadUrlRPC: RPC called failed: wrong filepath")
				}
				return "file://" + req.Filepath, nil
			}).Build()
			defer mockey.UnPatchAll()

			resp := ut.PerformRequest(router, consts.MethodGet, tc.Url, nil)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, tc.ExpectedResult, string(resp.Result().Body()))
		})
	}
}

func TestListDirFiles(t *testing.T) {
	type TestCase struct {
		Name              string
		ExpectedError     bool
		ExpectedResult    interface{}
		ExpectUpYunResult *model.UpYunFileDir
		Path              string
	}
	basePath := "/C语言"
	expectedUpYunResult := &model.UpYunFileDir{
		BasePath: &basePath,
		Files: []string{
			"10份练习.zip",
			"200912填空题(含答案).doc",
			"200912改错题(含答案).doc",
			"200912程序题(含答案).doc",
			"201012真题C语言(含答案).doc",
			"2010年4月试题.doc",
			"2010年4月试题答案.doc",
			"2011-06改错填空编程(含答案).doc",
			"2011-06真题C语言(含答案).doc",
			"2012年福建省C语言二级考试大题及答案.doc",
			"2015C语言课件.zip",
			"2015上机题.zip",
			"2015选择题.zip",
			"C语言-张莹.zip",
			"C语言模拟卷含答案 (1).doc",
			"C语言模拟卷含答案 (2).doc",
			"C语言模拟卷含答案 (3).doc",
			"C语言模拟卷含答案 (4).doc",
			"C语言模拟卷含答案 (5).doc",
			"C语言模拟卷含答案 (6).doc",
			"C语言模拟卷含答案 (7).doc",
			"C语言要点综合.zip",
			"C语言试题汇编.doc",
			"c语言真题.zip",
			"c语言选择题.doc",
			"ugee.tablet.driver.zip",
			"全国计算机二级考试复习资料.doc",
			"全国计算机等级考试二级笔试复习资料.doc",
			"实验.zip",
			"林秋月课件.zip",
			"王林课件.zip",
			"王鸿课件.zip",
			"省考.zip",
			"省考2级.zip",
			"福建省c语言考试试题c题库选择题答案06-08(最新).doc",
			"福建省计算机二级c语言模拟卷汇总.doc",
			"福建省计算机二级c语言选择题题库.doc",
			"福建省高等学校2013年计算机二级C语言试题库.doc",
			"福建省高等学校计算机二级C语言试题库大题部分.doc",
			"计算机2级.doc",
			"计算机等级考试二级C语言超级经典400道题目[1].doc.doc",
			"谢丽聪课件.zip",
			"选择50题-第二次完善.zip",
		},
		Folders: []string{
			"c语试题",
			"王鸿",
			"省考（期末考）真题",
			"谢丽聪",
		},
	}

	data, err := json.Marshal(expectedUpYunResult)
	assert.NoError(t, err)

	testCases := []TestCase{
		{
			Name:              "GetListDirFilesSuccessfully",
			ExpectedError:     false,
			ExpectedResult:    `{"code":"10000","message":"Success","data":` + string(data) + `}`,
			ExpectUpYunResult: expectedUpYunResult,
			Path:              "/C语言",
		},
		{
			Name:              "EmptyPath",
			ExpectedError:     false,
			ExpectedResult:    `{"code":"20001","message":"参数错误, path is empty"}`,
			ExpectUpYunResult: nil,
			Path:              "",
		},
		{
			Name:              "GetListDirFilesFailed",
			ExpectedError:     true,
			ExpectedResult:    `{"code":"50001","message":"GetListDirFilesRPC: RPC called failed: wrong path"}`,
			ExpectUpYunResult: nil,
			Path:              "/C",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/paper/list", ListDirFiles)

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(rpc.GetDirFilesRPC).To(func(ctx context.Context, req *paper.ListDirFilesRequest) (resp *model.UpYunFileDir, err error) {
				if tc.ExpectedError {
					return nil, errors.New("GetListDirFilesRPC: RPC called failed: wrong path")
				}
				return tc.ExpectUpYunResult, nil
			}).Build()
			defer mockey.UnPatchAll()

			url := "/api/v1/paper/list" + "?path=" + tc.Path
			resp := ut.PerformRequest(router, consts.MethodGet, url, nil)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, tc.ExpectedResult, string(resp.Result().Body()))
		})
	}
}

func TestGetDownloadUrlForAndroid(t *testing.T) {
	type TestCase struct {
		Name           string
		ExpectedError  bool
		ExpectedResult interface{}
		Url            string
	}

	testCases := []TestCase{
		{
			Name:           "GetUrlSuccessfully",
			ExpectedError:  false,
			ExpectedResult: `{"code":2000,"data":{"url":"file://url"},"msg":"Success"}`,
			Url:            "/api/v1/paper/download?filepath=url",
		},
		{
			Name:           "GetUrlFailed",
			ExpectedError:  true,
			ExpectedResult: `{"code":50001,"data":null,"msg":"GetDownloadUrlRPC: RPC called failed: wrong filepath"}`,
			Url:            "/api/v1/paper/download?filepath=",
		},
		{
			Name:          "BindAndValidateError",
			ExpectedError: false,
			ExpectedResult: `{"code":20001,"data":null,"msg":"参数错误, 'filepath' field is a 'required' parameter,` +
				` but the request body does not have this parameter 'filepath'"}`,
			Url: "/api/v1/paper/download",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/paper/download", GetDownloadUrlForAndroid)

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(rpc.GetDownloadUrlRPC).To(func(ctx context.Context, req *paper.GetDownloadUrlRequest) (url string, err error) {
				if tc.ExpectedError {
					return "", errors.New("GetDownloadUrlRPC: RPC called failed: wrong filepath")
				}
				return "file://" + req.Filepath, nil
			}).Build()
			defer mockey.UnPatchAll()

			resp := ut.PerformRequest(router, consts.MethodGet, tc.Url, nil)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, tc.ExpectedResult, string(resp.Result().Body()))
		})
	}
}

func TestListDirFilesForAndroid(t *testing.T) {
	type TestCase struct {
		Name              string
		ExpectedError     bool
		ExpectedResult    interface{}
		ExpectUpYunResult *model.UpYunFileDir
		Path              string
	}
	basePath := "/C语言"
	expectedUpYunResult := &model.UpYunFileDir{
		BasePath: &basePath,
		Files: []string{
			"10份练习.zip",
			"200912填空题(含答案).doc",
			"200912改错题(含答案).doc",
			"200912程序题(含答案).doc",
			"201012真题C语言(含答案).doc",
			"2010年4月试题.doc",
			"2010年4月试题答案.doc",
			"2011-06改错填空编程(含答案).doc",
			"2011-06真题C语言(含答案).doc",
			"2012年福建省C语言二级考试大题及答案.doc",
			"2015C语言课件.zip",
			"2015上机题.zip",
			"2015选择题.zip",
			"C语言-张莹.zip",
			"C语言模拟卷含答案 (1).doc",
			"C语言模拟卷含答案 (2).doc",
			"C语言模拟卷含答案 (3).doc",
			"C语言模拟卷含答案 (4).doc",
			"C语言模拟卷含答案 (5).doc",
			"C语言模拟卷含答案 (6).doc",
			"C语言模拟卷含答案 (7).doc",
			"C语言要点综合.zip",
			"C语言试题汇编.doc",
			"c语言真题.zip",
			"c语言选择题.doc",
			"ugee.tablet.driver.zip",
			"全国计算机二级考试复习资料.doc",
			"全国计算机等级考试二级笔试复习资料.doc",
			"实验.zip",
			"林秋月课件.zip",
			"王林课件.zip",
			"王鸿课件.zip",
			"省考.zip",
			"省考2级.zip",
			"福建省c语言考试试题c题库选择题答案06-08(最新).doc",
			"福建省计算机二级c语言模拟卷汇总.doc",
			"福建省计算机二级c语言选择题题库.doc",
			"福建省高等学校2013年计算机二级C语言试题库.doc",
			"福建省高等学校计算机二级C语言试题库大题部分.doc",
			"计算机2级.doc",
			"计算机等级考试二级C语言超级经典400道题目[1].doc.doc",
			"谢丽聪课件.zip",
			"选择50题-第二次完善.zip",
		},
		Folders: []string{
			"c语试题",
			"王鸿",
			"省考（期末考）真题",
			"谢丽聪",
		},
	}

	data, err := json.Marshal(expectedUpYunResult)
	assert.NoError(t, err)

	testCases := []TestCase{
		{
			Name:              "GetListDirFilesSuccessfully",
			ExpectedError:     false,
			ExpectedResult:    `{"code":2000,"data":` + strings.Replace(string(data), "basePath", "base_path", 1) + `,"msg":"Success"` + `}`,
			ExpectUpYunResult: expectedUpYunResult,
			Path:              "/C语言",
		},
		{
			Name:              "EmptyPath",
			ExpectedError:     false,
			ExpectedResult:    `{"code":20001,"data":null,"msg":"参数错误, path is empty"}`,
			ExpectUpYunResult: nil,
			Path:              "",
		},
		{
			Name:              "GetListDirFilesFailed",
			ExpectedError:     true,
			ExpectedResult:    `{"code":50001,"data":null,"msg":"GetListDirFilesRPC: RPC called failed: wrong path"}`,
			ExpectUpYunResult: nil,
			Path:              "/C",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/paper/list", ListDirFilesForAndroid)

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(rpc.GetDirFilesRPC).To(func(ctx context.Context, req *paper.ListDirFilesRequest) (resp *model.UpYunFileDir, err error) {
				if tc.ExpectedError {
					return nil, errors.New("GetListDirFilesRPC: RPC called failed: wrong path")
				}
				return tc.ExpectUpYunResult, nil
			}).Build()
			defer mockey.UnPatchAll()

			url := "/api/v1/paper/list" + "?path=" + tc.Path
			resp := ut.PerformRequest(router, consts.MethodGet, url, nil)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, tc.ExpectedResult, string(resp.Result().Body()))
		})
	}
}
