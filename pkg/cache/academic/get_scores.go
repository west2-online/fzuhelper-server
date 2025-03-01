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

package academic

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (c *CacheAcademic) GetScoresCache(ctx context.Context, key string) (scores []*jwch.Mark, err error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errno.Errorf(errno.InternalJSONErrorCode, "dal.GetScoresCache: Get scores info failed: %v", err)
	}
	err = sonic.Unmarshal(data, &scores)
	if err != nil {
		return nil, errno.Errorf(errno.InternalJSONErrorCode, "dal.GetScoresCache: Unmarshal scores info failed: %v", err)
	}
	return scores, nil
}

func (c *CacheAcademic) GetScoresCacheYjsy(ctx context.Context, key string) (scores []*yjsy.Mark, err error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errno.Errorf(errno.InternalJSONErrorCode, "dal.GetScoresCacheYjsy: Get scores info failed: %v", err)
	}
	err = sonic.Unmarshal(data, &scores)
	if err != nil {
		return nil, errno.Errorf(errno.InternalJSONErrorCode, "dal.GetScoresCacheYjsy: Unmarshal scores info failed: %v", err)
	}
	return scores, nil
}
