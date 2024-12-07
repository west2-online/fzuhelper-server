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

package mw

import (
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func TestCreateAllToken(t *testing.T) {
	type testCase struct {
		name             string
		mockAccessToken  string
		mockRefreshToken string
		mockError        error
		expectingError   bool
	}

	testCases := []testCase{
		{
			name:             "CreateAllTokenSuccess",
			mockAccessToken:  "access_token_example",
			mockRefreshToken: "refresh_token_example",
			mockError:        nil,
			expectingError:   false,
		},
		{
			name:             "CreateAccessTokenError",
			mockAccessToken:  "",
			mockRefreshToken: "",
			mockError:        fmt.Errorf("failed to create access token"),
			expectingError:   true,
		},
		{
			name:             "CreateRefreshTokenError",
			mockAccessToken:  "access_token_example",
			mockRefreshToken: "",
			mockError:        fmt.Errorf("failed to create refresh token"),
			expectingError:   true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(CreateToken).When(func(tokenType int64) bool { return tokenType == constants.TypeAccessToken }).
				Return(tc.mockAccessToken, tc.mockError).
				When(func(tokenType int64) bool { return tokenType == constants.TypeRefreshToken }).
				Return(tc.mockRefreshToken, tc.mockError).
				Build()

			accessToken, refreshToken, err := CreateAllToken()

			if tc.expectingError {
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockAccessToken, accessToken)
				assert.Equal(t, tc.mockRefreshToken, refreshToken)
			}
		})
	}
}
