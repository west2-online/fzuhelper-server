// This file is be designed to define any common error so that we can use it in any service simply.

package errno

var (
	// Success
	Success = NewErrNo(SuccessCode, "Success")

	ParamError = NewErrNo(ParamErrorCode, "parameter error")
	ParamEmpty = NewErrNo(ParamEmptyCode, "some params that required are empty")

	AuthFailedError      = NewErrNo(AuthErrorCode, "authorization failed")
	BizError             = NewErrNo(BizErrorCode, "business error")
	InternalServiceError = NewErrNo(InternalServiceErrorCode, "internal service error")

	UserExistedError  = NewErrNo(InternalDatabaseErrorCode, "user existed")
	UserNonExistError = NewErrNo(InternalDatabaseErrorCode, "user didn't exist")
	SuffixError       = NewErrNo(ParamErrorCode, "invalid file")
	UpcloudError      = NewErrNo(BizFileUploadErrorCode, "upload to upcloud error")
)
