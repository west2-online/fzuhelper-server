// This file is designed to define any error code
package errno

const (
	// For microservices
	SuccessCode = 10000
	SuccessMsg  = "ok"

	// Error
	/*
		200xx: 参数错误，Param 打头
		300xx: 鉴权错误，Auth 打头
		400xx: 业务错误，Biz 打头
		500xx: 内部错误，Internal 打头
	*/
	ParamErrorCode         = 20001 // 参数错误
	ParamEmptyCode         = 20002 // 参数为空
	ParamInvalidCode       = 20003 // 参数无效
	ParamMissingCode       = 20004 // 参数缺失
	ParamTooLongCode       = 20005 // 参数过长
	ParamTooShortCode      = 20006 // 参数过短
	ParamTypeCode          = 20007 // 参数类型错误
	ParamFormatCode        = 20008 // 参数格式错误
	ParamRangeCode         = 20009 // 参数范围错误
	ParamValueCode         = 20010 // 参数值错误
	ParamFileNotExistCode  = 20011 // 文件不存在
	ParamFileReadErrorCode = 20012 // 文件读取错误

	AuthErrorCode     = 30001 // 鉴权错误
	AuthInvalidCode   = 30002 // 鉴权无效
	AuthExpiredCode   = 30003 // 鉴权过期
	AuthMissingCode   = 30004 // 鉴权缺失
	AuthNotEnoughCode = 30005 // 鉴权不足

	BizErrorCode           = 40001 // 业务错误
	BizLogicCode           = 40002 // 业务逻辑错误
	BizLimitCode           = 40003 // 业务限制错误
	BizNotExist            = 40005 // 业务不存在错误
	BizFileUploadErrorCode = 40006 // 文件上传错误(service 层)

	InternalServiceErrorCode   = 50001 // 未知服务错误
	InternalDatabaseErrorCode  = 50002 // 数据库错误
	InternalRedisErrorCode     = 50003 // Redis错误
	InternalNetworkErrorCode   = 50004 // 网络错误
	InternalTimeoutErrorCode   = 50005 // 超时错误
	InternalIOErrorCode        = 50006 // IO错误
	InternalJSONErrorCode      = 50007 // JSON错误
	InternalXMLErrorCode       = 50008 // XML错误
	InternalURLEncodeErrorCode = 50009 // URL编码错误
	InternalHTTPErrorCode      = 50010 // HTTP错误
	InternalHTTP2ErrorCode     = 50011 // HTTP2错误
	InternalGRPCErrorCode      = 50012 // GRPC错误
	InternalThriftErrorCode    = 50013 // Thrift错误
	InternalProtobufErrorCode  = 50014 // Protobuf错误
	InternalSQLErrorCode       = 50015 // SQL错误
	InternalNoSQLErrorCode     = 50016 // NoSQL错误
	InternalORMErrorCode       = 50017 // ORM错误
	InternalQueueErrorCode     = 50018 // 队列错误
	InternalETCDErrorCode      = 50019 // ETCD错误
	InternalTraceErrorCode     = 50020 // Trace错误
)
