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

	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

func TestGetLoginDataRPC(t *testing.T) {
	type testCase struct {
		name            string
		mockRespId      string
		mockRespCookies []string
		mockError       error
		expectedId      string
		expectedCookies []string
		expectingError  bool
	}

	testCases := []testCase{
		{
			name:            "GetLoginDataSuccess",
			mockRespId:      "123456",
			mockRespCookies: []string{"session_token=abc123", "path=/"},
			mockError:       nil,
			expectedId:      "123456",
			expectedCookies: []string{"session_token=abc123", "path=/"},
			expectingError:  false,
		},
		{
			name:            "GetLoginDataRPCError",
			mockRespId:      "",
			mockRespCookies: nil,
			mockError:       fmt.Errorf("RPC call failed"),
			expectedId:      "",
			expectedCookies: nil,
			expectingError:  true,
		},
	}

	req := &user.GetLoginDataRequest{
		Id:       "test_user",
		Password: "test_password",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(GetLoginDataRPC).Return(tc.mockRespId, tc.mockRespCookies, tc.mockError).Build()
			id, cookies, err := GetLoginDataRPC(context.Background(), req)
			if tc.expectingError {
				assert.Empty(t, id)
				assert.Nil(t, cookies)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedId, id)
				assert.Equal(t, tc.expectedCookies, cookies)
			}
		})
	}
}
