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
	"fmt"

	"github.com/west2-online/fzuhelper-server/cmd/user/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/pwd"
)

func (s *UserService) Register(req *user.RegisterRequest) (*db.User, error) {
	PwdDigest := pwd.SetPassword(req.Password)
	id, err := db.SF.NextVal()
	if err != nil {
		return nil, fmt.Errorf("User.Register SFCreateIDError:%v", err.Error())
	}
	userModel := &db.User{
		ID:       id,
		Number:   req.Number,
		Password: PwdDigest,
	}
	return db.Register(s.ctx, userModel)
}
