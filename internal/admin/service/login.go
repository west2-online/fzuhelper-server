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
	"crypto/rand"
	"encoding/base64"
	"slices"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

const oauthStateBytes = 32

// Login 生成 state，并返回授权页面 url（目前为飞书）
func (s *AdminService) Login(req *admin.LoginRequest) (string, error) {
	if config.Admin == nil || config.Admin.Feishu.AppID == "" || config.Admin.Feishu.RedirectURI == "" {
		return "", errno.InternalServiceError.WithMessage("incomplete Feishu OAuth config")
	}
	if !slices.Contains(config.Admin.Feishu.AllowedReturnUrls, req.ReturnTo) {
		return "", errno.ParamError.WithMessage("returnTo is not allowed")
	}

	state, err := createOAuthState()
	if err != nil {
		return "", errno.InternalServiceError.WithError(err)
	}
	if err = s.cache.Admin.SetOAuthStateCache(s.ctx, state, req.ReturnTo); err != nil {
		return "", errno.RedisError.WithError(err)
	}

	return s.oauthConfig.AuthCodeURL(state), nil
}

func createOAuthState() (string, error) {
	randomBytes := make([]byte, oauthStateBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}
