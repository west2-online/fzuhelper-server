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
	"net/http"

	"github.com/west2-online/fzuhelper-server/pkg/utils"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
)

type UserService struct {
	ctx        context.Context
	Identifier string
	cookies    []*http.Cookie
	db         *db.Database
	sf         *utils.Snowflake
	cache      *cache.Cache
}

func NewUserService(ctx context.Context, identifier string, cookies []*http.Cookie, clientset *base.ClientSet) *UserService {
	return &UserService{
		ctx:        ctx,
		Identifier: identifier,
		cookies:    cookies,
		db:         clientset.DBClient,
		cache:      clientset.CacheClient,
		sf:         clientset.SFClient,
	}
}
