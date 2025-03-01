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

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (c *CacheAcademic) SetScoresCache(ctx context.Context, key string, scores []*jwch.Mark) error {
	data, err := sonic.Marshal(scores)
	if err != nil {
		logger.Errorf("dal.SetScoresCache: Marshal scores info failed: %v", err)
		return err
	}
	err = c.client.Set(ctx, key, data, constants.AcademicScoresExpire).Err()
	if err != nil {
		logger.Errorf("dal.SetScoresCache: Set scores info failed: %v", err)
		return err
	}
	return nil
}

func (c *CacheAcademic) SetScoresCacheYjsy(ctx context.Context, key string, scores []*yjsy.Mark) error {
	data, err := sonic.Marshal(scores)
	if err != nil {
		logger.Errorf("dal.SetScoresCacheYjsy: Marshal scores info failed: %v", err)
		return err
	}
	err = c.client.Set(ctx, key, data, constants.AcademicScoresExpire).Err()
	if err != nil {
		logger.Errorf("dal.SetScoresCacheYjsy: Set scores info failed: %v", err)
		return err
	}
	return nil
}
