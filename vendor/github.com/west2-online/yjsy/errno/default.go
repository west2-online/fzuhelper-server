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

package errno

var (
	// Success
	Success = NewErrNo(SuccessCode, "Success")

	ServiceError             = NewErrNo(ServiceErrorCode, "service is unable to start successfully")
	ServiceInternalError     = NewErrNo(ServiceErrorCode, "service Internal Error")
	ParamError               = NewErrNo(ParamErrorCode, "parameter error")
	AuthorizationFailedError = NewErrNo(AuthorizationFailedErrCode, "authorization failed")

	// User
	AccountConflictError  = NewErrNo(AuthorizationFailedErrCode, "account conflict")
	CookieError           = NewErrNo(AuthorizationFailedErrCode, "session expired")
	SystemError           = NewErrNo(AuthorizationFailedErrCode, "system error")
	LoginCheckFailedError = NewErrNo(AuthorizationFailedErrCode, "login check failed")

	// HTTP
	HTTPQueryError = NewErrNo(HTTPQueryErrorCode, "HTTP query failed")
	HTMLParseError = NewErrNo(HTTPQueryErrorCode, "HTML parse failed")
)
