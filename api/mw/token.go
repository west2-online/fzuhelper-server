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

package mw

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type Claims struct {
	StudentID string `json:"student_id"`
	Type      int64  `json:"type"`
	jwt.StandardClaims
}

// CreateAllToken 创建一对 token，第一个是 access token，第二个是 refresh token
func CreateAllToken() (string, string, error) {
	accessToken, err := CreateToken(constants.TypeAccessToken, "")
	if err != nil {
		return "", "", err
	}
	refreshToken, err := CreateToken(constants.TypeRefreshToken, "")
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// CreateToken 会通过不同 Token 类型创建不同的 Token
func CreateToken(tokenType int64, stuID string) (string, error) {
	if config.Server == nil {
		return "", errno.AuthError.WithMessage("server config not found")
	}

	var expireTime time.Time
	nowTime := time.Now()
	var token string
	var err error

	switch tokenType {
	case constants.TypeAccessToken:
		expireTime = nowTime.Add(constants.AccessTokenTTL)
	case constants.TypeRefreshToken:
		expireTime = nowTime.Add(constants.RefreshTokenTTL)
	case constants.TypeCalendarToken:
		expireTime = nowTime.Add(constants.CalendarTokenTTL)
	}

	claims := Claims{
		StudentID: stuID,
		Type:      tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间戳
			IssuedAt:  nowTime.Unix(),    // 当前时间戳
			Issuer:    constants.Issuer,  // 颁发者签名
		},
	}

	// 选择 Ed25519 是出于兼顾性能和安全性的考虑，PS512 安全性太高但性能不好，ES512 速度没有 Ed25519 快
	// 这里不考虑旧版的对称加密
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	key, err := jwt.ParseEdPrivateKeyFromPEM([]byte(config.Server.Secret))
	if err != nil {
		return "", errno.AuthError.WithMessage(fmt.Sprintf("parse private key failed, err: %v", err))
	}

	token, err = tokenStruct.SignedString(key)
	if err != nil {
		return "", errno.AuthError.WithMessage(fmt.Sprintf("sign token failed, err: %v", err))
	}
	return token, nil
}

// CheckToken 会检查 token 是否有效，如果有效则返回 token 类型，否则返回错误(type 会返回 -1)
// Check 成功后返回 token 中的 stu_id
func CheckToken(token string) (int64, string, error) {
	if token == "" {
		return -1, "", errno.AuthMissing
	}
	// 解析 token，但不进行签名验证
	tokenStruct, _, err := new(jwt.Parser).ParseUnverified(token, &Claims{})
	if err != nil {
		return -1, "", errno.AuthInvalid.WithError(err)
	}

	unverifiedClaims, ok := tokenStruct.Claims.(*Claims)
	if !ok {
		return -1, "", errno.AuthError.WithMessage("cannot handle claims")
	}

	secret, err := jwt.ParseEdPublicKeyFromPEM([]byte(constants.PublicKey))
	if err != nil {
		return -1, "", errno.AuthError.WithMessage(fmt.Sprintf("parse public key failed, err: %v", err))
	}

	// 使用正确的密钥再次解析 token
	response, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errno.AuthError.WithMessage(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return secret, nil
	})
	// 验证 token 是否有效
	if err != nil {
		return unverifiedClaims.Type, "", checkError(err, unverifiedClaims.Type)
	}

	if _, ok := response.Claims.(*Claims); ok && response.Valid {
		return unverifiedClaims.Type, unverifiedClaims.StudentID, nil
	}

	return -1, "", errno.AuthInvalid
}

// checkError 会检查错误类型并返回对应的错误(含过期)
func checkError(err error, tokenType int64) error {
	var ve *jwt.ValidationError
	if errors.As(err, &ve) {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			if tokenType == constants.TypeAccessToken {
				return errno.AuthAccessExpired
			}
			return errno.AuthRefreshExpired
		}
	}
	return errno.AuthError.WithMessage(err.Error())
}
