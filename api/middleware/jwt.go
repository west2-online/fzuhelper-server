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

package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	jwtgo "github.com/golang-jwt/jwt"
	"github.com/hertz-contrib/jwt"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

var JwtMiddleware *jwt.HertzJWTMiddleware

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:                 "fzuhelper-server jwt",
		SigningAlgorithm:      "HS256",
		Key:                   []byte(config.JwtKeys.AccessTokenKey),
		MaxRefresh:            constants.RefreshTokenTTL,
		Timeout:               constants.AccessTokenTTL,
		TokenLookup:           "header:Authorization",
		TokenHeadName:         "Bearer",
		IdentityKey:           constants.Identity,
		LoginResponse:         loginResponse,
		Authenticator:         authenticator,
		PayloadFunc:           payloadFunc,
		IdentityHandler:       identityHandler,
		Unauthorized:          unauthorizedHandler,
		HTTPStatusMessageFunc: httpStatusMessageFunc,
	})
	if err != nil {
		panic(err)
	}
}

// 用于设置登陆成功后为向 token 中添加自定义负载信息的函数
func payloadFunc(data interface{}) jwt.MapClaims {
	claims := jwt.MapClaims{
		"id": data,
	}
	return claims
}

// 用于设置登录时认证用户信息的函数
func authenticator(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	var err error
	var req api.GetAccessTokenRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("middleware.jwt: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return nil, err
	}

	identifier := c.Request.Header.Get("id")
	id := identifier[len(identifier)-9:]
	cookies := c.Request.Header.GetAll("cookies")

	err = jwch.NewStudent().
		WithUser(id, "").
		WithLoginData(identifier, utils.ParseCookies(cookies)).
		CheckSession()
	if err != nil {
		return nil, err
	}

	return id, nil
}

// 用于设置登录的响应函数
func loginResponse(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
	refreshToken, _, err := generateRefreshToken(ctx, c)
	if err != nil {
		logger.Errorf("Generate refresh token failed: %v", err)
		pack.RespError(c, err)
		return
	}
	c.Header("access_token", "Bearer "+token) // 加上 Bearer 前缀
	c.Header("refresh_token", refreshToken)
	pack.RespSuccess(c)
}

// 从token提取用户信息的函数
func identityHandler(ctx context.Context, c *app.RequestContext) interface{} {
	claims := jwt.ExtractClaims(ctx, c)
	return claims["id"]
}

// 用于设置 jwt 验证流程失败的响应函数
func unauthorizedHandler(ctx context.Context, c *app.RequestContext, code int, message string) {
	pack.RespError(c, errors.New(message))
}

// 生成refresh_token
func generateRefreshToken(ctx context.Context, c *app.RequestContext) (string, time.Time, error) {
	refreshExpire := time.Now().Add(constants.RefreshTokenTTL)
	id := c.Request.Header.Get("id")

	claims := jwtgo.MapClaims{
		"exp": refreshExpire.Unix(),
		"id":  id,
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	RefreshToken, err := token.SignedString([]byte(config.JwtKeys.RefreshTokenKey))
	return RefreshToken, refreshExpire, err
}

// 使用refresh_token刷新access_token
func RefreshTokenHandler(ctx context.Context, c *app.RequestContext) {
	refreshToken := c.Request.Header.Get("refresh_token")
	if refreshToken == "" {
		pack.RespError(c, errors.New("refresh_token is required"))
		return
	}
	// 解析 refresh_token
	token, err := jwtgo.Parse(refreshToken, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(config.JwtKeys.RefreshTokenKey), nil
	})
	if err != nil || !token.Valid {
		pack.RespError(c, err)
		return
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		pack.RespError(c, errors.New("invalid claims in refresh token"))
		return
	}

	identity, exists := claims["id"]

	if !exists {
		pack.RespError(c, errors.New("identity not found in refresh token"))
		return
	}

	// 使用 identity 生成新的 access_token
	newToken, _, err := JwtMiddleware.TokenGenerator(identity)
	if err != nil {
		pack.RespError(c, err)
		return
	}

	c.Header("access_token", newToken)
	pack.RespSuccess(c)
}

func httpStatusMessageFunc(e error, ctx context.Context, c *app.RequestContext) string {
	return e.Error()
}
