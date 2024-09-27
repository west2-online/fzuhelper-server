namespace go user

struct BaseResp {
    1: i64 code,
    2: string msg,
}



//just for backend testing
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required BaseResp base,
    2: required string id
    3: required list<string> cookies
}

service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest req)
}
