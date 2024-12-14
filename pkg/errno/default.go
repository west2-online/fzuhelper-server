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

// This file is be designed to define any common error so that we can use it in any service simply.

package errno

import "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	Success                   = NewErrNo(SuccessCode, "Success")
	CustomLaunchScreenSuccess = NewErrNo(consts.StatusOK, "Success") // 兼容处理

	AuthError          = NewErrNo(AuthErrorCode, "Auth Failed")                    // 鉴权失败，通常是内部错误，如解析失败
	AuthInvalid        = NewErrNo(AuthInvalidCode, "Auth Invalid")                 // 鉴权无效，如令牌颁发者不是 west2-online
	AuthAccessExpired  = NewErrNo(AuthAccessExpiredCode, "Access Token Expired")   // 访问令牌过期
	AuthRefreshExpired = NewErrNo(AuthRefreshExpiredCode, "Refresh Token Expired") // 刷新令牌过期
	AuthMissing        = NewErrNo(AuthMissingCode, "Auth Missing")                 // 鉴权缺失，如访问令牌缺失

	ParamError         = NewErrNo(ParamErrorCode, "parameter error") // 参数校验失败，可能是参数为空、参数类型错误等
	ParamMissingHeader = NewErrNo(ParamMissingHeaderCode, "missing request header data (id or cookies)")

	BizError             = NewErrNo(BizErrorCode, "business error")
	InternalServiceError = NewErrNo(InternalServiceErrorCode, "internal service error")

	SuffixError           = NewErrNo(ParamErrorCode, "invalid file")
	NoAccessError         = NewErrNo(AuthErrorCode, "user don't have authority to this biz")
	NoRunningPictureError = NewErrNo(BizErrorCode, "no valid picture")
	NoMatchingPlanError   = NewErrNo(BizErrorCode, "no matching plan")

	// internal error
	UpcloudError    = NewErrNo(BizFileUploadErrorCode, "upload to upcloud error")
	SFCreateIDError = NewErrNo(InternalDatabaseErrorCode, "sf create id failed")

	// redis
	RedisError = NewErrNo(InternalRedisErrorCode, "redis error")
)
