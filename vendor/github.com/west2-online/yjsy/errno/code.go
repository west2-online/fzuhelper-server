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

const (
	// Success
	SuccessCode = 10000
	SuccessMsg  = "success"

	// Error
	ServiceErrorCode           = 10001 // 默认服务错误
	ParamErrorCode             = 10002 // 参数错误
	HTTPQueryErrorCode         = 10003 // HTTP请求出错
	AuthorizationFailedErrCode = 10004 // 鉴权失败
	UnexpectedTypeErrorCode    = 10005 // 未知类型
	NotImplementErrorCode      = 10006 // 未实装

)
