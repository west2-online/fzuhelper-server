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
    2: optional string token,
}

service UserService {
    LoginResponse Login(1: LoginRequest req)(api.post="/launch_screen/api/login"),

    //test for backend
    RegisterResponse Register(1: RegisterRequest req)(api.post="/launch_screen/api/register"),
}

struct CreateImageRequest {
    1: required i64 pic_type,//1为空，2为页面跳转，3为app跳转
    2: optional i64 duration,
    3: optional string href,//连接
    4: binary image,
    5: required i64 start_at,
    6: required i64 end_at,
    7: required i64 s_type,
    8: required i64 frequency,//单日最大展示次数
    9: required i64 start_time,//比如6表示6点
    10:required i64 end_time,
    11:required string text,//描述图片
    12:string regex,//正则匹配项

    13:i64 user_id,//get by token
}

struct CreateImageResponse{
    1:BaseResp base,
    2:optional Picture picture,
}

struct GetImageRequest{
    1:required i64 picture_id,

    2:i64 user_id,
}

struct GetImageResponse{
    1:BaseResp base,
    2:optional Picture picture,
}

struct GetImagesByUserIdRequest{
    1:i64 user_id,
}

struct GetImagesByUserIdResponse{
    1:BaseResp base,
    2:optional list<Picture> picture_list,
}

struct ChangeImagePropertyRequest {
    1: required i64 pic_type,//1为空，2为页面跳转，3为app跳转
    2: optional i64 duration,
    3: optional string href,//连接
    4: required i64 start_at,
    5: required i64 end_at,
    6: required i64 s_type,
    7: required i64 frequency,
    8: required i64 start_time,//比如6表示6点
    9:required i64 end_time,
    10:required string text,//描述图片
    11:string regex,//正则匹配项

    12:i64 user_id,//get by token
}

struct ChangeImagePropertyResponse{
    1:BaseResp base,
    2:optional Picture picture,
}

struct ChangeImageRequest {
    1:required i64 picture_id,
    2:required binary image,

    3:i64 user_id,
}

struct ChangeImageResponse{
    1:BaseResp base,
    2:optional Picture picture,
}

struct DeleteImageRequest{
    1:required i64 picture_id,

    2:i64 user_id,
}

struct DeleteImageResponse{
    1:BaseResp base,
    2:optional Picture picture,
}

struct MobileGetImageRequest{
    1:required i64 type,
    2:required i64 student_id,
    3:required string college,
    4:required string device,
}

struct MobileGetImageResponse{
    1:BaseResp base,
    2:optional Picture picture,
}

struct AddImagePointTimeRequest{
    1:required i64 picture_id,
}

struct AddImagePointTimeResponse{
    1:BaseResp base,
    2:optional Picture picture,
}

service LaunchScreenService{
    CreateImageResponse CreateImage(1:CreateImageRequest req)(api.post="/launch_screen/api/image"),
    GetImagesByUserIdResponse GetImage(1:GetImageRequest req)(api.get="/launch_screen/api/image"),
    GetImagesByUserIdResponse GetImagesByUserId(1:GetImagesByUserIdRequest req)(api.get="/launch_screen/api/images"),
    ChangeImagePropertyResponse ChangeImageProperty(1:ChangeImagePropertyRequest req)(api.put="/launch_screen/api/image"),
    ChangeImageResponse ChangeImage(1:ChangeImageRequest req)(api.put="/launch_screen/api/image/img"),
    DeleteImageResponse DeleteImage(1:DeleteImageRequest req)(api.delete="/launch_screen/api/image"),
    MobileGetImageResponse MobileGetImage(1:MobileGetImageRequest req)(api.get="/launch_screen/api/screen"),
    AddImagePointTimeResponse AddImagePointTime(1:AddImagePointTimeRequest req)(api.get="/launch_screen/api/image/point"),
}
