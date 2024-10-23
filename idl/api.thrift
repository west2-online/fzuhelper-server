namespace go api
include "model.thrift"

//classroom
struct EmptyClassroomRequest {
    1: required string date
    2: required string campus
    3: required string startTime;//节数
    4: required string endTime;
}

struct EmptyClassroomResponse {
    1: required list<model.Classroom> classrooms
}

service ClassRoomService {
    EmptyClassroomResponse GetEmptyClassrooms(1: EmptyClassroomRequest request)(api.get="/api/v1/common/classroom/empty")
}

//user
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required string id
    2: required list<string> cookies
}

service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest request)(api.get="/api/v1/jwch/user/login"),
}


//launch_screen
struct CreateImageRequest {
    1: required i64 pic_type,
    2: optional i64 duration,
    3: optional string href,
    4: binary image,
    5: required i64 start_at,
    6: required i64 end_at,
    7: required i64 s_type,
    8: required i64 frequency,
    9: required i64 start_time,
    10:required i64 end_time,
    11:required string text,
    12:required string regex,
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
    11:required i64 picture_id,
    12:required string regex,
}

struct ChangeImagePropertyResponse{
    1:model.BaseResp base,
    2:optional model.Picture picture,
}

struct ChangeImageRequest {
    1:required i64 picture_id,
    2: binary image,
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
    2:optional model.Picture picture,
}

struct MobileGetImageRequest{
    1:required i64 type,
    2:required i64 student_id,
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
    CreateImageResponse CreateImage(1:CreateImageRequest req)(api.post="/launch_screen/api/image"),
    GetImageResponse GetImage(1:GetImageRequest req)(api.get="/launch_screen/api/image"),
    ChangeImagePropertyResponse ChangeImageProperty(1:ChangeImagePropertyRequest req)(api.put="/launch_screen/api/image"),
    ChangeImageResponse ChangeImage(1:ChangeImageRequest req)(api.put="/launch_screen/api/image/img"),
    DeleteImageResponse DeleteImage(1:DeleteImageRequest req)(api.delete="/launch_screen/api/image"),
    MobileGetImageResponse MobileGetImage(1:MobileGetImageRequest req)(api.get="/launch_screen/api/screen"),
    AddImagePointTimeResponse AddImagePointTime(1:AddImagePointTimeRequest req)(api.get="/launch_screen/api/image/point"),
}



//paper


struct ListDirFilesRequest {
    1: required string path,

}

struct ListDirFilesResponse {
    1: required model.UpYunFileDir dir,
}

struct GetDownloadUrlRequest {
    1: required string url,
}

struct GetDownloadUrlResponse {
    1: required string url,
}


service PaperService {
    ListDirFilesResponse ListDirFiles(1: ListDirFilesRequest req) (api.get="/api/v1/paper/list"),
    GetDownloadUrlResponse GetDownloadUrl(1: GetDownloadUrlRequest req) (api.get="/api/v1/paper/download"),
}

//academic
struct GetScoresRequest {
}

struct GetScoresResponse {
    1: required list<model.Score> scores
}

struct GetGPARequest {
}

struct GetGPAResponse {
    1: required model.GPABean gpa
}

struct GetCreditRequest {
}

struct GetCreditResponse {
    1: required list<model.Credit> major
}

struct GetUnifiedExamRequest {
}

struct GetUnifiedExamResponse {
    1: required list<model.UnifiedExam> unifiedExam
}

service AcademicService {
    GetScoresResponse GetScores(1:GetScoresRequest req)(api.get="/api/v1/jwch/academic/scores")
    GetGPAResponse GetGPA(1:GetGPARequest req)(api.get="/api/v1/jwch/academic/gpa")
    GetCreditResponse GetCredit(1:GetCreditRequest req)(api.get="/api/v1/jwch/academic/credit")
    GetUnifiedExamResponse GetUnifiedExam(1:GetUnifiedExamRequest req)(api.get="/api/v1/jwch/academic/unifiedExam")
}