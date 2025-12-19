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
	"math/rand"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// GenerateRandomCode 生成指定位数的随机验证码(字母+数字）
func GenerateRandomCode(length int) string {
	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = constants.CharSet[r.Intn(len(constants.CharSet))]
	}

	return string(code)
}
