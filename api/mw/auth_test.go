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
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	serverConstants "github.com/west2-online/fzuhelper-server/pkg/constants"
)

func TestAuthBindsTokenToHeaderStudentID(t *testing.T) {
	tests := []struct {
		name          string
		tokenType     int64
		tokenStuID    string
		headerID      string
		expectHandler bool
	}{
		{
			name:          "matching student id",
			tokenType:     serverConstants.TypeAccessToken,
			tokenStuID:    "102300217",
			headerID:      "20268514814102300217",
			expectHandler: true,
		},
		{
			name:       "different student id",
			tokenType:  serverConstants.TypeAccessToken,
			tokenStuID: "102300217",
			headerID:   "20268514814102300218",
		},
		{
			name:       "anonymous legacy token",
			tokenType:  serverConstants.TypeAccessToken,
			tokenStuID: "",
			headerID:   "20268514814102300217",
		},
		{
			name:       "refresh token on protected endpoint",
			tokenType:  serverConstants.TypeRefreshToken,
			tokenStuID: "102300217",
			headerID:   "20268514814102300217",
		},
	}

	defer mockey.UnPatchAll()
	for _, tt := range tests {
		mockey.PatchConvey(tt.name, t, func() {
			mockey.Mock(CheckToken).Return(tt.tokenType, tt.tokenStuID, nil).Build()
			mockey.Mock(CreateAllToken).To(func(stuID string) (string, string, error) {
				assert.Equal(t, tt.tokenStuID, stuID)
				return "new-access-token", "new-refresh-token", nil
			}).Build()

			called := false
			router := route.NewEngine(&config.Options{})
			router.GET("/protected", Auth(), GetHeaderParams(), func(_ context.Context, c *app.RequestContext) {
				called = true
				c.String(consts.StatusOK, "passed")
			})

			res := ut.PerformRequest(router, consts.MethodGet, "/protected", nil,
				ut.Header{Key: serverConstants.AuthHeader, Value: "token"},
				ut.Header{Key: "Id", Value: tt.headerID},
				ut.Header{Key: "Cookies", Value: "session=value"})

			assert.Equal(t, tt.expectHandler, called)
			if tt.expectHandler {
				assert.Equal(t, "passed", string(res.Result().Body()))
				assert.Equal(t, "new-access-token", string(res.Result().Header.Peek(serverConstants.AccessTokenHeader)))
			}
		})
	}
}
