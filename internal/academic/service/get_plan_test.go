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
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	meta "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/jwch"
)

func TestGetPlan(t *testing.T) {
	type testCase struct {
		name           string
		mockFileResult *[]byte
		mockUrl        string
		mockFileError  error
		expectedResult string
		expectedError  error
	}
	mockUrl := "https://www.example.com&id=123456789"
	cutUrl, _, _ := strings.Cut(mockUrl, "&id")
	mockHtml := []byte(`body { background-color: #fff; }`)
	testCases := []testCase{
		{
			name:           "SuccessCase",
			mockFileResult: &mockHtml,
			mockUrl:        mockUrl,
			mockFileError:  nil,
			expectedResult: cutUrl,
			expectedError:  nil,
		},
		{
			name:           "NotFound",
			mockFileResult: nil,
			mockUrl:        "",
			mockFileError:  fmt.Errorf("%s", "cultivate plan not found"),
			expectedResult: "",
			expectedError: fmt.Errorf("%s", strings.Join([]string{
				"AcademicService.GetPlan",
			}, "")),
		},
	}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*jwch.Student).WithLoginData).To(func(identifier string, cookies []*http.Cookie) *jwch.Student {
				return jwch.NewStudent()
			}).Build()
			mockey.Mock(meta.GetLoginData).To(func(ctx context.Context) (*model.LoginData, error) {
				return &model.LoginData{
					Id:      "123456789",
					Cookies: "",
				}, nil
			}).Build()
			mockey.Mock((*jwch.Student).GetCultivatePlan).To(func() (string, error) {
				return tc.mockUrl, tc.mockFileError
			}).Build()
			mockey.Mock(getHtmlSource).To(func() (*[]byte, error) {
				return tc.mockFileResult, tc.mockFileError
			}).Build()
			academicService := AcademicService{}
			result, err := academicService.GetPlan()
			if tc.expectedError != nil {
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				// fmt.Println(string(*result))
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
