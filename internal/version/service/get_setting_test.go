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
		expectedResult      *[]byte
		expectedError       error
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
			expectedResult:      &mockPlanResult,
			expectedError:       nil,
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
			expectedResult:      nil,
			expectedError:       fmt.Errorf("VersionService.GetCloudSetting AddVisit error"),
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*versionCache.CacheVersion).AddVisit).Return(tc.expectedError).Build()
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
			mockey.Mock(upyun.URlUploadFile).To(func(data []byte, filename string) error {
				return tc.mockUploadError
			}).Build()

			// Mock getJSONWithoutComments
			mockey.Mock(getJSONWithoutComments).To(func(json string) (string, error) {
				return tc.mockNoCommentJson, tc.mockNoCommentError
			}).Build()

			// Mock findMatchingPlan
			mockey.Mock(findMatchingPlan).To(func(planList *[]pack.Plan, criteria *pack.Plan) (*pack.Plan, error) {
				if len(tc.mockPlanList) == 0 || tc.mockCriteria == nil || *tc.mockCriteria.Name != *tc.mockPlanList[0].Name {
					return nil, errno.NoMatchingPlanError
				}
				return &tc.mockPlanList[0], nil
			}).Build()

			mockClientSet := new(base.ClientSet)
			mockClientSet.CacheClient = new(cache.Cache)
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

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
