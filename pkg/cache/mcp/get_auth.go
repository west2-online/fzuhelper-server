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
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func (c *CacheMCP) GetAuth(ctx context.Context, key string) (*model.MCPAuth, error) {
	data, err := c.client.GetEx(ctx, mcpAuthKey(key), constants.MCPAuthKeyExpire).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dal.GetAuth: get cache failed: %w", err)
	}

	auth := new(model.MCPAuth)
	if err = sonic.Unmarshal(data, auth); err != nil {
		return nil, fmt.Errorf("dal.GetAuth: unmarshal failed: %w", err)
	}
	return auth, nil
}
