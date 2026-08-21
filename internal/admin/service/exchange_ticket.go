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
	"github.com/west2-online/fzuhelper-server/api/mw"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// ExchangeTicket 原子消费一次性 ticket，并签发管理面板 access token
func (s *AdminService) ExchangeTicket(req *admin.ExchangeTicketRequest) (string, error) {
	adminId, exists, err := s.cache.Admin.ConsumeLoginTicketCache(s.ctx, req.Ticket)
	if err != nil {
		return "", errno.RedisError.WithError(err)
	}
	if !exists || adminId == "" {
		return "", errno.AuthError.WithMessage("invalid or expired login ticket")
	}

	accessToken, err := mw.CreateAdminToken(adminId)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
