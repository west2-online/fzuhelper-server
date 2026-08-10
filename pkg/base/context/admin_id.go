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

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const adminIDKey = "adminID"

func WithAdminID(ctx context.Context, adminID string) context.Context {
	value, err := sonic.MarshalString(adminID)
	if err != nil {
		logger.Infof("Failed to marshal adminID: %v", err)
	}
	return newContext(ctx, adminIDKey, value)
}

func GetAdminID(ctx context.Context) (string, error) {
	json, ok := fromContext(ctx, adminIDKey)
	if !ok {
		return "", errno.ParamMissingHeader.WithMessage("Failed to get header in context")
	}
	adminID := ""
	err := sonic.UnmarshalString(json, &adminID)
	if err != nil {
		return "", errno.InternalServiceError.WithMessage("Failed to get header in context when unmarshalling adminID")
	}
	return adminID, nil
}
