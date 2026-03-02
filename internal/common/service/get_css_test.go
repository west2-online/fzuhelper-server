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
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func TestGetCSS(t *testing.T) {
	type testCase struct {
		name           string
		mockFileResult *[]byte
		mockFileError  error
		expectResult   *[]byte
		expectError    string
	}

	mockCSS := []byte(`body { background-color: #fff; }`)

	testCases := []testCase{
		{
			name:           "SuccessCase",
			mockFileResult: &mockCSS,
			expectResult:   &mockCSS,
		},
		{
			name:          "FileNotFound",
			mockFileError: errno.UpcloudError,
			expectError:   errno.UpcloudError.ErrorMsg,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{}

			// Mock upyun.URlGetFile
			mockey.Mock(upyun.URlGetFile).Return(tc.mockFileResult, tc.mockFileError).Build()
			// Mock upyun.JoinFileName
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// Initialize CommonService
			commonService := NewCommonService(context.Background(), mockClientSet)
			// Call the method
			result, err := commonService.GetCSS()

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
