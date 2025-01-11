namespace go user

include "model.thrift"

// for backend testing
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required model.BaseResp base,
    2: required string id
    3: required list<string> cookies
}

struct GetUserInfoRequest{
}

struct GetUserInfoResponse{
    1: required model.BaseResp base,
    2: optional model.UserInfo data,
}

service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest req),
    GetUserInfoResponse GetUserInfo(1: GetUserInfoRequest request),
}
