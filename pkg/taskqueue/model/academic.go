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

package model

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/jwch"
)

// SetScoresCacheTask 定义
type SetScoresCacheTask struct {
	key     string
	scores  []*jwch.Mark
	cache   *cache.Cache
	context context.Context
}

func NewSetScoresCacheTask(key string, scores []*jwch.Mark, cache *cache.Cache, context context.Context) *SetScoresCacheTask {
	return &SetScoresCacheTask{
		key:     key,
		scores:  scores,
		cache:   cache,
		context: context,
	}
}

func (t *SetScoresCacheTask) Execute() error {
	return t.cache.Academic.SetScoresCache(t.context, t.key, t.scores)
}
