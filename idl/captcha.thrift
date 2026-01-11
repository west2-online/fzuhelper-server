namespace go captcha

include "model.thrift"

struct ValidateCodeRequest {
    1: required string image,
}

struct ValidateCodeResponse {
    1: required model.BaseResp base,
    2: required string data,
}

struct ValidateCodeForAndroidRequest {
    1: required string validateCode,
}

struct ValidateCodeForAndroidResponse {
    1: required string code,
    2: required string message,
}

service CaptchaService {
    ValidateCodeResponse ValidateCode(1: ValidateCodeRequest req)(api.post='/api/v1/user/validate-code'),
    ValidateCodeForAndroidResponse ValidateCodeForAndroid(1: ValidateCodeForAndroidRequest req)(api.post='/api/login/validateCode'),
}