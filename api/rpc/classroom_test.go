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

package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestGetEmptyRoomRPC(t *testing.T) {
	type testCase struct {
		name           string
		mockResp       []*model.Classroom
		mockError      error
		expectedResult []*model.Classroom
		expectingError bool
	}

	testCases := []testCase{
		{
			name: "GetEmptyRoomSuccess",
			mockResp: []*model.Classroom{
				{
					Build:    "东 1",
					Location: "旗山东1-201",
					Capacity: "60",
					Type:     "机房",
				},
			},
			mockError: nil,
			expectedResult: []*model.Classroom{
				{
					Build:    "东 1",
					Location: "旗山东1-201",
					Capacity: "60",
					Type:     "机房",
				},
			},
		},
		{
			name:           "GetEmptyRoomRPCError",
			mockResp:       nil,
			mockError:      fmt.Errorf("RPC call failed"),
			expectedResult: nil,
			expectingError: true,
		},
	}

	req := &classroom.EmptyRoomRequest{
		Date:      "2024-12-01",
		Campus:    "旗山校区",
		StartTime: "1",
		EndTime:   "1",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(GetEmptyRoomRPC).Return(tc.mockResp, tc.mockError).Build()
			result, err := GetEmptyRoomRPC(context.Background(), req)
			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetExamRoomInfoRPC(t *testing.T) {
	type testCase struct {
		name           string
		mockResp       []*model.ExamRoomInfo
		mockError      error
		expectedResult []*model.ExamRoomInfo
		expectingError bool
	}

	testCases := []testCase{
		{
			name: "GetExamRoomInfoSuccess",
			mockResp: []*model.ExamRoomInfo{
				{
					Name:     "高等数学",
					Credit:   "4.5",
					Teacher:  "张三",
					Location: "旗山校区教一201",
					Time:     "08:00-10:00",
					Date:     "2024-12-15",
				},
			},
			mockError: nil,
			expectedResult: []*model.ExamRoomInfo{
				{
					Name:     "高等数学",
					Credit:   "4.5",
					Teacher:  "张三",
					Location: "旗山校区教一201",
					Time:     "08:00-10:00",
					Date:     "2024-12-15",
				},
			},
		},
		{
			name:           "GetExamRoomInfoRPCError",
			mockResp:       nil,
			mockError:      fmt.Errorf("RPC call failed"),
			expectedResult: nil,
			expectingError: true,
		},
	}

	req := &classroom.ExamRoomInfoRequest{}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(GetExamRoomInfoRPC).Return(tc.mockResp, tc.mockError).Build()
			result, err := GetExamRoomInfoRPC(context.Background(), req)
			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
