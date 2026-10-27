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

package mcp

import (
	"context"
	"net/http"
	"strings"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
)

const (
	HeaderStudentID       = "Stu-Id"
	HeaderStudentPassword = "Stu-Password"
	HeaderStudentType     = "Stu-Type"

	studentTypeUndergraduate = "1"
	studentTypePostgraduate  = "2"
)

type Credentials struct {
	StuID       string
	Password    string
	StudentType string
}

func CredentialsFromHeader(header http.Header) *Credentials {
	return &Credentials{
		StuID:       strings.TrimSpace(header.Get(HeaderStudentID)),
		Password:    header.Get(HeaderStudentPassword),
		StudentType: normalizeStudentType(header.Get(HeaderStudentType)),
	}
}

func normalizeStudentType(studentType string) string {
	if strings.TrimSpace(studentType) == studentTypePostgraduate {
		return studentTypePostgraduate
	}
	return studentTypeUndergraduate
}

// AuthParams 认证参数
type AuthParams struct {
	UserID      string
	UserCookies string
}

// WithLoginData 将认证参数添加到 context 中
// 用于统一处理 LoginData 的注入
func WithLoginData(ctx context.Context, params *AuthParams) context.Context {
	return metainfoContext.WithLoginData(ctx, &model.LoginData{
		Id:      params.UserID,
		Cookies: params.UserCookies,
	})
}
