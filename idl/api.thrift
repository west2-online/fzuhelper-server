namespace go api
include "model.thrift"
# 重构的服务 url 统一前缀为 /api/v1，兼容部分不做任何修改
# 其中有使用鉴权的前缀为 /jwch，主要表现为 Header 需要 id 和 cookies 的接口

## ----------------------------------------------------------------------------
## classroom 空教室、考表
## ----------------------------------------------------------------------------
struct EmptyClassroomRequest {
    1: required string date
    2: required string campus
    3: required string startTime;
    4: required string endTime;
}

struct EmptyClassroomResponse {
    1: optional list<model.Classroom> classrooms
}

struct ExamRoomInfoRequest {
    1: required string term
}

struct ExamRoomInfoResponse {
    1: optional list<model.ExamRoomInfo> examRoomInfos
}

service ClassRoomService {
    // 查询空教室
    EmptyClassroomResponse GetEmptyClassrooms(1: EmptyClassroomRequest request)(api.get="/api/v1/common/classroom/empty")
    // 查询考表
    ExamRoomInfoResponse GetExamRoomInfo(1: ExamRoomInfoRequest request)(api.get="/api/v1/jwch/classroom/exam")
}

## ----------------------------------------------------------------------------
## user 用户（如登录、鉴权）
## ----------------------------------------------------------------------------
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required string id
    2: required string cookies
}

struct ValidateCodeRequest {
    1: required string image
}

struct ValidateCodeResponse {}

// Android兼容
struct ValidateCodeForAndroidRequest {
    1: required string validateCode
}

struct ValidateCodeForAndroidResponse {
}

struct GetAccessTokenRequest {
}

struct GetAccessTokenResponse {
    1: string code;
    2: string message;
}

struct RefreshTokenRequest {
}

struct RefreshTokenResponse {
    1: string code;
    2: string message;
}

struct TestAuthRequest{
}

struct TestAuthResponse{
    1: string message
}

struct GetUserInfoRequest{
}

struct GetUserInfoResponse{
    1: required model.BaseResp base,
    2: optional model.UserInfo data,
}

struct GetLoginDataForYJSYRequest{
    1: required string id
    2: required string password
}

struct GetLoginDataForYJSYResponse{
    1: required string id
    2: required string cookies
}
struct GetInvitationCodeRequest{
        1: optional bool isRefresh // 刷新邀请码
}
struct GetInvitationCodeResponse{
        1: required model.BaseResp base,
        2: required string invitation_code,
}
struct BindInvitationRequest{
        1: required string invitation_code
}
struct BindInvitationResponse{
        1: required model.BaseResp base,
}
struct GetFriendListRequest{

}
struct GetFriendListResponse{
     1: required model.BaseResp base,
    2: required list<model.UserInfo> data
}
struct DeleteFriendRequest{
    1:required string id
}
struct DeleteFriendResponse{
         1: required model.BaseResp base,
}
service UserService {
    // 后端自动登录（含验证码识别），该接口默认不提供给客户端，仅供测试
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest request)(api.get="/api/v1/internal/user/login"), # 后端内部测试接口使用，使用 internal 前缀做区别
    // 后端自动登录（研究生，无需验证码），该接口默认不提供给客户端，仅供测试
    GetLoginDataForYJSYResponse GetGetLoginDataForYJSY(1:GetLoginDataForYJSYRequest request)(api.get="/api/v1/internal/yjsy/user/login"), # 后端内部测试接口使用，使用 internal 前缀做区别
    // 自动识别验证码
    ValidateCodeResponse ValidateCode(1: ValidateCodeRequest request)(api.post="/api/v1/user/validate-code")
    // 自动识别验证码（安卓兼容）
    ValidateCodeForAndroidResponse ValidateCodeForAndroid(1: ValidateCodeForAndroidRequest request)(api.post="/api/login/validateCode") # 兼容安卓端
    // 获取 Access-Token
    GetAccessTokenResponse GetToken(1: GetAccessTokenRequest request)(api.get="/api/v1/login/access-token"),
    // 获取 Refresh-Token
    RefreshTokenResponse RefreshToken(1: RefreshTokenRequest request)(api.get="/api/v1/login/refresh-token"),
    // 测试含鉴权的 ping 功能
    TestAuthResponse TestAuth(1: TestAuthRequest request)(api.get="/api/v1/jwch/ping")
    // 获取用户信息
    GetUserInfoResponse GetUserInfo(1: GetUserInfoRequest request)(api.get="/api/v1/jwch/user/info")
    // 获取邀请码
    GetInvitationCodeResponse GetInvitationCode(1:GetInvitationCodeRequest request)(api.get="/api/v1/user/friend/invite")
    // 绑定邀请关系
    BindInvitationResponse BindInvitation(1:BindInvitationRequest request)(api.post = "/api/v1/user/friend/bind")
    // 查看好友列表
    GetFriendListResponse GetFriendList(1:GetFriendListRequest request)(api.get = "/api/v1/user/friend/info")
    // 删除好友
    DeleteFriendResponse DeleteFriend(1:DeleteFriendRequest request)(api.delete = "/api/v1/user/friend/delete")
}
## ----------------------------------------------------------------------------
## course 课表
## ----------------------------------------------------------------------------
struct CourseListRequest {
    1: required string term
    2: optional bool is_refresh
}

struct CourseListResponse {
    1: required model.BaseResp base
    2: required list<model.Course> data
}

struct CourseTermListRequest{}

struct CourseTermListResponse{
    1: required model.BaseResp base
    2: required list<string> data
}

struct GetCalendarTokenRequest {
}

struct GetCalendarTokenResponse {
    1: required string token
}

struct SubscribeCalendarRequest {
    1:required string token
}

struct SubscribeCalendarResponse {
    1: binary ics
}

struct GetLocateDateRequest{}

struct GetLocateDateResponse{
    1: optional model.LocateDate locateDate
}

struct GetFriendCourseRequest {
    1: required string term
    2: required string id
}

struct GetFriendCourseResponse {
    1: required model.BaseResp base
    2: required list<model.Course> data
}

service CourseService {
    // 获取课表
    CourseListResponse GetCourseList(1: CourseListRequest req)(api.get="/api/v1/jwch/course/list")
    // 获取学期
    CourseTermListResponse GetTermList(1: CourseTermListRequest req)(api.get="/api/v1/jwch/term/list")
    // 获取日历订阅 token
    GetCalendarTokenResponse GetCalendar(1: GetCalendarTokenRequest req)(api.get="/api/v1/jwch/course/calendar/token")

    // 由手机端的日历 app 直接发起的请求，无双 token 保护（即 url "/jwch" 前缀）
    SubscribeCalendarResponse SubscribeCalendar(1: SubscribeCalendarRequest req)(api.get="/api/v1/course/calendar/subscribe")
    // 获取当前周数、学期、学年
    GetLocateDateResponse GetLocateDate(1:GetLocateDateRequest req)(api.get="/api/v1/course/date")
    // 获取好友课表
    GetFriendCourseResponse GetFriendCourse(1:GetFriendCourseRequest req)(api.get="/api/v1/course/friend")
}

## ----------------------------------------------------------------------------
## launch_screen 开屏页
## ----------------------------------------------------------------------------
struct CreateImageRequest {
    1: required i64 pic_type,
    2: optional i64 duration,
    3: string href,
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

service LaunchScreenService {
    // 创建一张开屏页
    CreateImageResponse CreateImage(1: CreateImageRequest req) (api.post="/api/v1/launch-screen/image"),
    // 获取开屏页
    GetImageResponse GetImage(1: GetImageRequest req) (api.get="/api/v1/launch-screen/image"),
    // 更改指定开屏页属性
    ChangeImagePropertyResponse ChangeImageProperty(1: ChangeImagePropertyRequest req) (api.put="/api/v1/launch-screen/image/property"),
    // 更改指定开屏页
    ChangeImageResponse ChangeImage(1: ChangeImageRequest req) (api.put="/api/v1/launch-screen/image"),
    // 删除指定开屏页
    DeleteImageResponse DeleteImage(1: DeleteImageRequest req) (api.delete="/api/v1/launch-screen/image"),
    // （移动端）获取开屏页
    MobileGetImageResponse MobileGetImage(1: MobileGetImageRequest req) (api.get="/api/v1/launch-screen/screen"),
    // 添加图片展示时间
    AddImagePointTimeResponse AddImagePointTime(1: AddImagePointTimeRequest req) (api.get="/api/v1/launch-screen/image/point-time"),
}

## ----------------------------------------------------------------------------
## paper 历年卷
## ----------------------------------------------------------------------------
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

// 以下是旧版本兼容
struct ListDirFilesForAndroidRequest {
    1: required string path,
}

struct ListDirFilesForAndroidResponse {}

struct GetDownloadUrlForAndroidRequest {
    1: required string filepath,
}

struct GetDownloadUrlForAndroidResponse {}


service PaperService {
    // 罗列指定文件夹下的文件
    ListDirFilesResponse ListDirFiles(1: ListDirFilesRequest req) (api.get="/api/v1/paper/list"),
    // 获取指定文件下载地址
    GetDownloadUrlResponse GetDownloadUrl(1: GetDownloadUrlRequest req) (api.get="/api/v1/paper/download"),

    // 兼容安卓
    ListDirFilesForAndroidResponse ListDirFilesForAndroid(1: ListDirFilesForAndroidRequest req) (api.get="/api/v1/list")
    GetDownloadUrlForAndroidResponse GetDownloadUrlForAndroid(1: GetDownloadUrlForAndroidRequest req) (api.get="/api/v1/downloadUrl")
}

## ----------------------------------------------------------------------------
## academic 学业信息
## ----------------------------------------------------------------------------
struct GetScoresRequest {}

struct GetScoresResponse {
    1: required list<model.Score> scores
}

struct GetGPARequest {}

struct GetGPAResponse {
    1: required model.GPABean gpa
}

struct GetCreditRequest {}

struct GetCreditResponse {
    1: required list<model.Credit> major
}

struct GetUnifiedExamRequest {}

struct GetUnifiedExamResponse {
    1: required list<model.UnifiedExam> unifiedExam
}

struct GetCreditV2Request {
}

struct GetCreditV2Response {
    1: required model.BaseResp base
    2: optional model.CreditResponse credit
}

struct GetPlanRequest{
    1: required string id
    2: required string cookies
}

struct GetPlanResponse{
    1: model.BaseResp base
}

service AcademicService {
    // 获取课程成绩
    GetScoresResponse GetScores(1:GetScoresRequest req)(api.get="/api/v1/jwch/academic/scores")
    // 获取 GPA 信息
    GetGPAResponse GetGPA(1:GetGPARequest req)(api.get="/api/v1/jwch/academic/gpa")
    // 获取学分统计
    GetCreditResponse GetCredit(1:GetCreditRequest req)(api.get="/api/v1/jwch/academic/credit")
    // 获取联考成绩
    GetUnifiedExamResponse GetUnifiedExam(1:GetUnifiedExamRequest req)(api.get="/api/v1/jwch/academic/unified-exam")
    // 获取培养计划
    GetPlanResponse GetPlan(1:GetPlanRequest req)(api.get="/api/v1/jwch/academic/plan")
    // 获取学分统计 V2
    GetCreditV2Response GetCreditV2(1:GetCreditV2Request req)(api.get="/api/v2/jwch/academic/credit")
}

## ----------------------------------------------------------------------------
## version（原url，版本控制相关）
## ----------------------------------------------------------------------------
struct LoginRequest{
    1: required string password,
}

struct LoginResponse{
    1: model.BaseResp base,
}

struct UploadRequest{
    1: required string version,
    2: required string code,
    3: required string url,
    4: required string feature,
    5: required string type,
    6: required string password,
    7: required bool force,
}

struct UploadResponse{
    1: model.BaseResp base,

}

struct UploadParamsRequest{
    1: required string password,
}

struct UploadParamsResponse{
    1: optional model.BaseResp base,
    2: optional string policy,
    3: optional string authorization,
}

struct DownloadReleaseApkRequest{
}

struct DownloadReleaseApkResponse{
    1: optional binary file,
}

struct DownloadBetaApkRequest{
}

struct DownloadBetaApkResponse{
    1: optional binary file,
}

struct GetReleaseVersionRequest{
}

struct GetReleaseVersionResponse{
    1: optional model.BaseResp base,
    2: optional string code,
    3: optional string feature,
    4: optional string url,
    5: optional string version,
    6: optional bool force,

}

struct GetBetaVersionRequest{
}

struct GetBetaVersionResponse{
    1: optional model.BaseResp base,
    2: optional string code,
    3: optional string feature,
    4: optional string url,
    5: optional string version,
    6: optional bool force,

}

struct GetSettingRequest{
    1: optional string account,
    2: optional string version,
    3: optional bool beta,
    4: optional string phone,
    5: optional bool isLogin,
    6: optional string loginType,
}

struct GetSettingResponse{
    1: optional model.BaseResp base,
    2: binary cloud_setting,
}

struct GetTestRequest{
    1: optional string account,
    2: optional string version,
    3: optional bool beta,
    4: optional string phone,
    5: optional bool isLogin,
    6: optional string loginType,
    7: optional string setting,
}

struct GetTestResponse{
    1: model.BaseResp base,
    2: binary cloud_setting,
}

struct GetCloudRequest{
}

struct GetCloudResponse{
    1: optional model.BaseResp base,
    2: binary cloud_setting,
}

struct SetCloudRequest{
    1: required string password,
    2: required string setting,
}

struct SetCloudResponse{
    1: model.BaseResp base,
}

struct GetDumpRequest{
}

struct GetDumpResponse{
    1: optional model.BaseResp base,
    2: string data,
}

struct AndroidGetVersioneRequest{
}

struct AndroidGetVersionResponse{
    1: optional model.BaseResp base,
    2: optional model.Version release,
    3: optional model.Version beta,
}

service VersionService{
    LoginResponse Login(1:LoginRequest req)(api.post="/api/v2/url/login")
    UploadResponse UploadVersion(1:UploadRequest req)(api.post="/api/v2/url/upload")
    UploadParamsResponse UploadParams(1:UploadParamsRequest req)(api.post="/api/v2/url/upload-params")
    DownloadReleaseApkResponse DownloadReleaseApk(1:DownloadReleaseApkRequest req)(api.get="/api/v2/url/release.apk")
    DownloadBetaApkResponse DownloadBetaApk(1:DownloadBetaApkRequest req)(api.get="/api/v2/url/beta.apk")
    GetReleaseVersionResponse GetReleaseVersion(1:GetReleaseVersionRequest req)(api.get="/api/v2/url/version.json")
    GetBetaVersionResponse GetBetaVersion(1:GetBetaVersionRequest req)(api.get="/api/v2/url/versionbeta.json")
    GetSettingResponse GetSetting(1:GetSettingRequest req)(api.get="/api/v2/url/settings.php")
    GetTestResponse GetTest(1:GetTestRequest req)(api.post="/api/v2/url/test")
    GetCloudResponse GetCloud(1:GetCloudRequest req)(api.get="/api/v2/url/getcloud")
    SetCloudResponse SetCloud(1:SetCloudRequest req)(api.post="/api/v2/url/setcloud")
    GetDumpResponse GetDump(1:GetDumpRequest req)(api.get="/api/v2/url/dump")
    AndroidGetVersionResponse AndroidGetVersion(1:AndroidGetVersioneRequest req)(api.get="/api/v2/version/android"),

}

## ----------------------------------------------------------------------------
## common（通用内容，如隐私政策等信息）
## ----------------------------------------------------------------------------
struct GetCSSRequest{
}

struct GetCSSResponse{
    1: binary css,
}

struct GetHtmlRequest{
}

struct GetHtmlResponse{
    1: binary html,
}

struct GetUserAgreementRequest{
}

struct GetUserAgreementResponse{
    1: binary user_agreement,
}

// 学期列表
struct TermListRequest {}

struct TermListResponse {
    1: required model.BaseResp base
    2: required model.TermList term_lists
}

// 学期信息
struct TermRequest {
    1: required string term
}

struct TermResponse {
    1: required model.BaseResp base
    2: required model.TermInfo term_info
}

struct GetNoticeRequst {
    1: required i64 pageNum
}

struct GetNoticeResponse {
    1: required list<model.NoticeInfo> notices
    2: required i64 total
}

struct GetContributorInfoRequest {
}

struct GetContributorInfoResponse {
    1: required list<model.Contributor> fzuhelper_app
    2: required list<model.Contributor> fzuhelper_server
    3: required list<model.Contributor> jwch
    4: required list<model.Contributor> yjsy
}

struct GetToolboxConfigRequest {
    1: optional i64 version
    2: optional string student_id
    3: optional string platform
}

struct GetToolboxConfigResponse {
    1: required list<model.ToolboxConfig> config
}

struct PutToolboxConfigRequest {
    1: required string secret
    2: required i64 tool_id
    3: optional string student_id
    4: optional string platform
    5: optional i64 version
    6: optional bool visible
    7: optional string name
    8: optional string icon
    9: optional string type
    10: optional string message
    11: optional string extra
}

struct PutToolboxConfigResponse {
    1: optional i64 config_id
}

service CommonService {
    // （兼容）获取隐私政策 css
    GetCSSResponse GetCSS(1:GetCSSRequest req)(api.get="/api/v2/common/fzu-helper.css"),
    // （兼容）获取隐私政策 html
    GetHtmlResponse GetHtml(1:GetHtmlRequest req)(api.get="/api/v2/common/fzu-helper.html"),
    // 获取用户协议
    GetUserAgreementResponse GetUserAgreement(1: GetUserAgreementRequest req) (api.get="/api/v2/common/user-agreement.html")
    // 学期信息：学期列表
    TermListResponse GetTermsList(1: TermListRequest req) (api.get="/api/v1/terms/list")
    // 学期信息：学期详情
    TermResponse GetTerm(1: TermRequest req) (api.get="/api/v1/terms/info")
    // 获取教务处通知
    GetNoticeResponse GetNotice(1: GetNoticeRequst req) (api.get="/api/v1/common/notice")
    // 获取贡献者列表
    GetContributorInfoResponse GetContributorInfo(1: GetContributorInfoRequest req)(api.get="/api/v1/common/contributor")
     // 获取工具箱配置
    GetToolboxConfigResponse GetToolboxConfig(1:GetToolboxConfigRequest req)( api.get="/api/v1/toolbox/config")
    // 更新工具箱配置
    PutToolboxConfigResponse PutToolboxConfig(1:PutToolboxConfigRequest req)(api.put="/api/v1/toolbox/config")
}

## ----------------------------------------------------------------------------
## oa（目前只有feedback）
## ----------------------------------------------------------------------------
struct CreateFeedbackRequest {
    1: required string stu_id,
    2: required string name,
    3: required string college,
    4: required string contact_phone,
    5: required string contact_qq,
    6: required string contact_email,

    7:  required string network_env,    // "2G"/"3G"/"4G"/"5G"/"wifi"/"unknown"
    8:  required bool   is_on_campus,    // true/false
    9:  required string os_name,
    10: required string os_version,
    11: required string manufacturer,
    12: required string device_model,

    13: required string problem_desc,

    14: required string screenshots,     // JSON 字符串文本，如 "[]"
    15: required string app_version,
    16: required string version_history,  // JSON，建议 "[]"

    17: required string network_traces,   // JSON，允许对象或数组，建议 "[]"
    18: required string events,          // JSON，建议 "[]"
    19: required string user_settings     // JSON，建议 "{}"
}

struct CreateFeedbackResponse {
    1: required model.BaseResp base,
    2: required i64 report_id
}

struct GetFeedbackByIDRequest{
    1: required i64   report_id,
}

struct FeedbackDetailResponse {
    1: required model.BaseResp base,
    2: optional model.Feedback data,
}

struct GetListFeedbackRequest{
    1: optional string stu_id,
    2: optional string name,

    3: optional string network_env,    // "2G"/"3G"/"4G"/"5G"/"wifi"/"unknown"
    4: optional bool   is_on_campus,    // true/false
    5: optional string os_name,
    6: optional string problem_desc,
    7: optional string app_version,
    8: optional i64    begin_time_ms
    9: optional i64    end_time_ms

    10: optional i64 limit
    11: optional i64 page_token
    12: optional bool order_desc
}

struct GetListFeedbackResponse{
    1: required model.BaseResp base,
    2: optional list<model.FeedbackListItem> data,
    3: optional i64 page_token
}

service FeedbackService {
    CreateFeedbackResponse CreateFeedback(1: CreateFeedbackRequest request)
        (api.post="/api/v1/feedback/create", api.body="request");
    FeedbackDetailResponse GetFeedbackByID(1: GetFeedbackByIDRequest request)
        (api.get="/api/v1/feedbacks/detail");
    GetListFeedbackResponse ListFeedback(1: GetListFeedbackRequest request)
      (api.get="/api/v1/feedbacks/list");
}
