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

var (
	// Success
	Success = NewErrNo(SuccessCode, "Success")

	ParamError         = NewErrNo(ParamErrorCode, "parameter error")
	ParamEmpty         = NewErrNo(ParamEmptyCode, "some params that required are empty")
	ParamMissingHeader = NewErrNo(ParamMissingHeaderCode, "missing request header data (id or cookies)")

	AuthFailedError      = NewErrNo(AuthErrorCode, "authorization failed")
	BizError             = NewErrNo(BizErrorCode, "business error")
	InternalServiceError = NewErrNo(InternalServiceErrorCode, "internal service error")

	UserExistedError  = NewErrNo(InternalDatabaseErrorCode, "user existed")
	UserNonExistError = NewErrNo(InternalDatabaseErrorCode, "user didn't exist")
	SuffixError       = NewErrNo(ParamErrorCode, "invalid file")
	UpcloudError      = NewErrNo(BizFileUploadErrorCode, "upload to upcloud error")

	// redis
	RedisError = NewErrNo(InternalRedisErrorCode, "redis error")
)
