// This file is be designed to define any common error so that we can use it in any service simply.

package errno

var (
	// Success
	Success = NewErrNo(SuccessCode, "Success")

	ParamError   = NewErrNo(ParamErrorCode, "parameter error")
	ParamEmpty   = NewErrNo(ParamEmptyCode, "some params that required are empty")
	ServiceError = NewErrNo(ServiceErrorCode, "service is unable to start successfully")

	AuthFailedError      = NewErrNo(AuthErrorCode, "authorization failed")
	BizError             = NewErrNo(BizErrorCode, "business error")
	InternalServiceError = NewErrNo(InternalServiceErrorCode, "internal service error")
)
