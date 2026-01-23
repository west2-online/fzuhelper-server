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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/internal/version/pack"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func TestAndroidGetVersion(t *testing.T) {
	type testCase struct {
		name              string  // 测试用例名称
		mockReleaseBytes  *[]byte // mock返回的Release版本JSON数据
		mockBetaBytes     *[]byte // mock返回的Beta版本JSON数据
		mockReleaseError  error   // mock返回的Release版本错误
		mockBetaError     error   // mock返回的Beta版本错误
		expectingError    bool    // 是否期望抛出错误
		expectedErrorInfo string  // 期望的错误信息
	}

	mockReleaseVersion := &pack.Version{Url: "http://example.com/release.apk", Version: "2.0.0"}
	mockReleaseBytes, _ := json.Marshal(mockReleaseVersion)

	mockBetaVersion := &pack.Version{Url: "http://example.com/beta.apk", Version: "2.1.0"}
	mockBetaBytes, _ := json.Marshal(mockBetaVersion)

	testCases := []testCase{
		{
			name:             "SuccessCase",
			mockReleaseBytes: &mockReleaseBytes,
			mockBetaBytes:    &mockBetaBytes,
			mockReleaseError: nil,
			mockBetaError:    nil,
			expectingError:   false,
		},
		{
			name:              "ReleaseFileNotFound",
			mockReleaseBytes:  nil,
			mockBetaBytes:     &mockBetaBytes,
			mockReleaseError:  fmt.Errorf("file not found"),
			mockBetaError:     nil,
			expectingError:    true,
			expectedErrorInfo: "VersionService.AndroidGetVersion.GetReleaseVersion error:file not found",
		},
		{
			name:              "BetaFileNotFound",
			mockReleaseBytes:  &mockReleaseBytes,
			mockBetaBytes:     nil,
			mockReleaseError:  nil,
			mockBetaError:     fmt.Errorf("file not found"),
			expectingError:    true,
			expectedErrorInfo: "VersionService.AndroidGetVersion.GetBetaVersion error:file not found",
		},
		{
			name:              "ReleaseUnmarshalError",
			mockReleaseBytes:  func() *[]byte { b := []byte("invalid json"); return &b }(),
			mockBetaBytes:     &mockBetaBytes,
			mockReleaseError:  nil,
			mockBetaError:     nil,
			expectingError:    true,
			expectedErrorInfo: "VersionService.AndroidGetVersion.GetReleaseVersion error",
		},
		{
			name:              "BetaUnmarshalError",
			mockReleaseBytes:  &mockReleaseBytes,
			mockBetaBytes:     func() *[]byte { b := []byte("invalid json"); return &b }(),
			mockReleaseError:  nil,
			mockBetaError:     nil,
			expectingError:    true,
			expectedErrorInfo: "VersionService.AndroidGetVersion.GetBetaVersion error",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(upyun.URlGetFile).To(func(filename string) (*[]byte, error) {
				if filename == releaseVersionFileName {
					return tc.mockReleaseBytes, tc.mockReleaseError
				}
				return tc.mockBetaBytes, tc.mockBetaError
			}).Build()
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			urlService := &VersionService{}

			release, beta, err := urlService.AndroidGetVersion()

			if tc.expectingError {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorInfo)
				assert.Nil(t, release)
				assert.Nil(t, beta)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, mockReleaseVersion, release)
				assert.Equal(t, mockBetaVersion, beta)
			}
		})
	}
}
