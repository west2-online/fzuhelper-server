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

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	launchScreenDB "github.com/west2-online/fzuhelper-server/pkg/db/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/oss"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestListImage(t *testing.T) {
	type testCase struct {
		name           string
		req            *launch_screen.ListImageRequest
		mockCheckPwd   bool
		mockDBResult   *[]model.Picture
		mockDBTotal    int64
		mockDBError    error
		expectPageNum  int
		expectPageSize int
		expectError    string
	}

	pictures := &[]model.Picture{
		{
			ID:        2024,
			Url:       "newUrl",
			Href:      "href",
			Text:      "text",
			PicType:   3,
			SType:     3,
			Frequency: 4,
			StartAt:   time.Now().Add(-24 * time.Hour),
			EndAt:     time.Now().Add(24 * time.Hour),
		},
	}

	testCases := []testCase{
		{
			name:           "ListImage_Success_DefaultPage",
			req:            &launch_screen.ListImageRequest{Secret: "secret"},
			mockCheckPwd:   true,
			mockDBResult:   pictures,
			mockDBTotal:    1,
			expectPageNum:  1,
			expectPageSize: 20,
		},
		{
			name:           "ListImage_Success_CustomPage",
			req:            &launch_screen.ListImageRequest{Secret: "secret", PageNum: new(int64(2)), PageSize: new(int64(10))},
			mockCheckPwd:   true,
			mockDBResult:   pictures,
			mockDBTotal:    12,
			expectPageNum:  2,
			expectPageSize: 10,
		},
		{
			name:           "ListImage_Success_PageSizeTooLarge",
			req:            &launch_screen.ListImageRequest{Secret: "secret", PageNum: new(int64(1)), PageSize: new(int64(1000))},
			mockCheckPwd:   true,
			mockDBResult:   pictures,
			mockDBTotal:    1,
			expectPageNum:  1,
			expectPageSize: 20,
		},
		{
			name:         "ListImage_AuthFailed",
			req:          &launch_screen.ListImageRequest{Secret: "wrong-secret"},
			mockCheckPwd: false,
			expectError:  "LaunchScreenService.ListImage error: AuthFailedError",
		},
		{
			name:         "ListImage_PageOffsetTooLarge",
			req:          &launch_screen.ListImageRequest{Secret: "secret", PageNum: new(int64(1 << 62)), PageSize: new(int64(100))},
			mockCheckPwd: true,
			expectError:  "page offset is too large",
		},
		{
			name:         "ListImage_DBError",
			req:          &launch_screen.ListImageRequest{Secret: "secret"},
			mockCheckPwd: true,
			mockDBError:  errno.BizError,
			expectError:  "LaunchScreenService.ListImage error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
				OssSet: &oss.OSSSet{
					Provider: oss.UpYunProvider,
					Upyun:    new(oss.UpYunConfig),
				},
			}
			launchScreenService := NewLaunchScreenService(context.Background(), mockClientSet)

			mockey.Mock(utils.CheckPwd).Return(tc.mockCheckPwd).Build()

			var gotPageNum, gotPageSize int
			mockey.Mock((*launchScreenDB.DBLaunchScreen).ListImage).To(func(ctx context.Context, pageNum, pageSize int) (*[]model.Picture, int64, error) {
				gotPageNum, gotPageSize = pageNum, pageSize
				return tc.mockDBResult, tc.mockDBTotal, tc.mockDBError
			}).Build()

			result, total, err := launchScreenService.ListImage(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.mockDBResult, result)
			assert.Equal(t, tc.mockDBTotal, total)
			assert.Equal(t, tc.expectPageNum, gotPageNum)
			assert.Equal(t, tc.expectPageSize, gotPageSize)
		})
	}
}
