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

package course

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *CacheCourse) GetDateCache(ctx context.Context, key string) (date *model.LocateDate, err error) {
	res, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetDateCache: get cache failed: %v", err)
	}
	if err = sonic.Unmarshal(res, &date); err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetDateCache: Unmarshal failed: %v", err)
	}
	return date, nil
}
