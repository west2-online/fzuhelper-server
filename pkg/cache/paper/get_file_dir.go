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
	"github.com/pkg/errors"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func (c *CachePaper) GetFileDirCache(ctx context.Context, key string) (bool, *model.UpYunFileDir, error) {
	ret := &model.UpYunFileDir{}

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return false, ret, errors.Errorf("dal.GetFileDirCache: get dir info failed: %v", err)
	}
	err = sonic.Unmarshal(data, &ret)
	if err != nil {
		return false, ret, errors.Errorf("dal.GetFileDirCache: Unmarshal dir info failed: %v", err)
	}
	return true, ret, nil
}
