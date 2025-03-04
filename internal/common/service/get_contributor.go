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

package service

import (
	"fmt"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func (s *CommonService) GetContributorInfo() (map[string][]*model.Contributor, error) {
	contributorKeys := []string{
		constants.ContributorFzuhelperAppKey,
		constants.ContributorFzuhelperServerKey,
		constants.ContributorJwchKey,
		constants.ContributorYJSYKey,
	}

	contributors := make(map[string][]*model.Contributor)

	// 遍历四个 key，依次从缓存中获取数据
	for _, key := range contributorKeys {
		if ok := s.cache.IsKeyExist(s.ctx, key); !ok {
			return nil, fmt.Errorf("service.GetContributorInfo: %s not exist", key)
		}

		// 获取当前 key 对应的 contributor 数据
		contributorInfo, err := s.cache.Common.GetContributorInfo(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetContributorInfo: failed to get contributor info for key %s: %w", key, err)
		}

		// 将数据存入返回结果 map 中
		contributors[key] = contributorInfo
	}

	return contributors, nil
}
