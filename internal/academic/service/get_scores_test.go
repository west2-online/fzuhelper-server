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

func TestAcademicService_GetScores(t *testing.T) {
	type testCase struct {
		name           string
		mockReturn     []*jwch.Mark
		mockError      error
		expectedResult []*jwch.Mark
		expectingError bool
	}

	expectedResult := []*jwch.Mark{
		{
			Name:    "Mathematics",
			Score:   "90",
			Credits: "4.0",
			GPA:     "3.9",
		},
		{
			Name:    "Physics",
			Score:   "85",
			Credits: "3.0",
			GPA:     "3.6",
		},
	}

	testCases := []testCase{
		{
			name:           "GetScoresSuccess",
			mockReturn:     expectedResult,
			mockError:      nil,
			expectedResult: expectedResult,
		},
		{
			name:           "GetScoresFailure",
			mockReturn:     nil,
			mockError:      fmt.Errorf("Get scores info fail"),
			expectedResult: nil,
			expectingError: true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockey.PatchConvey(tc.name, t, func() {
				mockey.Mock((*jwch.Student).GetMarks).Return(tc.mockReturn, tc.mockError).Build()
				academicService := NewAcademicService(context.Background())
				result, err := academicService.GetScores()
				if tc.expectingError {
					assert.Nil(t, result)
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "Get scores info fail")
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedResult, result)
				}
			})
		})
	}
}
