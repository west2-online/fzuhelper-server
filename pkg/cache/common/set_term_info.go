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

package common

import (
	"context"
	"log"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/jwch"
)

func (c *CacheCommon) SetTermInfo(ctx context.Context, key string, value *jwch.CalTermEvents) error {
	data, err := sonic.Marshal(value)
	if err != nil {
		return errno.Errorf(errno.InternalJSONErrorCode, "dal.SetFileDirCache: Unmarshal dir info failed: %v", err)
	}

	if err = c.client.Set(ctx, key, data, constants.TermInfoKeyExpire).Err(); err != nil {
		log.Println(err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "%v", err)
	}
	return nil
}
