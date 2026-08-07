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
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bytedance/sonic"
	"golang.org/x/oauth2"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type feishuUserInfo struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

type feishuUserInfoResponse struct {
	Code int64           `json:"code"`
	Msg  string          `json:"msg"`
	Data *feishuUserInfo `json:"data"`
}

// Callback 由 oauth 回调，校验 state 并通过 code 拿用户信息
func (s *AdminService) Callback(req *admin.CallbackRequest) (string, error) {
	if config.Admin == nil || config.Admin.Feishu.AppID == "" || config.Admin.Feishu.AppSecret == "" ||
		config.Admin.Feishu.RedirectURI == "" {
		return "", errno.InternalServiceError.WithMessage("incomplete Feishu OAuth config")
	}

	// 原子消费（读取并删除）缓存中的 state，并拿到要返回的 url
	returnTo, exists, err := s.cache.Admin.ConsumeOAuthStateCache(s.ctx, req.State)
	if err != nil {
		return "", errno.RedisError.WithError(err)
	}
	if !exists {
		return "", errno.AuthError.WithMessage("invalid or expired OAuth state")
	}

	// 用 code 换取飞书 user_access_token
	token, err := s.getFeishuUserAccessToken(req.Code)
	if err != nil {
		return "", err
	}

	// 获取飞书用户并验证
	userInfo, err := s.getFeishuUserInfo(token)
	if err != nil {
		return "", err
	}
	allowed, adminUser, err := s.db.Admin.GetAdminByFeishuUserId(s.ctx, userInfo.UserId)
	if err != nil {
		return "", err
	}
	if !allowed {
		return buildCallbackRedirectUrl(returnTo, "error", "forbidden")
	}

	// 生成 ticket 并缓存
	ticket, err := createOAuthState()
	if err != nil {
		return "", errno.InternalServiceError.WithError(err)
	}
	if err = s.cache.Admin.SetLoginTicketCache(s.ctx, ticket, strconv.FormatInt(adminUser.Id, 10)); err != nil {
		return "", errno.RedisError.WithError(err)
	}

	return buildCallbackRedirectUrl(returnTo, "ticket", ticket)
}

// getFeishuUserAccessToken 用 code 换取飞书 user_access_token
func (s *AdminService) getFeishuUserAccessToken(code string) (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(s.ctx, constants.AdminOAuthRequestTimeout)
	defer cancel()
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil || token == nil || !token.Valid() {
		return nil, errno.AuthError.WithMessage("Feishu OAuth authorization failed")
	}
	return token, nil
}

// getFeishuUserInfo 获取飞书用户信息
func (s *AdminService) getFeishuUserInfo(token *oauth2.Token) (*feishuUserInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, constants.AdminOAuthRequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.userInfoUrl, nil)
	if err != nil {
		return nil, errno.InternalServiceError.WithMessage("create Feishu user info request failed")
	}
	resp, err := s.oauthConfig.Client(ctx, token).Do(req)
	if err != nil {
		return nil, errno.InternalServiceError.WithMessage("request Feishu user info failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errno.AuthError.WithMessage("get Feishu user info failed")
	}
	var userResp feishuUserInfoResponse
	if err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, errno.InternalServiceError.WithMessage("decode Feishu user info response failed")
	}
	if userResp.Code != 0 || userResp.Data == nil {
		return nil, errno.AuthError.WithMessage("get Feishu user info failed")
	}
	if userResp.Data.UserId == "" {
		return nil, errno.InternalServiceError.WithMessage("Feishu user_id is empty")
	}
	return userResp.Data, nil
}

func buildCallbackRedirectUrl(returnTo, key, value string) (string, error) {
	redirectUrl, err := url.Parse(returnTo)
	if err != nil {
		return "", fmt.Errorf("parse callback return url failed: %w", err)
	}
	query := redirectUrl.Query()
	query.Set(key, value)
	redirectUrl.RawQuery = query.Encode()
	return redirectUrl.String(), nil
}
