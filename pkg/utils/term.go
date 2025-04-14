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
	"fmt"
	"strconv"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// TransformSemester 将学期字符串转换为 "startYear-endYear" 和 "term" 的格式
// e.g: 202401 -> "2024-2025-1"
func TransformSemester(data string) (string, error) {
	// 检查输入是否合法
	if len(data) != constants.YJSTERMLEN {
		return "", fmt.Errorf("invalid input length: expected 6 characters, got %d", len(data))
	}

	// 提取年份和学期编号
	year := data[:4] // 前四位作为起始年份
	term := data[4:] // 后两位作为学期编号

	// 检查年份是否为有效数字
	startYear, err := strconv.Atoi(year)
	if err != nil {
		return "", fmt.Errorf("invalid year format: %s", year)
	}

	// 检查学期编号是否有效（"01" 或 "02"）
	termInt, err := strconv.Atoi(term) // 转换为整数以去掉前导 0
	if err != nil {
		return "", fmt.Errorf("invalid term: expected '01' or '02', got %s", term)
	}

	// 计算结束年份
	endYear := startYear + 1

	// 返回结果
	return fmt.Sprintf("%d-%d-%d", startYear, endYear, termInt), nil
}
