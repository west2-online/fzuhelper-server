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
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/jwch"
)

func TestGetTermList(t *testing.T) {
	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    *jwch.SchoolCalendar
	}

	expectedResult := &jwch.SchoolCalendar{
		CurrentTerm: "202401",
		Terms: []jwch.CalTerm{
			{
				TermId:     "2024012024082620250117",
				SchoolYear: "2024",
				Term:       "202401",
				StartDate:  "2024-08-26",
				EndDate:    "2025-01-17",
			},
			{
				TermId:     "2024022025022420250704",
				SchoolYear: "2024",
				Term:       "202402",
				StartDate:  "2025-02-24",
				EndDate:    "2025-07-04",
			},
		},
	}
	defer mockey.UnPatchAll()
	testCases := []TestCase{
		{
			Name:              "GetTermListSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
		},
		{
			Name:              "GetTermListError",
			expectedError:     true,
			expectedErrorInfo: errors.New("get term list failed"),
			expectedResult:    nil,
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock((*jwch.Student).GetSchoolCalendar).To(func() (*jwch.SchoolCalendar, error) {
				return tc.expectedResult, tc.expectedErrorInfo
			}).Build()
			commonService := NewTermService(context.Background())
			result, err := commonService.GetTermList()
			if tc.expectedError {
				assert.EqualError(t, tc.expectedErrorInfo, err.Error())
				assert.Nil(t, result)
			} else {
				assert.Nil(t, tc.expectedErrorInfo, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetTerm(t *testing.T) {
	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    *jwch.CalTermEvents
	}

	expectedResult := &jwch.CalTermEvents{
		TermId:     "201501",
		Term:       "201501",
		SchoolYear: "2015",
		Events: []jwch.CalTermEvent{
			{
				Name:      "学生注册",
				StartDate: "2015-08-29",
				EndDate:   "2015-08-30",
			},
			{
				Name:      "学生补考",
				StartDate: "2015-08-29",
				EndDate:   "2015-09-06",
			},
			{
				Name:      "正式上课",
				StartDate: "2015-08-31",
				EndDate:   "2015-08-31",
			},
			{
				Name:      "新生报到",
				StartDate: "2015-09-07",
				EndDate:   "2015-09-07",
			},
			{
				Name:      "校运会",
				StartDate: "2015-11-12",
				EndDate:   "2015-11-14",
			},
			{
				Name:      "期末考试",
				StartDate: "2016-01-16",
				EndDate:   "2016-01-22",
			},
			{
				Name:      "寒假",
				StartDate: "2016-01-23",
				EndDate:   "2016-02-28",
			},
		},
	}

	testCases := []TestCase{
		{
			Name:              "GetTermSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    expectedResult,
		},
		{
			Name:              "GetTermError",
			expectedError:     true,
			expectedErrorInfo: errors.New("get term events failed"),
			expectedResult:    nil,
		},
	}

	defer mockey.UnPatchAll()
	req := &common.TermRequest{Term: "201501"}
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock((*jwch.Student).GetTermEvents).To(func(termId string) (*jwch.CalTermEvents, error) {
				return tc.expectedResult, tc.expectedErrorInfo
			}).Build()
			commonService := NewTermService(context.Background())
			result, err := commonService.GetTerm(req)
			if tc.expectedError {
				assert.EqualError(t, tc.expectedErrorInfo, err.Error())
				assert.Nil(t, result)
			} else {
				assert.Nil(t, tc.expectedErrorInfo, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
