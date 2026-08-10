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
