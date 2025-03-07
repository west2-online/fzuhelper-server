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

	"github.com/bytedance/sonic"
)

func (s *VersionService) GetDump() (string, error) {
	vs, err := s.db.Version.GetVersionList(s.ctx)
	if err != nil {
		return "", fmt.Errorf("GetDump: get version list error: %w", err)
	}
	result := make(map[string]int64)
	for _, v := range vs {
		result[v.Date] = v.Visits
	}
	jsonData, err := sonic.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("GetDump: marshal error: %w", err)
	}
	return string(jsonData), nil
}
