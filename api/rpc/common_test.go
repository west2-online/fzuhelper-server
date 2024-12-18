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

package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestGetTermsListRPC(t *testing.T) {
	f := func(str string) *string {
		return &str
	}
	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    *model.TermList
		mockResult        *model.TermList
	}

	expectedResult := &model.TermList{
		CurrentTerm: f("202401"),
		Terms: []*model.Term{
			{
				TermId:     f("2024012024082620250117"),
				SchoolYear: f("2024"),
				Term:       f("202401"),
				StartDate:  f("2024-08-26"),
				EndDate:    f("2025-01-17"),
			},
			{
				TermId:     f("2024022025022420250704"),
				SchoolYear: f("2024"),
				Term:       f("202402"),
				StartDate:  f("2025-02-24"),
				EndDate:    f("2025-07-04"),
			},
		},
	}

	testCases := []TestCase{
		{
			Name:              "GetTermsListRPCSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
			mockResult:        expectedResult,
		},
		{
			Name:              "GetTermsListRPCError",
			expectedError:     true,
			expectedErrorInfo: errors.New("RPC call failed"),
			expectedResult:    nil,
			mockResult:        nil,
		},
	}

	req := &common.TermListRequest{}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(GetTermsListRPC).To(func(ctx context.Context, req *common.TermListRequest) (*model.TermList, error) {
				return tc.mockResult, tc.expectedErrorInfo
			}).Build()

			result, err := GetTermsListRPC(context.Background(), req)
			if tc.expectedError {
				assert.EqualError(t, tc.expectedErrorInfo, err.Error())
				assert.Equal(t, tc.expectedResult, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetTermRPC(t *testing.T) {
	f := func(str string) *string {
		return &str
	}
	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    *model.TermInfo
		mockResult        *model.TermInfo
	}

	expectedResult := &model.TermInfo{
		TermId:     f("201501"),
		Term:       f("201501"),
		SchoolYear: f("2015"),
		Events: []*model.TermEvent{
			{
				Name:      f("学生注册"),
				StartDate: f("2015-08-29"),
				EndDate:   f("2015-08-30"),
			},
			{
				Name:      f("学生补考"),
				StartDate: f("2015-08-29"),
				EndDate:   f("2015-09-06"),
			},
			{
				Name:      f("正式上课"),
				StartDate: f("2015-08-31"),
				EndDate:   f("2015-08-31"),
			},
			{
				Name:      f("新生报到"),
				StartDate: f("2015-09-07"),
				EndDate:   f("2015-09-07"),
			},
			{
				Name:      f("校运会"),
				StartDate: f("2015-11-12"),
				EndDate:   f("2015-11-14"),
			},
			{
				Name:      f("期末考试"),
				StartDate: f("2016-01-16"),
				EndDate:   f("2016-01-22"),
			},
			{
				Name:      f("寒假"),
				StartDate: f("2016-01-23"),
				EndDate:   f("2016-02-28"),
			},
		},
	}

	testCases := []TestCase{
		{
			Name:              "GetTermsListRPCSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
			mockResult:        expectedResult,
		},
		{
			Name:              "GetTermsListRPCError",
			expectedError:     true,
			expectedErrorInfo: errors.New("RPC call failed"),
			expectedResult:    nil,
			mockResult:        nil,
		},
	}

	req := &common.TermRequest{
		Term: "201501",
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(GetTermRPC).To(func(ctx context.Context, req *common.TermRequest) (*model.TermInfo, error) {
				return tc.mockResult, tc.expectedErrorInfo
			}).Build()

			result, err := GetTermRPC(context.Background(), req)
			if tc.expectedError {
				assert.EqualError(t, tc.expectedErrorInfo, err.Error())
				assert.Equal(t, tc.expectedResult, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
