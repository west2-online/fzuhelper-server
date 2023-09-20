namespace go screen

struct BaseResp {
    1: i64 code
    2: string msg
}

struct Picture {
    1: i64 picture_id
    2: string create_at
    3: string update_at
    4: string url
    5: string herf
    6: string text
    7: i8 pic_type
    8: i64 show_times
    9: i64 point_times
    10: i64 duration
    11: string start_at
    12: string end_at
    13: i8 s_type
    14: i64 frequency
}

struct CreatePictureRequest {
    1: string href 
    2: string text 
    3: i8 picType 
    4: i64 duration 
    5: i64 start_at 
    6: i64 end_at 
    7: i64 start_time 
    8: i64 end_time 
    9: i64 frequency 
    10: i8 s_type 
    11: string token
    12: binary imgfile
}

struct CreatePictureResponse{
    1: BaseResp base
    2: Picture picture
}

struct GetPictureRequest{
    1: string token
    2: i64 picture_id // 0-获取所有 其他-获取指定id图片
}

struct GetPictureResponse{
    1: BaseResp base
    2: list<Picture> picture
    3: i64 total
}

struct PutPictureRequset{
    1: string href 
    2: string text 
    3: i8 picType 
    4: i64 duration 
    5: i64 start_at 
    6: i64 end_at 
    7: i64 start_time 
    8: i64 end_time 
    9: i64 frequency 
    10: i8 s_type 
    11: i64 picture_id
    12: string token 
}

struct PutPictureResponse{
    1: BaseResp base
    2: Picture picture
}

struct PutPictureImgRequset{
    1: string token
    2: binary imgfile
    3: i64 picture_id
}

struct DeletePictureRequest{
    1: string token
    2: i64 picture_id
}

struct DeletePictureResponse{
    1: BaseResp base
    2: Picture picture
}

struct RetPictureRequest{ // android
    1: string token
    2: string type
}

struct RetPictureResponse{
    1: BaseResp base
    2: list<Picture> picture
    3: i64 total
}

struct AddPointRequest{
    1: string token
    2: i64 picture_id
}

struct AddPointResponse{
    1: BaseResp base
}

service LaunchScreenService {
    CreatePictureResponse PictureCreate(1:CreatePictureRequest req)
    GetPictureResponse PictureGet(1:GetPictureRequest req)
    PutPictureResponse PictureUpdate(1:PutPictureRequset req)
    PutPictureResponse PictureImgUpdate(1:PutPictureImgRequset req)
    DeletePictureResponse PictureDelete(1:DeletePictureRequest req)

    RetPictureResponse RetPicture(1:RetPictureRequest req)
    AddPointResponse AddPoint(1:AddPointRequest req)
}