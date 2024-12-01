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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestGenerateDownloadUrl(t *testing.T) {
	const basePath = "http://fzuhelper-paper-cos.test.upcdn.net/"
	const filePath = "/C语言/10份练习.zip"

	type testCase struct {
		name           string
		expectedResult interface{}
	}

	expectedResult := basePath + utils.UriEncode(filePath) + "?_upt=newUrl"

	testCases := []testCase{
		{
			name:           "GetDownloadUrl",
			expectedResult: expectedResult,
		},
	}

	req := &paper.GetDownloadUrlRequest{Filepath: filePath}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockClientSet.CacheClient = new(cache.Cache)
			paperService := NewPaperService(context.Background(), mockClientSet)

			mockey.Mock(upyun.GetDownloadUrl).To(func(uri string) (string, error) {
				return expectedResult, nil
			}).Build()

			result, err := paperService.GetDownloadUrl(req)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
