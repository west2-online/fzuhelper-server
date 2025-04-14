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

package utils

import (
	"strings"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// IsGraduate 根据 id 判断是否是研究生
func IsGraduate(id string) (bool, error) {
	if len(id) == 0 {
		return false, errno.NewErrNo(errno.ParamErrorCode, "学号 id 不符合规范")
	}
	return strings.HasPrefix(id[:5], "00000"), nil
}
