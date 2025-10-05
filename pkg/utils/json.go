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
	"encoding/json"

	"github.com/bytedance/sonic"
)

func JSONEncode(v interface{}) (string, error) {
	data, err := sonic.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// EnsureJSONArray 确保传入字符串是“合法的 JSON”且用于表示“数组”场景
func EnsureJSONArray(s string) string {
	if s == "" {
		return "[]"
	}
	if json.Valid([]byte(s)) {
		return s
	}
	return "[]"
}

// EnsureJSONObject 确保传入字符串是“合法的 JSON”且用于表示“对象”场景
func EnsureJSONObject(s string) string {
	if s == "" {
		return "{}"
	}
	if json.Valid([]byte(s)) {
		return s
	}
	return "{}"
}

// EnsureJSON 确保传入字符串是“合法的 JSON”（不限定必须为数组或对象）
func EnsureJSON(s string) string {
	if s == "" {
		return "[]"
	}
	if json.Valid([]byte(s)) {
		return s
	}
	return "[]"
}
