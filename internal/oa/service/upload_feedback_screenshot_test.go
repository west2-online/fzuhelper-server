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

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/cos"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type screenshotCOSRepo struct {
	category      string
	suffix        string
	uploaded      []byte
	generateError error
	uploadError   error
}

func (r *screenshotCOSRepo) GenerateFileName(category, suffix string) (string, string, error) {
	r.category = category
	r.suffix = suffix
	return "https://img.example.com/feedback/img/1." + suffix, "/feedback/img/1." + suffix, r.generateError
}

func (r *screenshotCOSRepo) Upload(file []byte, _ string) error {
	r.uploaded = append([]byte(nil), file...)
	return r.uploadError
}

func TestUploadFeedbackScreenshot(t *testing.T) {
	jpeg := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00}
	png := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}

	type testCase struct {
		name          string
		file          []byte
		generateError error
		uploadError   error
		expectSuffix  string
		expectError   string
	}
	testCases := []testCase{
		{name: "upload jpeg", file: jpeg, expectSuffix: "jpg"},
		{name: "upload png", file: png, expectSuffix: "png"},
		{name: "empty file", expectError: "反馈截图不能为空"},
		{name: "file too large", file: make([]byte, constants.FeedbackScreenshotMaxSize+1), expectError: "不能超过 5MB"},
		{name: "invalid image", file: []byte("not an image"), expectError: "仅支持 JPEG 和 PNG"},
		{name: "generate name error", file: jpeg, generateError: errno.InternalServiceError, expectError: "生成反馈截图名称失败"},
		{name: "upload error", file: jpeg, uploadError: errno.UpcloudError, expectError: "上传反馈截图失败"},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			repo := &screenshotCOSRepo{generateError: tc.generateError, uploadError: tc.uploadError}
			service := &OAService{ctx: context.Background(), cosClient: repo}

			url, err := service.UploadFeedbackScreenshot(tc.file)
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Empty(t, url)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, url)
			assert.Equal(t, cos.FeedbackImageCategory, repo.category)
			assert.Equal(t, tc.expectSuffix, repo.suffix)
			assert.Equal(t, tc.file, repo.uploaded)
		})
	}

	mockey.PatchConvey("COS client injected through client set", t, func() {
		repo := &screenshotCOSRepo{}
		service := NewOAService(context.Background(), "", nil, &base.ClientSet{FeedbackCOSClient: repo})
		_, err := service.UploadFeedbackScreenshot(jpeg)
		assert.NoError(t, err)
		assert.Equal(t, jpeg, repo.uploaded)
	})

	mockey.PatchConvey("cos not initialized", t, func() {
		service := &OAService{ctx: context.Background()}
		_, err := service.UploadFeedbackScreenshot(jpeg)
		assert.ErrorContains(t, err, "反馈文件存储未初始化")
	})
}
