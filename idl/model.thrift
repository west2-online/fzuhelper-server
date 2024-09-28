namespace go model

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
