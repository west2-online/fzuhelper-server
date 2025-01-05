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

package login_data

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type Key string

var loginDataKey Key = "loginData"

func GetLoginData(ctx context.Context) (*model.LoginData, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errno.ParamMissingHeader.WithMessage("Failed to get header in context")
	}
	return user, nil
}

func NewContext(ctx context.Context, value *model.LoginData) context.Context {
	return context.WithValue(ctx, loginDataKey, value)
}

func FromContext(ctx context.Context) (*model.LoginData, bool) {
	u, ok := ctx.Value(loginDataKey).(*model.LoginData)
	return u, ok
}
