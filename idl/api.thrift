namespace go api
include "model.thrift"

// classroom
struct EmptyClassroomRequest {
    1: required string date
    2: required string campus
    3: required string startTime;
    4: required string endTime;
}

struct EmptyClassroomResponse {
    1: required list<model.Classroom> classrooms
}

service ClassRoomService {
    EmptyClassroomResponse GetEmptyClassrooms(1: EmptyClassroomRequest request)(api.get="/api/v1/common/classroom/empty")
}

// user
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required string id
    2: required list<string> cookies
}

struct ValidateCodeRequest {
    1: required string image
}

struct ValidateCodeResponse {
}
// Android兼容
struct ValidateCodeForAndroidRequest {
    1: required string validateCode
}

struct ValidateCodeForAndroidResponse {
}

service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest request)(api.get="/api/v1/jwch/user/login"),
    ValidateCodeResponse ValidateCode(1: ValidateCodeRequest request)(api.post="/api/v1/jwch/user/validateCode")
    ValidateCodeForAndroidResponse ValidateCodeForAndroid(1: ValidateCodeForAndroidRequest request)(api.post="/api/login/validateCode")
}

// course
struct CourseListRequest {
    1: required string term
}

struct CourseListResponse {
    1: required model.BaseResp base
    2: required list<model.Course> data
}

service CourseService {
    CourseListResponse GetCourseList(1: CourseListRequest req)(api.get="/api/v1/jwch/course/list")
}

// launch_screen
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
    10: required i64 end_time,
    11: required string text,
    12: required string regex,
}

struct CreateImageResponse{
    1: model.BaseResp base,
    2: optional model.Picture picture,
}

struct GetImageRequest{
    1: required i64 picture_id,

}

struct GetImageResponse{
    1: model.BaseResp base,
    2: optional model.Picture picture,
}

struct ChangeImagePropertyRequest {
    1: required i64 pic_type, // 1 为空，2 为页面跳转，3 为 APP 跳转
    2: optional i64 duration,
    3: optional string href, // 链接
    4: required i64 start_at,
    5: required i64 end_at,
    6: required i64 s_type,
    7: required i64 frequency,
    8: required i64 start_time, // 例：6 表示 6点
    9: required i64 end_time,
    10: required string text, // 描述图片
    11: required i64 picture_id,
    12: required string regex,
}

struct ChangeImagePropertyResponse{
    1: model.BaseResp base,
    2: optional model.Picture picture,
}

struct ChangeImageRequest {
    1: required i64 picture_id,
    2: binary image,
}

struct ChangeImageResponse{
    1: model.BaseResp base,
    2: optional model.Picture picture,
}

struct DeleteImageRequest{
    1: required i64 picture_id,
}

struct DeleteImageResponse{
    1: model.BaseResp base,
}

struct MobileGetImageRequest{
    1: required i64 type,
    2: required string student_id,
    3: optional string college,
    4: required string device,
}

struct MobileGetImageResponse{
    1: model.BaseResp base,
    2: optional i64 count,
    3: optional list<model.Picture> picture_list,
}

struct AddImagePointTimeRequest{
    1: required i64 picture_id,
}

struct AddImagePointTimeResponse{
    1: model.BaseResp base,
    2: optional model.Picture picture,
}

service LaunchScreenService{
    CreateImageResponse CreateImage(1: CreateImageRequest req) (api.post="/launch_screen/api/image"),
    GetImageResponse GetImage(1: GetImageRequest req) (api.get="/launch_screen/api/image"),
    ChangeImagePropertyResponse ChangeImageProperty(1: ChangeImagePropertyRequest req) (api.put="/launch_screen/api/image"),
    ChangeImageResponse ChangeImage(1: ChangeImageRequest req) (api.put="/launch_screen/api/image/img"),
    DeleteImageResponse DeleteImage(1: DeleteImageRequest req) (api.delete="/launch_screen/api/image"),
    MobileGetImageResponse MobileGetImage(1: MobileGetImageRequest req) (api.get="/launch_screen/api/screen"),
    AddImagePointTimeResponse AddImagePointTime(1: AddImagePointTimeRequest req) (api.get="/launch_screen/api/image/point"),
}

// paper
struct ListDirFilesRequest {
    1: required string path,
}

struct ListDirFilesResponse {
    1: required model.UpYunFileDir dir,
}

struct GetDownloadUrlRequest {
    1: required string filepath,
}

struct GetDownloadUrlResponse {
    1: required string url,
}

// 兼容
struct ListDirFilesForAndroidRequest {
    1: required string path,
}

struct ListDirFilesForAndroidResponse {

}

struct GetDownloadUrlForAndroidRequest {
    1: required string filepath,
}

struct GetDownloadUrlForAndroidResponse {

}


service PaperService {
    ListDirFilesResponse ListDirFiles(1: ListDirFilesRequest req) (api.get="/api/v1/paper/list"),
    GetDownloadUrlResponse GetDownloadUrl(1: GetDownloadUrlRequest req) (api.get="/api/v1/paper/download"),

    // 兼容安卓
    ListDirFilesForAndroidResponse ListDirFilesForAndroid(1: ListDirFilesForAndroidRequest req) (api.get="/api/v1/list")
    GetDownloadUrlForAndroidResponse GetDownloadUrlForAndroid(1: GetDownloadUrlForAndroidRequest req) (api.get="/api/v1/downloadUrl")
}

// academic
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



// url

struct APILoginRequest {
    1: required string password
}

struct APILoginResponse {
    
}

struct UploadVersionInfoRequest {
    1: required string password
    2: required string type
    3: required string version
    4: required string code
    5: required string feature
    6: required string url
}

struct UploadVersionInfoResponse {

}

struct GetUploadParamsRequest {
    1: required string password
}

struct GetUploadParamsResponse {
    1: required string policy,
    2: required string authorization,
}


struct GetDownloadReleaseRequest {

}

struct GetDownloadReleaseResponse {

}

struct GetDownloadBetaRequest {

}

struct GetDownloadBetaResponse {

}

struct GetReleaseVersionRequest {

}

struct GetReleaseVersionResponse {


}
struct GetBetaVersionRequest {

}

struct GetBetaVersionResponse{

}

struct GetCloudSettingRequest {
    1: optional string account,
    2: optional string version,
    3: optional string beta,
    4: optional string phone,
    5: optional string isLogin,
    6: optional string loginType,
}

struct GetCloudSettingResponse {

}

struct GetAllCloudSettingRequest {

}

struct GetAllCloudSettingResponse {

}

struct SetAllCloudSettingRequest {
    1: required string password
    2: required string setting
}

struct SetAllCloudSettingResponse {

}

struct TestSettingRequest {
    1: required string setting
    2: required string account
    3: required string version
    4: required string beta
    5: required string phone
    6: required string isLogin
    7: required string loginType
}

struct TestSettingResponse {

}

struct DumpVisitRequest {

}

struct DumpVisitResponse {

}

struct FZUHelperCSSRequest{

}

struct FZUHelperCSSResponse {

}

struct FZUHelperHTMLRequest {

}

struct FZUHelperHTMLResponse {

}

struct UserAgreementHTMLRequest {

}

struct UserAgreementHTMLResponse {

}

service UrlService {
    APILoginResponse APILogin(1:APILoginRequest req) (api.post="/api/v1/url/login")
    UploadVersionInfoResponse UploadVersionInfo(1:UploadVersionInfoRequest req) (api.post="/api/v1/url/api/upload")
    GetUploadParamsResponse GetUploadParams(1:GetUploadParamsRequest req) (api.post="/api/v1/url/api/uploadparams")
    GetDownloadReleaseResponse GetDownloadRelease(1:GetDownloadReleaseRequest req) (api.get="/api/v1/url/release.apk")
    GetDownloadBetaResponse GetDownloadBeta(1: GetDownloadBetaRequest req) (api.get="/api/v1/url/beta.apk")
    GetReleaseVersionResponse GetReleaseVersion(1:GetReleaseVersionRequest req) (api.get="/api/v1/url/version.json")
    GetBetaVersionResponse GetBetaVersion(1: GetBetaVersionRequest req) (api.get="/api/v1/url/versionbeta.json")
    GetCloudSettingResponse GetCloudSetting(1: GetCloudSettingRequest req) (api.get="/api/v1/url/settings.php")
    GetAllCloudSettingResponse GetAllCloudSetting(1: GetAllCloudSettingRequest req) (api.get="/api/v1/url/api/getcloud")
    SetAllCloudSettingResponse SetAllCloudSetting(1: SetAllCloudSettingRequest req) (api.post="/api/v1/url/api/setcloud")
    TestSettingResponse TestSetting(1: TestSettingRequest req) (api.post="/api/v1/url/api/test")
    DumpVisitResponse DumpVisit(1: DumpVisitRequest req) (api.get="/api/v1/url/dump")
    FZUHelperCSSResponse FZUHelperCSS(1: FZUHelperCSSRequest req) (api.get="/api/v1/url/onekey/FZUHelper.css")
    FZUHelperHTMLResponse FZUHelperHTML(1: FZUHelperHTMLRequest req) (api.get="/api/v1/url/onekey/FZUHelper.html")
    UserAgreementHTMLResponse UserAgreementHTML(1: UserAgreementHTMLRequest req) (api.get="/api/v1/url/onekey/UserAgreement.html")
}