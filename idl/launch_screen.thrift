namespace go launch_screen
include"model.thrift"

struct CreateImageRequest {
    1: required i64 pic_type,
    2: optional i64 duration,
    3: optional string href,
    4: required binary image,
    5: required i64 start_at,
    6: required i64 end_at,
    7: required i64 s_type,
    8: required i64 frequency,
    9: required i64 start_time,
    10:required i64 end_time,
    11:required string text,
    12:required string regex,
    13:i64 buffer_count,
    14:string suffix,
}

struct CreateImageResponse{
    1:model.BaseResp base,
    2:optional model.Picture picture,
}

struct GetImageRequest{
    1:required i64 picture_id,

}

struct GetImageResponse{
    1:model.BaseResp base,
    2:optional model.Picture picture,
}


struct ChangeImagePropertyRequest {
    1: required i64 pic_type,// 1为空，2为页面跳转，3为app跳转
    2: optional i64 duration,
    3: optional string href,// 连接
    4: required i64 start_at,
    5: required i64 end_at,
    6: required i64 s_type,
    7: required i64 frequency,
    8: required i64 start_time,// 比如6表示6点
    9:required i64 end_time,
    10:required string text,// 描述图片
    11:required i64 picture_id,
    12:required string regex,

}

struct ChangeImagePropertyResponse{
    1:model.BaseResp base,
    2:optional model.Picture picture,
}

struct ChangeImageRequest {
    1:required i64 picture_id,
    2:required binary image,
    3:i64 buffer_count,
    4:string suffix,
}

struct ChangeImageResponse{
    1:model.BaseResp base,
    2:optional model.Picture picture,
}

struct DeleteImageRequest{
    1:required i64 picture_id,
}

struct DeleteImageResponse{
    1:model.BaseResp base,
}

struct MobileGetImageRequest{
    1:required i64 s_type,
    2:required string student_id,
    3:optional string college,
    4:required string device,
}


struct MobileGetImageResponse{
    1:model.BaseResp base,
    2:optional i64 count,
    3:optional list<model.Picture> picture_list,
}

struct AddImagePointTimeRequest{
    1:required i64 picture_id,
}

struct AddImagePointTimeResponse{
    1:model.BaseResp base,
    2:optional model.Picture picture,
}

service LaunchScreenService{
    CreateImageResponse CreateImage(1:CreateImageRequest req)(streaming.mode="client"), // 开启流式传输
    GetImageResponse GetImage(1:GetImageRequest req),
    ChangeImagePropertyResponse ChangeImageProperty(1:ChangeImagePropertyRequest req),
    ChangeImageResponse ChangeImage(1:ChangeImageRequest req)(streaming.mode="client"),
    DeleteImageResponse DeleteImage(1:DeleteImageRequest req),
    MobileGetImageResponse MobileGetImage(1:MobileGetImageRequest req),
    AddImagePointTimeResponse AddImagePointTime(1:AddImagePointTimeRequest req),
}
