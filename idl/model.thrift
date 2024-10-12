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

//Classroom 前端想要返回的fields
struct Classroom {
    1: required string build
    2: required string location
    3: required string capacity
    4: required string type
}

struct User{
    1:i64 id,
    2:string account,
    3:string name,
}

struct Picture{
    1:i64 id,
    2:i64 user_id,
    3:string url,
    4:string href,
    5:string text,
    6:i64 pic_type,
    7:optional i64 show_times,
    8:optional i64 point_times,
    9:i64 duration,
    10:optional i64 s_type,
    11:i64 frequency,
    12:i64 start_at,
    13:i64 end_at,
    14:i64 start_time,
    15:i64 end_time,
    16:i64 student_id,
    17:i64 device_type,
}
