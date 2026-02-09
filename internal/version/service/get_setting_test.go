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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/internal/version/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	versionCache "github.com/west2-online/fzuhelper-server/pkg/cache/version"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func TestGetCloudSetting(t *testing.T) {
	type testCase struct {
		name                string
		mockVisitsData      *[]byte
		mockVisitsError     error
		mockUploadError     error
		mockCloudSetting    *[]byte
		mockCloudSettingErr error
		mockNoCommentJson   string
		mockNoCommentError  error
		mockCriteria        *pack.Plan
		mockPlanList        []pack.Plan
		expectResult        *[]byte
		expectError         string
	}

	mockVisits := []byte(`{"2024-12-13": 100}`)
	mockSettings := []byte(`{"Plans":[{"Name":"Test Plan","Plan":{"key":"value"}}]}`)
	mockPlanResult := []byte(`{"key":"value"}`)

	testCases := []testCase{
		{
			name:                "SuccessCase",
			mockVisitsData:      &mockVisits,
			mockVisitsError:     nil,
			mockUploadError:     nil,
			mockCloudSetting:    &mockSettings,
			mockCloudSettingErr: nil,
			mockNoCommentJson:   `{"Plans":[{"Name":"Test Plan","Plan":{"key":"value"}}]}`,
			mockNoCommentError:  nil,
			mockCriteria:        &pack.Plan{Name: strPtr("Test Plan")},
			mockPlanList:        []pack.Plan{{Name: strPtr("Test Plan"), Plan: json.RawMessage(mockPlanResult)}},
			expectResult:        &mockPlanResult,
			expectError:         "",
		},
		{
			name:                "NoMatchingPlan",
			mockVisitsData:      &mockVisits,
			mockVisitsError:     nil,
			mockUploadError:     nil,
			mockCloudSetting:    &mockSettings,
			mockCloudSettingErr: nil,
			mockNoCommentJson:   `{"Plans":[{"Name":"Other Plan","Plan":{"key":"value"}}]}`,
			mockNoCommentError:  nil,
			mockCriteria:        &pack.Plan{Name: strPtr("Non-Matching Plan")},
			mockPlanList:        []pack.Plan{{Name: strPtr("Other Plan"), Plan: json.RawMessage(mockPlanResult)}},
			expectResult:        nil,
			expectError:         "VersionService.GetCloudSetting AddVisit error",
		},
		{
			name:                "FindMatchingPlanError",
			mockVisitsData:      &mockVisits,
			mockVisitsError:     nil,
			mockUploadError:     nil,
			mockCloudSetting:    &mockSettings,
			mockCloudSettingErr: nil,
			mockNoCommentJson:   `{"Plans":[{"Name":"Other Plan","Plan":{"key":"value"}}]}`,
			mockNoCommentError:  nil,
			mockCriteria:        &pack.Plan{Name: strPtr("Test Plan")},
			mockPlanList:        []pack.Plan{{Name: strPtr("Other Plan"), Plan: json.RawMessage(mockPlanResult)}},
			expectResult:        nil,
			expectError:         "VersionService.GetCloudSetting error",
		},
		{
			name:                "URLGetFileError",
			mockVisitsData:      &mockVisits,
			mockVisitsError:     nil,
			mockUploadError:     nil,
			mockCloudSetting:    nil,
			mockCloudSettingErr: fmt.Errorf("network error"),
			mockNoCommentJson:   "",
			mockNoCommentError:  nil,
			mockCriteria:        &pack.Plan{Name: strPtr("Test Plan")},
			mockPlanList:        []pack.Plan{},
			expectResult:        nil,
			expectError:         "VersionService.GetCloudSetting error:network error",
		},
		{
			name:                "GetJSONWithoutCommentsError",
			mockVisitsData:      &mockVisits,
			mockVisitsError:     nil,
			mockUploadError:     nil,
			mockCloudSetting:    &mockSettings,
			mockCloudSettingErr: nil,
			mockNoCommentJson:   "",
			mockNoCommentError:  fmt.Errorf("json processing error"),
			mockCriteria:        &pack.Plan{Name: strPtr("Test Plan")},
			mockPlanList:        []pack.Plan{},
			expectResult:        nil,
			expectError:         "VersionService.GetCloudSetting error:json processing error",
		},
		{
			name:                "UnmarshalError",
			mockVisitsData:      &mockVisits,
			mockVisitsError:     nil,
			mockUploadError:     nil,
			mockCloudSetting:    &mockSettings,
			mockCloudSettingErr: nil,
			mockNoCommentJson:   `this is not valid json at all`,
			mockNoCommentError:  nil,
			mockCriteria:        &pack.Plan{Name: strPtr("Test Plan")},
			mockPlanList:        []pack.Plan{},
			expectResult:        nil,
			expectError:         "VersionService.GetCloudSetting error",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				CacheClient: new(cache.Cache),
			}

			// AddVisit should only fail for specific test cases
			addVisitErr := tc.mockVisitsError
			if tc.name == "NoMatchingPlan" {
				addVisitErr = fmt.Errorf("%s", tc.expectError)
			}

			mockey.Mock((*versionCache.CacheVersion).AddVisit).Return(addVisitErr).Build()

			// Mock upyun.URlGetFile for visits data
			mockey.Mock(upyun.URlGetFile).To(func(filename string) (*[]byte, error) {
				if filename == visitsFileName {
					return tc.mockVisitsData, tc.mockVisitsError
				}
				return tc.mockCloudSetting, tc.mockCloudSettingErr
			}).Build()

			// Mock upyun.JoinFileName
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// Mock upyun.URlUploadFile
			mockey.Mock(upyun.URlUploadFile).Return(tc.mockUploadError).Build()

			// Mock getJSONWithoutComments - but let it return value to test unmarshal
			if tc.mockNoCommentError != nil || tc.name != "UnmarshalError" {
				mockey.Mock(getJSONWithoutComments).Return(tc.mockNoCommentJson, tc.mockNoCommentError).Build()
			} else {
				// For UnmarshalError case, mock it to return invalid JSON without error
				mockey.Mock(getJSONWithoutComments).Return(tc.mockNoCommentJson, nil).Build()
			}

			// Mock findMatchingPlan
			mockey.Mock(findMatchingPlan).To(func(planList *[]pack.Plan, criteria *pack.Plan) (*pack.Plan, error) {
				if len(tc.mockPlanList) == 0 || tc.mockCriteria == nil || *tc.mockCriteria.Name != *tc.mockPlanList[0].Name {
					return nil, errno.NoMatchingPlanError
				}
				return &tc.mockPlanList[0], nil
			}).Build()

			versionService := NewVersionService(context.Background(), mockClientSet)

			// Call the method
			result, err := versionService.GetCloudSetting(&version.GetSettingRequest{
				Account:   tc.mockCriteria.Account,
				Version:   tc.mockCriteria.Version,
				Beta:      tc.mockCriteria.Beta,
				Phone:     tc.mockCriteria.Phone,
				IsLogin:   tc.mockCriteria.IsLogin,
				LoginType: tc.mockCriteria.LoginType,
			})

			if tc.expectError != "" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestFindMatchingPlan(t *testing.T) {
	type testCase struct {
		name        string
		planList    []pack.Plan
		criteria    *pack.Plan
		expectPlan  *pack.Plan
		expectError error
	}

	testCases := []testCase{
		{
			name: "MatchByName",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), Plan: json.RawMessage([]byte(`{"key":"value1"}`))},
				{Name: strPtr("Other.*"), Plan: json.RawMessage([]byte(`{"key":"value2"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("TestPlan")},
			expectPlan:  &pack.Plan{Name: strPtr("Test.*"), Plan: json.RawMessage([]byte(`{"key":"value1"}`))},
			expectError: nil,
		},
		{
			name: "NoMatchingPlan",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), Plan: json.RawMessage([]byte(`{"key":"value1"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("Other")},
			expectPlan:  nil,
			expectError: errno.NoMatchingPlanError,
		},
		{
			name: "AccountMismatch",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), Account: strPtr("acc.*"), Plan: json.RawMessage([]byte(`{"key":"value"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("Test123"), Account: strPtr("other")},
			expectPlan:  nil,
			expectError: errno.NoMatchingPlanError,
		},
		{
			name: "VersionMismatch",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), Version: strPtr("1\\.0.*"), Plan: json.RawMessage([]byte(`{"key":"value"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("Test123"), Version: strPtr("2.0.0")},
			expectPlan:  nil,
			expectError: errno.NoMatchingPlanError,
		},
		{
			name: "PhoneMismatch",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), Phone: strPtr("133.*"), Plan: json.RawMessage([]byte(`{"key":"value"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("Test123"), Phone: strPtr("1440000")},
			expectPlan:  nil,
			expectError: errno.NoMatchingPlanError,
		},
		{
			name: "LoginTypeMismatch",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), LoginType: strPtr("student"), Plan: json.RawMessage([]byte(`{"key":"value"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("Test123"), LoginType: strPtr("teacher")},
			expectPlan:  nil,
			expectError: errno.NoMatchingPlanError,
		},
		{
			name: "BetaMismatch",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), Beta: boolPtr(true), Plan: json.RawMessage([]byte(`{"key":"value"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("Test123"), Beta: boolPtr(false)},
			expectPlan:  nil,
			expectError: errno.NoMatchingPlanError,
		},
		{
			name: "IsLoginMismatch",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), IsLogin: boolPtr(true), Plan: json.RawMessage([]byte(`{"key":"value"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("Test123"), IsLogin: boolPtr(false)},
			expectPlan:  nil,
			expectError: errno.NoMatchingPlanError,
		},
		{
			name: "MatchByMultipleCriteria",
			planList: []pack.Plan{
				{Name: strPtr("Test.*"), Version: strPtr("1\\.0.*"), Plan: json.RawMessage([]byte(`{"key":"value1"}`))},
				{Name: strPtr("Test.*"), Version: strPtr("2\\.0.*"), Plan: json.RawMessage([]byte(`{"key":"value2"}`))},
			},
			criteria:    &pack.Plan{Name: strPtr("TestPlan"), Version: strPtr("2.0.0")},
			expectPlan:  &pack.Plan{Name: strPtr("Test.*"), Version: strPtr("2\\.0.*"), Plan: json.RawMessage([]byte(`{"key":"value2"}`))},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := findMatchingPlan(&tc.planList, tc.criteria)
			if tc.expectError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectError, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectPlan, result)
			}
		})
	}
}

func TestGetJSONWithoutComments(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		checkContains []string
		expectError   bool
	}

	testCases := []testCase{
		{
			name: "NoComments",
			input: `{
				"key": "value",
				"number": 123
			}`,
			checkContains: []string{`"key": "value"`, `"number": 123`},
			expectError:   false,
		},
		{
			name: "WithComments",
			input: `{
				"key": "value", // This is a comment
				"number": 123 // Another comment
			}`,
			checkContains: []string{`"key": "value"`, `"number": 123`},
			expectError:   false,
		},
		{
			name: "CommentsInString",
			input: `{
				"url": "http://example.com", // URL should not be affected
				"comment": "// This is not a comment"
			}`,
			checkContains: []string{`"url": "http://example.com"`, `"comment": "// This is not a comment"`},
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getJSONWithoutComments(tc.input)
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for _, contain := range tc.checkContains {
					assert.Contains(t, result, contain)
				}
			}
		})
	}
}

func TestRemoveComments(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		expectOutput string
	}

	testCases := []testCase{
		{
			name:         "NoComments",
			input:        `"key": "value"`,
			expectOutput: `"key": "value"`,
		},
		{
			name:         "WithComment",
			input:        `"key": "value" // This is a comment`,
			expectOutput: `"key": "value" `,
		},
		{
			name:         "URLNotAffected",
			input:        `"url": "http://example.com"`,
			expectOutput: `"url": "http://example.com"`,
		},
		{
			name:         "CommentInString",
			input:        `"text": "some // text" // actual comment`,
			expectOutput: `"text": "some // text" `,
		},
		{
			name:         "EmptyString",
			input:        `""`,
			expectOutput: `""`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeComments(tc.input)
			assert.Equal(t, tc.expectOutput, result)
		})
	}
}
