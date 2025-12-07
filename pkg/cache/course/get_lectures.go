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
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/west2-online/jwch"
)

func (c *CacheCourse) GetLecturesCache(ctx context.Context, key string) ([]*jwch.Lecture, error) {
	var lects []*jwch.Lecture
	buf, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("dal.GetLecturesCache: cache failed: %w", err)
	}
	if err = sonic.Unmarshal(buf, &lects); err != nil {
		return nil, fmt.Errorf("dal.GetLecturesCache: Unmarshal failed: %w", err)
	}
	return lects, nil
}
