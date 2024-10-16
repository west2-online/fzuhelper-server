namespace go user
include "model.thrift"

//just for backend testing
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required model.BaseResp base,
    2: required string id
    3: required list<string> cookies
}

struct GetValidateCodeRequest{
    1:required string image,
}

struct GetValidateCodeResponse{
    1:model.BaseResp base,
    2:optional string code,
}


service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest req),
    GetValidateCodeResponse GetValidateCode(1: GetValidateCodeRequest req),
}
