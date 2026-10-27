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

package mcp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	cachemcp "github.com/west2-online/fzuhelper-server/pkg/cache/mcp"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

var mcpCache *cachemcp.CacheMCP

func WithAuth(
	ctx context.Context,
	request mcp.CallToolRequest,
	handler func(context.Context, *AuthParams) (*mcp.CallToolResult, error),
) (*mcp.CallToolResult, error) {
	creds := CredentialsFromHeader(request.Header)

	auth, err := getAuth(ctx, creds, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := handler(WithLoginData(ctx, auth), auth)
	if err == nil {
		return result, nil
	}
	if !isCookieExpired(err) {
		return mcp.NewToolResultError(err.Error()), nil
	}

	auth, err = getAuth(ctx, creds, auth)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("login session expired and re-login failed: %v", err)), nil
	}

	result, err = handler(WithLoginData(ctx, auth), auth)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return result, nil
}

func authCacheKey(creds *Credentials) string {
	mac := hmac.New(sha256.New, []byte(config.Server.Secret))
	mac.Write([]byte(creds.StuID))
	mac.Write([]byte{0})
	mac.Write([]byte(creds.Password))
	mac.Write([]byte{0})
	mac.Write([]byte(creds.StudentType))
	return hex.EncodeToString(mac.Sum(nil))
}

func getAuth(ctx context.Context, creds *Credentials, stale *AuthParams) (*AuthParams, error) {
	key := authCacheKey(creds)
	cached, err := mcpCache.GetAuth(ctx, key)
	if err != nil {
		return nil, err
	}
	if cached != nil && (stale == nil || cached.Cookies != stale.UserCookies) {
		return &AuthParams{
			UserID:      cached.UserId,
			UserCookies: cached.Cookies,
		}, nil
	}
	return login(ctx, creds, key)
}

func login(ctx context.Context, creds *Credentials, key string) (*AuthParams, error) {
	if creds.StuID == "" {
		return nil, fmt.Errorf("missing %s header, cannot log in or refresh the login session; "+
			"please configure Stu-Id/Stu-Password/Stu-Type headers in your MCP client", HeaderStudentID)
	}
	if creds.Password == "" {
		return nil, fmt.Errorf("missing %s header, cannot refresh the login session without the password; "+
			"please configure Stu-Id/Stu-Password/Stu-Type headers in your MCP client", HeaderStudentPassword)
	}

	var (
		id      string
		cookies string
		err     error
	)
	switch creds.StudentType {
	case studentTypePostgraduate:
		id, cookies, err = rpc.GetLoginDataForYJSYRPC(ctx, &user.GetLoginDataForYJSYRequest{
			Id:       creds.StuID,
			Password: creds.Password,
		})
	default:
		id, cookies, err = rpc.GetLoginDataRPC(ctx, &user.GetLoginDataRequest{
			Id:       creds.StuID,
			Password: creds.Password,
		})
	}
	if err != nil {
		return nil, err
	}

	if err = mcpCache.SetAuth(ctx, key, &model.MCPAuth{
		Account: creds.StuID,
		UserId:  id,
		Cookies: cookies,
	}); err != nil {
		return nil, err
	}

	return &AuthParams{
		UserID:      id,
		UserCookies: cookies,
	}, nil
}

func loginOnInitialize(ctx context.Context, _ any, request *mcp.InitializeRequest) {
	creds := CredentialsFromHeader(request.Header)
	if creds.StuID == "" && creds.Password == "" {
		return
	}
	if _, err := getAuth(ctx, creds, nil); err != nil {
		logger.WarnfWithCtx(ctx, "mcp: auto login on initialize failed: %v", err)
	}
}

func isCookieExpired(err error) bool {
	var errNo errno.ErrNo
	return errors.As(err, &errNo) && errNo.ErrorCode == errno.BizJwchCookieExceptionCode
}
