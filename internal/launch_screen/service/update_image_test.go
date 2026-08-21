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
	"strconv"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

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

func TestUpdateImagePath(t *testing.T) {
	type testCase struct {
		name             string
		mockIsExist      bool
		mockOriginReturn interface{}
		mockCloudReturn  interface{}
		mockReturn       interface{}
		expectResult     interface{}
		expectError      bool
	}

	origin := &model.Picture{
		ID:         2024,
		Url:        "oldUrl",
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
			name:             "UpdateImagePath",
			mockIsExist:      true,
			mockOriginReturn: origin,
			mockCloudReturn:  nil,
			mockReturn:       expectedResult,
			expectResult:     expectedResult,
		},
		{
			name:             "LaunchScreenNotExist",
			mockIsExist:      false,
			mockOriginReturn: gorm.ErrRecordNotFound,
			mockCloudReturn:  nil,
			mockReturn:       nil,
			expectResult:     nil,
			expectError:      true,
		},
		{
			name:             "cloudFail",
			mockIsExist:      true,
			mockCloudReturn:  errno.UpcloudError,
			mockOriginReturn: origin,
			mockReturn:       expectedResult,
			expectResult:     nil,
			expectError:      true,
		},
		{
			name:             "GetImageFileType error",
			mockIsExist:      true,
			mockOriginReturn: origin,
			mockCloudReturn:  nil,
			mockReturn:       expectedResult,
			expectResult:     nil,
			expectError:      true,
		},
		{
			name:             "GenerateImgName error",
			mockIsExist:      true,
			mockOriginReturn: origin,
			mockCloudReturn:  nil,
			mockReturn:       expectedResult,
			expectResult:     nil,
			expectError:      true,
		},
		{
			name:             "UploadImg error",
			mockIsExist:      true,
			mockOriginReturn: origin,
			mockCloudReturn:  nil,
			mockReturn:       expectedResult,
			expectResult:     nil,
			expectError:      true,
		},
	}

	req := &launch_screen.ChangeImageRequest{}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
				SFClient:    new(utils.Snowflake),
				OssSet: &oss.OSSSet{
					Provider: oss.UpYunProvider,
					Upyun:    new(oss.UpYunConfig),
				},
			}
			launchScreenService := NewLaunchScreenService(context.Background(), mockClientSet)

			if tc.mockIsExist {
				mockey.Mock((*launchScreenDB.DBLaunchScreen).GetImageById).Return(tc.mockOriginReturn, nil).Build()
			} else {
				mockey.Mock((*launchScreenDB.DBLaunchScreen).GetImageById).Return(nil, tc.mockOriginReturn).Build()
			}
			mockey.Mock(mockey.GetMethod(launchScreenService.ossClient, "DeleteImg")).Return(tc.mockCloudReturn).Build()

			// 根据测试用例设置UploadImg的mock返回值
			mockey.Mock(mockey.GetMethod(launchScreenService.ossClient, "UploadImg")).To(func(data []byte, remotePath string) error {
				if tc.name == "UploadImg error" {
					return assert.AnError
				}
				if tc.mockCloudReturn != nil {
					cloudErr, ok := tc.mockCloudReturn.(error)
					if !ok {
						return assert.AnError
					}
					return cloudErr
				}
				return nil
			}).Build()

			// 根据测试用例设置GenerateImgName的mock返回值
			mockey.Mock(mockey.GetMethod(launchScreenService.ossClient, "GenerateImgName")).To(func(suffix string) (string, string, error) {
				if tc.name == "GenerateImgName error" {
					return "", "", assert.AnError
				}
				return expectedResult.Url, expectedResult.Url, nil
			}).Build()

			// 根据测试用例设置GetImageFileType的mock返回值
			mockey.Mock(utils.GetImageFileType).To(func(fileBytes *[]byte) (string, error) {
				if tc.name == "GetImageFileType error" {
					return "", assert.AnError
				}
				return "jpg", nil
			}).Build()

			mockey.Mock((*launchScreenDB.DBLaunchScreen).UpdateImage).Return(tc.mockReturn, nil).Build()
			result, err := launchScreenService.UpdateImagePath(req)
			if tc.expectError {
				assert.Nil(t, result)
				switch {
				case !tc.mockIsExist:
					assert.EqualError(t, err, "LaunchScreenService.UpdateImagePath db.GetImageById error: record not found")
				case tc.name == "GetImageFileType error":
					assert.Error(t, err)
				case tc.name == "GenerateImgName error":
					assert.Error(t, err)
					assert.ErrorContains(t, err, "ossClient.GenerateImgName error")
				case tc.name == "UploadImg error":
					assert.Error(t, err)
					assert.ErrorContains(t, err, "LaunchScreenService.UpdateImagePath error")
				default:
					assert.EqualError(t, err, "LaunchScreenService.UpdateImagePath error: ["+strconv.Itoa(errno.BizFileUploadErrorCode)+"] "+errno.UpcloudError.ErrorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
