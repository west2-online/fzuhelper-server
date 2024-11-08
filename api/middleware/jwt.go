package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	jwtgo "github.com/golang-jwt/jwt"
	"github.com/hertz-contrib/jwt"
	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"time"
)

var (
	JwtMiddleware   *jwt.HertzJWTMiddleware
	Identity        = "identity"
	RefreshTokenKey = []byte("refresh_secret_key")
	AccessTokenKey  = []byte("access_token_key")
	RefreshTokenTTL = time.Minute * 15   // Access Token 有效期15分钟
	AccessTokenTTL  = time.Hour * 24 * 7 // Refresh Token 有效期7天
)

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:                 "fzuhelper-server jwt",
		SigningAlgorithm:      "HS256",
		Key:                   AccessTokenKey,
		MaxRefresh:            RefreshTokenTTL,
		TokenLookup:           "header:Authorization",
		TokenHeadName:         "Bearer",
		IdentityKey:           Identity,
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
	// 在 PayloadFunc 中添加过期时间
	claims := jwt.MapClaims{
		Identity: data,
		"exp":    time.Now().Add(AccessTokenTTL).Unix(),
	}
	return claims
}

// 用于设置登录时认证用户信息的函数
func authenticator(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	type VerifyInfoRequest struct {
	}
	var err error
	var req VerifyInfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("VeryfyInfo: BindAndValidate err: %v", err)
		return nil, err
	}
	loginData, err := api.GetLoginData(ctx)
	if err != nil {
		logger.Errorf("VeryfyInfo: GetLoginData err: %v", err)
		return nil, err
	}
	return loginData, nil
}

// 用于设置登录的响应函数
func loginResponse(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
	refreshToken, refreshExpire, err := generateRefreshToken(ctx, c)
	if err != nil {
		logger.Errorf("Generate refresh token failed: %v", err)
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, map[string]interface{}{
		"code":           code,
		"access_token":   token,
		"access_expire":  expire,
		"refresh_token":  refreshToken,
		"refresh_expire": refreshExpire,
	})
}

// 从token提取用户信息的函数
func identityHandler(ctx context.Context, c *app.RequestContext) interface{} {
	claims := jwt.ExtractClaims(ctx, c)
	return claims[Identity]
}

// 用于设置 jwt 验证流程失败的响应函数
func unauthorizedHandler(ctx context.Context, c *app.RequestContext, code int, message string) {
	pack.RespData(c, map[string]interface{}{
		"code":    code,
		"message": message,
	})
}

// 生成refresh_token
func generateRefreshToken(ctx context.Context, c *app.RequestContext) (string, time.Time, error) {
	refreshExpire := time.Now().Add(RefreshTokenTTL)
	claims := jwtgo.MapClaims{
		"exp":    refreshExpire.Unix(),
		Identity: c.Get(Identity),
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	RefreshToken, err := token.SignedString(RefreshTokenKey)
	return RefreshToken, refreshExpire, err
}

// 使用refresh_token刷新access_token
func refreshTokenHandler(ctx context.Context, c *app.RequestContext) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespError(c, err)
		return
	}
	token, err := jwtgo.Parse(req.RefreshToken, func(token *jwtgo.Token) (interface{}, error) {
		return RefreshTokenKey, nil
	})
	if err != nil || !token.Valid {
		pack.RespError(c, err)
		return
	}
	identity := token.Claims.(jwtgo.MapClaims)[Identity]
	newToken, expire, err := JwtMiddleware.TokenGenerator(identity)
	if err != nil || token.Valid {
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, map[string]interface{}{
		"access_token":  newToken,
		"access_expire": expire,
	})
}

func httpStatusMessageFunc(e error, ctx context.Context, c *app.RequestContext) string {
	return e.Error()
}
