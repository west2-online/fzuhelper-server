namespace go api

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct User{
    1:i64 id,
    2:string account,
    3:string name,
}

struct Picture{
    1:i64 id,
    2:string url,
    3:string href,
    4:string text,
    5:i64 pic_type,
    6:i64 show_times,
    7:i64 point_times,
    8:i64 duration,
    9:i64 s_type,
    10:i64 frequency,
}

struct RegisterRequest {
    1: required string account,
    2: required string name,
    3: required string password,
}

struct RegisterResponse {
    1: BaseResp base,
    2: optional i64 user_id,
}

struct LoginRequest {
    1: string account,
    2: string password,
}

struct LoginResponse {
    1: BaseResp base,
    2: optional string data,
}

service UserService {
    LoginResponse Login(1: LoginRequest req)(api.post="/launch_screen/api/login"),

    //test for backend
    RegisterResponse Register(1: RegisterRequest req)(api.post="/launch_screen/api/register"),
}
