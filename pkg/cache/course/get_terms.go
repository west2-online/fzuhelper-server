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

func (c *CacheCourse) GetTermsCache(ctx context.Context, key string) (terms *jwch.Term, err error) {
	terms = new(jwch.Term)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("dal.GetTermsCache: cache failed: %w", err)
	}
	if err = sonic.Unmarshal(data, terms); err != nil {
		return nil, fmt.Errorf("dal.GetTermsCache: Unmarshal failed: %w", err)
	}
	return terms, nil
}
