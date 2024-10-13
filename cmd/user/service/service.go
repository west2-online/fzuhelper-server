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

	"github.com/west2-online/fzuhelper-server/cmd/user/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

type UserService struct {
	ctx        context.Context
	Identifier string
	cookies    []*http.Cookie
}

func NewUserService(ctx context.Context, identifier string, cookies []*http.Cookie) *UserService {
	return &UserService{
		ctx:        ctx,
		Identifier: identifier,
		cookies:    cookies,
	}
}

func BuildUserResp(dbUser *db.User) *model.User {
	return &model.User{
		Id:      dbUser.ID,
		Name:    dbUser.Name,
		Account: dbUser.Account,
	}
}
