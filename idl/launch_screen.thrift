namespace go launch_screen
include"model.thrift"

struct CreateImageRequest {
    1: required i64 pic_type,//1为空，2为页面跳转，3为app跳转
    2: optional i64 duration,
    3: optional string href,//连接
    4: required binary image,
    5: required string start_at,
    6: required string end_at,
    7: required i64 s_type,
    8: required i64 frequency,
    9: required i64 start_time,//比如6表示6点
    10:required i64 end_time,
    11:required string text,//描述图片
    12:string regex,//正则匹配项

    13:i64 user_id,//get by token
}

struct CreateImageResponse{
    1:model.BaseResp base,
    2:optional model.Picture picture,
}

struct GetImageRequest{
    1:required i64 picture_id,

    2:i64 user_id,
}

struct GetImageResponse{
    1:required i64 picture_id,

    2:i64 user_id,
}
