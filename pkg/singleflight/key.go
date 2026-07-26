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

package singleflight

import (
	"fmt"
	"strings"
)

// 思考了一下，决定把 prefix 放constant里，动态key合并为一个函数

func Key(prefix string, parts ...any) string {
	if len(parts) == 0 {
		return prefix
	}

	keys := make([]string, 0, len(parts)+1)
	keys = append(keys, prefix)
	for _, part := range parts {
		keys = append(keys, fmt.Sprint(part))
	}
	return strings.Join(keys, ":")
}
