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

// MVC--Model
package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/cmd/user/pack/pwd"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type User struct {
	ID        int64
	Account   string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

func Register(ctx context.Context, userModel *User) (*User, error) {
	userResp := new(User)
	// WithContext(ctx)是将一个context.Context对象和数据库连接绑定，以实现在数据库操作中使用context.Context上下文传递。
	if err := DB.WithContext(ctx).Where("account = ? OR name = ?", userModel.Account, userModel.Name).First(&userResp).Error; err == nil {
		return nil, errno.UserExistedError
	}

	if err := DB.WithContext(ctx).Create(userModel).Error; err != nil {
		return nil, err
	}
	return userModel, nil
}

func Login(ctx context.Context, userModel *User) (*User, error) {
	userResp := new(User)
	if err := DB.WithContext(ctx).Where("account = ?", userModel.Account).
		First(&userResp).Error; err != nil {
		return nil, errno.UserNonExistError
	}

	if !pwd.CheckPassword(userResp.Password, userModel.Password) {
		return nil, errno.AuthFailedError
	}

	return userResp, nil
}
