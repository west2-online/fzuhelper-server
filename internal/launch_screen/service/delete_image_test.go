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
)

func TestDeleteImage(t *testing.T) {
	type testCase struct {
		name            string
		mockReturn      interface{}
		mockCloudReturn interface{}
		expectResult    interface{}
		expectError     bool
	}

	expectedResult := &model.Picture{
		ID:         2024,
		Url:        "newUrl",
		Href:       "href",
		Text:       "text",
		PicType:    3,
		ShowTimes:  0,
		PointTimes: 0,
		Duration:   0,
		StartAt:    time.Now().Add(-24 * time.Hour),
		EndAt:      time.Now().Add(24 * time.Hour),
		StartTime:  0,
		EndTime:    24,
		SType:      3,
		Frequency:  4,
		Regex:      "{\"device\": \"android,ios\", \"student_id\": \"102301517,102301544\"}",
	}

	testCases := []testCase{
		{
			name:         "AddPointTime",
			mockReturn:   expectedResult,
			expectResult: expectedResult,
		},
		{
			name:            "cloudFail",
			mockReturn:      expectedResult,
			mockCloudReturn: errno.UpcloudError,
			expectError:     true,
		},
		{
			name:        "DeleteImage error",
			expectError: true,
		},
	}

	req := &launch_screen.DeleteImageRequest{
		PictureId: 2024,
	}

	defer mockey.UnPatchAll() // 撤销所有mock操作，不会影响其他测试
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

			mockey.Mock((*launchScreenDB.DBLaunchScreen).DeleteImage).To(func(ctx context.Context, id int64) (*model.Picture, error) {
				if tc.name == "DeleteImage error" {
					return nil, errno.BizError
				}
				pic, ok := tc.mockReturn.(*model.Picture)
				if !ok {
					return nil, errno.BizError
				}
				return pic, nil
			}).Build()

			mockey.Mock(mockey.GetMethod(launchScreenService.ossClient, "GetRemotePathFromUrl")).Return(expectedResult.Url).Build()
			mockey.Mock(mockey.GetMethod(launchScreenService.ossClient, "DeleteImg")).Return(tc.mockCloudReturn).Build()

			err := launchScreenService.DeleteImage(req.PictureId)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
