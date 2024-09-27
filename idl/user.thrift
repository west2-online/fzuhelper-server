namespace go user

//just for backend testing
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required string id
    2: required list<string> cookies
}

service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest request)
}
