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

struct RegisterRequest {
    1: required string number,
    2: required string password,
}

struct RegisterResponse {
    1: model.BaseResp base,
    2: optional i64 user_id,
}

struct LoginRequest {
    1: string number,
    2: string password,
}

struct LoginResponse {
    1: model.BaseResp base,
    2: optional string token,
}

service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest req),

    LoginResponse Login(1: LoginRequest req)(api.post="/user/api/login"),
    //test for backend
    RegisterResponse Register(1: RegisterRequest req)(api.post="/uaer/api/register"),
}
