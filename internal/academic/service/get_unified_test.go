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

	"github.com/west2-online/jwch"
)

func TestAcademicService_GetUnifiedExam(t *testing.T) {
	type testCase struct {
		name             string
		mockCETReturn    []*jwch.UnifiedExam
		mockJSError      error
		mockJSReturn     []*jwch.UnifiedExam
		mockError        error
		expectedResult   []*jwch.UnifiedExam
		expectingError   bool
		expectedErrorMsg string
	}

	cetExam := []*jwch.UnifiedExam{
		{
			Name:  "CET-4",
			Score: "520",
			Term:  "2021年12月",
		},
	}
	jsExam := []*jwch.UnifiedExam{
		{
			Name:  "JS",
			Score: "90",
			Term:  "2022年6月",
		},
	}

	testCases := []testCase{
		{
			name:           "GetUnifiedExamSuccess",
			mockCETReturn:  cetExam,
			mockJSReturn:   jsExam,
			mockError:      nil,
			expectedResult: append(cetExam, jsExam...),
			expectingError: false,
		},
		{
			name:             "GetCETFailure",
			mockCETReturn:    nil,
			mockJSError:      fmt.Errorf("Get cet info fail"),
			expectedResult:   nil,
			expectingError:   true,
			expectedErrorMsg: "Get cet info fail",
		},
		{
			name:             "GetJSFailure",
			mockCETReturn:    cetExam,
			mockJSReturn:     nil,
			mockError:        fmt.Errorf("Get js info fail"),
			expectedResult:   nil,
			expectingError:   true,
			expectedErrorMsg: "Get js info fail",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockey.PatchConvey(tc.name, t, func() {
				mockey.Mock((*jwch.Student).GetCET).Return(tc.mockCETReturn, tc.mockError).Build()
				mockey.Mock((*jwch.Student).GetJS).Return(tc.mockJSReturn, tc.mockJSError).Build()
				academicService := NewAcademicService(context.Background())
				result, err := academicService.GetUnifiedExam()
				if tc.expectingError {
					assert.Nil(t, result)
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedErrorMsg)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedResult, result)
				}
			})
		})
	}
}
