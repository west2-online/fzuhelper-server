namespace go user

include"model.thrift"

struct RegisterRequest {
    1: required string account,
    2: required string name,
    3: required string password,
}

struct RegisterResponse {
    1: model.BaseResp base,
    2: optional i64 user_id,
}

struct LoginRequest {
    1: string account,
    2: string password,
}

struct LoginResponse {
    1: model.BaseResp base,
    2: optional string token,
}

service UserService {
    LoginResponse Login(1: LoginRequest req)(api.post="/launch_screen/api/login"),

    //test for backend
    RegisterResponse Register(1: RegisterRequest req)(api.post="/launch_screen/api/register"),
}
