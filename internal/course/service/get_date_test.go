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
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestGetLocateDate(t *testing.T) {
	type testCase struct {
		name          string
		cacheExist    bool
		cacheGetError error
		jwchReturn    *jwch.LocateDate
		jwchError     error
		expectYear    string
		expectWeek    string
		expectTerm    string
		expectError   string
	}

	testCases := []testCase{
		{
			name:       "cache exist success",
			cacheExist: true,
			expectYear: "2025",
			expectWeek: "10",
			expectTerm: "202501",
		},
		{
			name:          "cache exist but get error -> fallback",
			cacheExist:    true,
			cacheGetError: assert.AnError,
			jwchReturn:    &jwch.LocateDate{Year: "2025", Week: "11", Term: "202501"},
			expectYear:    "2025",
			expectWeek:    "11",
			expectTerm:    "202501",
		},
		{
			name:       "cache miss -> fetch from jwch",
			jwchReturn: &jwch.LocateDate{Year: "2025", Week: "12", Term: "202501"},
			expectYear: "2025",
			expectWeek: "12",
			expectTerm: "202501",
		},
		{
			name:        "jwch error",
			jwchError:   assert.AnError,
			expectError: "Get locate date fail",
		},
	}

	mockStudent := &jwch.Student{}
	mockLoginData := &model.LoginData{Id: "052106112", Cookies: "test_cookie"}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			// Mock cache IsKeyExist
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()

			if tc.cacheExist {
				// Mock GetDateCache when cache exists
				mockey.Mock((*coursecache.CacheCourse).GetDateCache).Return(&model.LocateDate{
					Year: tc.expectYear,
					Week: tc.expectWeek,
					Term: tc.expectTerm,
					Date: time.Now().In(constants.ChinaTZ).Format(time.DateTime),
				}, tc.cacheGetError).Build()
			}

			// Mock jwch.NewStudent() and GetLocateDate
			mockey.Mock(jwch.NewStudent).Return(mockStudent).Build()
			mockey.Mock((*jwch.Student).GetLocateDate).Return(tc.jwchReturn, tc.jwchError).Build()

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))

			result, err := courseService.GetLocateDate()
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.Equal(t, tc.expectYear, result.Year)
				assert.Equal(t, tc.expectWeek, result.Week)
				assert.Equal(t, tc.expectTerm, result.Term)
				assert.NotEmpty(t, result.Date)
			}
		})
	}
}
