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
	"time"
)

// StrOrEmpty 如果输入为空指针，返回空串；否则解引用
func StrOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(*p)
}

// I64OrZero 如果输入为空指针，返回0；否则解引用
func I64OrZero(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// TimePtrFromMillis 将一个指向毫秒时间戳的指针转换为 *time.Time。
// 如果入参指针为 nil 或时间戳 <= 0，则返回 nil。
// 返回的时间使用 UTC 时区。
func TimePtrFromMillis(p *int64) *time.Time {
	if p == nil || *p <= 0 {
		return nil
	}
	t := time.UnixMilli(*p).UTC()
	return &t
}
