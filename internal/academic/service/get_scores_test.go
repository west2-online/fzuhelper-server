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
	"github.com/west2-online/fzuhelper-server/pkg/base"
	meta "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	academicCache "github.com/west2-online/fzuhelper-server/pkg/cache/academic"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/jwch"
)

func TestAcademicService_GetScores(t *testing.T) {
	type testCase struct {
		name               string
		mockIsKeyExist     bool
		mockCacheReturn    []*jwch.Mark
		mockCacheError     error
		mockJwchReturn     []*jwch.Mark
		mockJwchError      error
		expectedResult     []*jwch.Mark
		expectingError     bool
		expectingCacheCall bool
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
			name:               "GetScores from cache success",
			mockIsKeyExist:     true,
			mockCacheReturn:    expectedResult,
			mockCacheError:     nil,
			expectedResult:     expectedResult,
			expectingCacheCall: true,
		},
		{
			name:               "GetScores from cache failure, fallback to GetMarks",
			mockIsKeyExist:     true,
			mockCacheReturn:    nil,
			mockCacheError:     fmt.Errorf("Get scores info fail"),
			mockJwchReturn:     expectedResult,
			mockJwchError:      nil,
			expectedResult:     expectedResult,
			expectingError:     true,
			expectingCacheCall: true,
		},
		{
			name:               "GetScores from GetMarks success",
			mockIsKeyExist:     false,
			mockJwchReturn:     expectedResult,
			mockJwchError:      nil,
			expectedResult:     expectedResult,
			expectingError:     false,
			expectingCacheCall: false,
		},
		{
			name:               "GetScores from GetMarks failure",
			mockIsKeyExist:     false,
			mockJwchReturn:     nil,
			mockJwchError:      fmt.Errorf("Get scores info fail"),
			expectedResult:     nil,
			expectingError:     true,
			expectingCacheCall: false,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockClientSet.CacheClient = new(cache.Cache)
			mockey.Mock((*jwch.Student).GetMarks).Return(tc.mockJwchReturn, tc.mockJwchError).Build()
			mockey.Mock(meta.GetLoginData).To(func(ctx context.Context) (*model.LoginData, error) {
				return &model.LoginData{
					Id:      "1111111111111111111111111111111111",
					Cookies: "",
				}, nil
			}).Build()
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.mockIsKeyExist).Build()
			if tc.expectingCacheCall {
				mockey.Mock((*academicCache.CacheAcademic).GetScoresCache).
					Return(tc.mockCacheReturn, tc.mockCacheError).
					Build()
			}
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()
			academicService := NewAcademicService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := academicService.GetScores(&model.LoginData{
				Id:      "123456789",
				Cookies: "cookie1=value1;cookie2=value2",
			})
			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Get scores info fail")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
