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

	// redis
	RedisError = NewErrNo(InternalRedisErrorCode, "redis error")
)
