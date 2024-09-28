namespace go model

struct BaseResp {
    1: i64 code,
    2: string msg,
}

//由前端给的登陆信息，包括id和cookies, 这个struct仅用于测试返回数据，因为登录实现在前端完成，不会在实际项目中使用
struct LoginData {
    1: required string id
    2: required list<string> cookies
}

//Classroom
struct Classroom {
    1: required string build
    2: required string location
    3: required string capacity
    4: required string type
}

