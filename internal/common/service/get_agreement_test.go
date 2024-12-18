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
	"fmt"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"strconv"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func TestGetUserAgreement(t *testing.T) {
	type testCase struct {
		name           string
		mockFileResult *[]byte
		mockFileError  error
		expectedResult *[]byte
		expectedError  error
	}

	mockAgreement := []byte(`{"agreement": "This is a user agreement."}`)

	testCases := []testCase{
		{
			name:           "SuccessCase",
			mockFileResult: &mockAgreement,
			mockFileError:  nil,
			expectedResult: &mockAgreement,
			expectedError:  nil,
		},
		{
			name:           "FileNotFound",
			mockFileResult: nil,
			mockFileError:  errno.UpcloudError,
			expectedResult: nil,
			expectedError:  fmt.Errorf("CommonService.GetUserAgreement error:[" + strconv.Itoa(errno.BizFileUploadErrorCode) + "] " + errno.UpcloudError.ErrorMsg),
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock upyun.URlGetFile
			mockey.Mock(upyun.URlGetFile).To(func(filename string) (*[]byte, error) {
				return tc.mockFileResult, tc.mockFileError
			}).Build()

			// Mock upyun.JoinFileName
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// Initialize CommonService
			commonService := &CommonService{}

			// Call the method
			result, err := commonService.GetUserAgreement()

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
