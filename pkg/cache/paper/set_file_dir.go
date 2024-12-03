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

package paper

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

const TwoDay = 2 * constants.OneDay

func (c *CachePaper) SetFileDirCache(ctx context.Context, key string, dir model.UpYunFileDir) error {
	data, err := sonic.Marshal(dir)
	if err != nil {
		return errno.Errorf(errno.InternalJSONErrorCode, "dal.SetFileDirCache: Unmarshal dir info failed: %v", err)
	}
	if err = c.client.Set(ctx, key, data, TwoDay).Err(); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "%v", err)
	}
	return nil
}
