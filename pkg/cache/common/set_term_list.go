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
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/jwch"
)

func (c *CacheCommon) SetTermListCache(ctx context.Context, key string, list *jwch.SchoolCalendar) error {
	termListJson, err := sonic.Marshal(list)
	if err != nil {
		return fmt.Errorf("dal.SetTermListCache: Marshal info failed: %w", err)
	}
	if err = c.client.Set(ctx, key, termListJson, constants.CommonTermListKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetTermListCache: Set cache failed: %w", err)
	}
	return nil
}
