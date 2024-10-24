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

	"gorm.io/gorm"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
)

func TestAddPointTime(t *testing.T) {
	type testCase struct {
		name           string
		mockReturn     interface{}
		expectingError bool
	}
	testCases := []testCase{
		{
			name:       "AddPointTime",
			mockReturn: nil,
		},
		{
			name:           "dbError",
			mockReturn:     gorm.ErrRecordNotFound,
			expectingError: true,
		},
	}
	req := &launch_screen.AddImagePointTimeRequest{
		PictureId: 2024,
	}
	defer mockey.UnPatchAll() // 撤销所有mock操作，不会影响其他测试

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background())

			mockey.Mock(db.AddPointTime).Return(tc.mockReturn).Build()

			err := launchScreenService.AddPointTime(req.PictureId)

			if tc.expectingError {
				assert.EqualError(t, err, "LaunchScreenService.AddPointTime err: record not found")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
