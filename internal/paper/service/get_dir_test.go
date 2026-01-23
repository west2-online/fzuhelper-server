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
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	paperCache "github.com/west2-online/fzuhelper-server/pkg/cache/paper"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func TestGetDir(t *testing.T) {
	type testCase struct {
		name string // 用例名
		// 控制返回值与mock函数行为
		mockIsCacheExist bool

		mockCacheReturn *model.UpYunFileDir

		mockUpYunReturn *model.UpYunFileDir
		// 期望输出
		expectedResult *model.UpYunFileDir
		// 此用例是否报错
		expectingError bool
		// 期望错误信息
		expectedErrorInfo error
		// 成功获取数据
		mockIsGetInfo bool
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

	// ResultWithIgnoredFolders为返回结果中包含被过滤文件夹的情况
	resultWithIgnoredFolders := &model.UpYunFileDir{
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
			"upyun_storage_log_AhYIBW15",
			"test",
		},
	}

	testCases := []testCase{
		{
			name:              "GetDirFromUpYunWithCacheNotFound",
			expectedResult:    expectedResult,
			mockIsCacheExist:  false,
			mockCacheReturn:   nil,
			expectingError:    false,
			mockUpYunReturn:   expectedResult,
			expectedErrorInfo: nil,
			mockIsGetInfo:     true,
		},
		{
			name:              "GetDirFromCache",
			expectedResult:    expectedResult,
			mockIsCacheExist:  true,
			mockCacheReturn:   expectedResult,
			expectingError:    false,
			mockUpYunReturn:   nil,
			expectedErrorInfo: nil,
			mockIsGetInfo:     true,
		},
		{
			name:              "GetDirError",
			expectedResult:    nil,
			mockIsCacheExist:  false,
			mockCacheReturn:   nil,
			expectingError:    true,
			mockUpYunReturn:   nil,
			expectedErrorInfo: errors.New("failed to get info from upyun"),
			mockIsGetInfo:     false,
		},
		{
			name:              "SetDataCacheFailed",
			expectedResult:    expectedResult,
			mockIsCacheExist:  false,
			mockCacheReturn:   nil,
			expectingError:    true,
			mockUpYunReturn:   expectedResult,
			expectedErrorInfo: errors.New("failed to set data in cache"),
			mockIsGetInfo:     true,
		},
		{
			name:              "GetCacheDataFailed",
			expectedResult:    nil,
			mockIsCacheExist:  true,
			mockCacheReturn:   nil,
			expectingError:    true,
			mockUpYunReturn:   nil,
			expectedErrorInfo: errors.New("failed to get data from cache"),
			mockIsGetInfo:     false,
		},
		{
			name:              "FilterIgnoredFolders",
			expectedResult:    expectedResult,
			mockIsCacheExist:  false,
			mockCacheReturn:   nil,
			expectingError:    false,
			mockUpYunReturn:   resultWithIgnoredFolders,
			expectedErrorInfo: nil,
			mockIsGetInfo:     true,
		},
	}

	req := &paper.ListDirFilesRequest{
		Path: "/C语言",
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockClientSet.CacheClient = new(cache.Cache)
			paperService := NewPaperService(context.Background(), mockClientSet)

			mockey.Mock(((*paperCache.CachePaper).GetFileDirKey)).To(func(path string) string {
				return path
			}).Build()
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.mockIsCacheExist
			}).Build()
			mockey.Mock((*paperCache.CachePaper).GetFileDirCache).To(func(ctx context.Context, key string) (bool, *model.UpYunFileDir, error) {
				if tc.name == "GetCacheDataFailed" {
					return false, nil, tc.expectedErrorInfo
				}
				return true, tc.mockCacheReturn, nil
			}).Build()
			mockey.Mock(upyun.GetDir).To(func(path string) (*model.UpYunFileDir, error) {
				if tc.mockIsGetInfo {
					return tc.mockUpYunReturn, nil
				}
				return tc.mockUpYunReturn, tc.expectedErrorInfo
			}).Build()
			mockey.Mock((*paperCache.CachePaper).SetFileDirCache).To(func(ctx context.Context, key string, dir model.UpYunFileDir) error {
				return tc.expectedErrorInfo
			}).Build()

			ret, result, err := paperService.GetDir(req)
			if tc.expectingError {
				if tc.mockIsGetInfo {
					assert.ErrorIs(t, err, tc.expectedErrorInfo)
				} else {
					assert.EqualError(t, err, "service.GetDir: get dir info failed: "+tc.expectedErrorInfo.Error())
				}
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.mockIsGetInfo, ret)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.mockIsGetInfo, ret)
			}
		})
	}
}
