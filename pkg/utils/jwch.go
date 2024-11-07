package utils

import (
	"errors"
	"fmt"
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
func ParseJwchStuId(id string) (string, error) {
	if len(id) != StuIDLen {
		return "", errors.New("invalid id")
	}

	return id[len(id)-9:], nil
}
