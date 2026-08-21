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

	"golang.org/x/oauth2"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
)

type AdminService struct {
	ctx         context.Context
	cache       *cache.Cache
	db          *db.Database
	oauthConfig *oauth2.Config
	userInfoUrl string
}

func NewAdminService(ctx context.Context, clientSet *base.ClientSet) *AdminService {
	return &AdminService{
		ctx:   ctx,
		cache: clientSet.CacheClient,
		db:    clientSet.DBClient,
		// 飞书文档 https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/authentication-management/access-token/get-user-access-token-v3
		oauthConfig: &oauth2.Config{
			ClientID:     config.Admin.Feishu.AppID,
			ClientSecret: config.Admin.Feishu.AppSecret,
			RedirectURL:  config.Admin.Feishu.RedirectURI,
			Scopes: []string{
				constants.FeishuAdminOAuthScope,
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:   constants.FeishuOAuthAuthorizeUrl,
				TokenURL:  constants.FeishuOAuthTokenUrl,
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
		userInfoUrl: constants.FeishuUserInfoUrl,
	}
}
