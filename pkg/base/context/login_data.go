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

package context

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const loginDataKey string = "loginData"

// WithLoginData 将LoginData加入到context中，通过metainfo传递到RPC server
func WithLoginData(ctx context.Context, loginData *model.LoginData) context.Context {
	value, err := sonic.MarshalString(*loginData)
	if err != nil {
		logger.Infof("Failed to marshal LoginData: %v", err)
	}
	return newContext(ctx, loginDataKey, value)
}

// GetLoginData 从context中取出LoginData
func GetLoginData(ctx context.Context) (*model.LoginData, error) {
	user, ok := fromContext(ctx, loginDataKey)
	if !ok {
		return nil, errno.ParamMissingHeader.WithMessage("Failed to get header in context")
	}
	value := new(model.LoginData)
	err := sonic.UnmarshalString(user, value)
	if err != nil {
		return nil, errno.InternalServiceError.WithMessage("Failed to get header in context when unmarshalling loginData")
	}
	return value, nil
}

// ExtractIDFromLoginData 从 LoginData 中提取学号
func ExtractIDFromLoginData(data *model.LoginData) string {
	if data == nil {
		return ""
	}
	return ExtractIDFromIdentifier(data.Id)
}

// ExtractIDFromIdentifier 从 Identifier 中提取学号
// 研究生：id直接是stuId，以 00000 为前缀（可能是 9 或 10 位）
// 本科生：identifier 末尾为学号，按入学年份判断是 9 位还是 10 位
// 2026年开始本科生学号由9位变到10位
// 本科生identifier规则：yyyymmdd + hhmmss + 学号，长度不定
// 比如2026年8月5日，开头就是202685，12月17日就是20261217
// 时间也是一样，24小时制，当日1时50分07秒就是1507，10时0分0秒就是1000，12时37分59秒就是123759
func ExtractIDFromIdentifier(id string) string {
	if len(id) < constants.StudentIDLength {
		return ""
	}

	// 研究生
	if utils.IsGraduate(id) {
		return utils.RemoveGraduatePrefix(id)
	}

	// 本科生
	return utils.RemoveUndergraduatePrefix(id)
}
