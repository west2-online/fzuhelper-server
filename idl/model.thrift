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
}
