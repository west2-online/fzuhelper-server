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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	meta "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/jwch"
)

func TestAcademicService_GetGPA(t *testing.T) {
	type testCase struct {
		name           string
		mockReturn     *jwch.GPABean
		mockError      error
		expectedResult *jwch.GPABean
		expectingError bool
	}

	expectedResult := &jwch.GPABean{
		Time: "2023-06-01",
		Data: []jwch.GPAData{
			{
				Type:  "Mathematics",
				Value: "4.0",
			},
			{
				Type:  "Physics",
				Value: "3.5",
			},
		},
	}

	testCases := []testCase{
		{
			name:           "GetGPASuccess",
			mockReturn:     expectedResult,
			mockError:      nil,
			expectedResult: expectedResult,
		},
		{
			name:           "GetGPAFailure",
			mockReturn:     nil,
			mockError:      fmt.Errorf("get gpa info fail"),
			expectedResult: nil,
			expectingError: true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*jwch.Student).GetGPA).Return(tc.mockReturn, tc.mockError).Build()
			mockey.Mock(meta.GetLoginData).To(func(ctx context.Context) (*model.LoginData, error) {
				return &model.LoginData{
					Id:      "1111111111111111111111111111111111",
					Cookies: "",
				}, nil
			}).Build()
			academicService := NewAcademicService(context.Background())
			result, err := academicService.GetGPA()
			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Get gpa info fail")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
