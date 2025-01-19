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

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/jwch"
)

func TestGetExamRoomInfo(t *testing.T) {
	type testCase struct {
		name           string
		mockReturn     interface{}
		expectedResult interface{}
		expectingError bool
	}

	tests := []testCase{
		{
			name: "GetExamRoomInfo",
			mockReturn: []*jwch.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectedResult: []*jwch.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectingError: true,
		},
	}

	req := &classroom.ExamRoomInfoRequest{
		Term: "202401",
	}

	defer mockey.UnPatchAll()

	// 运行所有测试用例
	for _, tc := range tests {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockey.Mock((*jwch.Student).WithLoginData).Return(jwch.NewStudent()).Build()
			mockey.Mock((*jwch.Student).GetExamRoom).Return(tc.mockReturn, nil).Build()
			// mock login data
			loginData := new(model.LoginData)
			ctx := customContext.WithLoginData(context.Background(), loginData)

			classroomService := NewClassroomService(ctx, mockClientSet)
			result, err := classroomService.GetExamRoomInfo(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
