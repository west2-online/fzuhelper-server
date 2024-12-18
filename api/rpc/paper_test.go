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

package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
)

func TestGetDownloadUrlRPC(t *testing.T) {
	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    string
		mockResult        string
	}

	testCases := []TestCase{
		{
			Name:              "GetDownloadUrlRPCSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    "url",
			mockResult:        "url",
		},
		{
			Name:              "GetDownloadUrlRPCError",
			expectedError:     true,
			expectedErrorInfo: errors.New("RPC call failed"),
			expectedResult:    "",
			mockResult:        "",
		},
	}

	req := &paper.GetDownloadUrlRequest{
		Filepath: "test_filepath",
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(GetDownloadUrlRPC).To(func(ctx context.Context, req *paper.GetDownloadUrlRequest) (url string, err error) {
				return tc.mockResult, tc.expectedErrorInfo
			}).Build()

			result, err := GetDownloadUrlRPC(context.Background(), req)
			if tc.expectedError {
				assert.EqualError(t, tc.expectedErrorInfo, err.Error())
				assert.Equal(t, tc.expectedResult, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetDirFilesRPC(t *testing.T) {
	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    *model.UpYunFileDir
		mockResult        *model.UpYunFileDir
	}
	basePath := "/C语言"
	expectedResult := &model.UpYunFileDir{
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

	testCases := []TestCase{
		{
			Name:              "GetDownloadUrlRPCSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
			mockResult:        expectedResult,
		},
		{
			Name:              "GetDownloadUrlRPCError",
			expectedError:     true,
			expectedErrorInfo: errors.New("RPC call failed"),
			expectedResult:    nil,
			mockResult:        nil,
		},
	}

	req := &paper.ListDirFilesRequest{
		Path: "test_path",
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(GetDirFilesRPC).To(func(ctx context.Context, req *paper.ListDirFilesRequest) (files *model.UpYunFileDir, err error) {
				return tc.mockResult, tc.expectedErrorInfo
			}).Build()

			result, err := GetDirFilesRPC(context.Background(), req)
			if tc.expectedError {
				assert.EqualError(t, tc.expectedErrorInfo, err.Error())
				assert.Equal(t, tc.expectedResult, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
