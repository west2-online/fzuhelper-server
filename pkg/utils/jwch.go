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
	"strings"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/jwch"
)

const StuIDLen = 23

// RetryLogin 将会重新尝试进行登录
func RetryLogin(stu *jwch.Student) error {
	var err error
	delay := constants.InitialDelay

	for attempt := 1; attempt <= constants.MaxRetries; attempt++ {
		err = stu.Login()
		if err == nil {
			return nil // 登录成功
		}

		if attempt < constants.MaxRetries {
			time.Sleep(delay) // 等待一段时间后再重试
			delay *= 2        // 指数退避，逐渐增加等待时间
		}
	}

	return fmt.Errorf("failed to login after %d attempts: %w", constants.MaxRetries, err)
}

// ParseJwchStuId 用于解析教务处 id 里的学号
// 如 20241025133150102401339 的后 9 位
func ParseJwchStuId(id string) string {
	return id[len(id)-9:]
}

// GenerateCourseHash 生成课程的唯一哈希
func GenerateCourseHash(name, term, teacher, electiveType string) string {
	input := strings.Join([]string{name, term, teacher, electiveType}, "|")
	return SHA256(input)
}
