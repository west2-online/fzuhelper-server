namespace go url_modify
include"model.thrift"

struct LoginRequest{
    1: optional string password,
}

struct LoginResponse{
    1: string msg,
}

struct UploadRequest{
    1: optional string version,
    2: optional string code,
    3: optional string url,
    4: optional string feature,
    5: optional string type,
    6: optional string password,
}

struct UploadResponse{
    1: string msg,
}

struct UploadParamsRequest{
    1: optional string password,
}

struct UploadParamsResponse{
    1: optional string msg,
    2: optional string policy,
    3: optional string authorization,
}

struct DownlowdReleaseApkRequest{
}

struct DownlowdReleaseApkResponse{
    1: binary file,
}

struct DownlowdBetaApkRequest{
}

struct DownlowdBetaApkResponse{
    1: binary file,
}

struct GetReleaseVersionRequest{
}

struct GetReleaseVersionResponse{
    1: string code,
    2: string feature,
    3: string url,
    4: string version,
}

struct GetBetaVersionRequest{
}

struct GetBetaVersionResponse{
    1: string code,
    2: string feature,
    3: string url,
    4: string version,
}

struct GetSettingRequest{
    1: optional string account,
    2: optional string version,
    3: optional string beta,
    4: optional string phone,
    5: optional string isLogin,
    6: optional string loginType,
}

struct GetSettingResponse{
    2: optional string accountMistakePlan2,
    3: optional string exceptionInterceptors,
    4: optional string notice,
    5: optional string trustAllCerts,
}

struct GetTestRequest{
    1: optional string account,
    2: optional string version,
    3: optional string beta,
    4: optional string phone,
    5: optional string isLogin,
    6: optional string loginType,
}

struct GetTestResponse{
    //todo:补全字段
}

struct GetCloudRequest{
}

struct GetCloudResponse{
    1: string code,
    2: optional string data,
    3: optional string msg,
}

struct SetCloudRequest{
    1: optional string password,
    2: optional string setting,
}

struct SetCloudResponse{
    1: string msg,
}

struct GetDumpRequest{
}

struct GetDumpResponse{
    1: string return_json,
}

struct GetCSSRequest{
}

struct GetCSSResponse{
    1: string css,
}

struct GetHtmlRequest{
}

struct GetHtmlResponse{
    1: string html,
}

struct GetUserAgreementRequest{
}

struct GetUserAgreementResponse{
    1: string user_agreement,
}

service UrlService{
    LoginResponse Login(1:LoginRequest req)(api.post="/api/v1/url/login"),
    UploadResponse Upload(1:UploadRequest req)(api.post="/api/v1/url/api/upload"),
    UploadResponse UploadParams(1:UploadParamsRequest req)(api.post="/api/v1/url/api/uploadparams"),
    DownlowdReleaseApkResponse DownlowdReleaseApk(1:DownlowdReleaseApkRequest req)(api.get="/api/v1/url/release.apk"),
    DownlowdBetaApkResponse DownlowdBetaApk(1:DownlowdBetaApkRequest req)(api.get="/api/v1/url/beta.apk"),
    GetReleaseVersionResponse GetReleaseVersion(1:GetReleaseVersionRequest req)(api.get="/api/v1/url/version.json"),
    GetBetaVersionResponse GetBetaVersion(1:GetBetaVersionRequest req)(api.get="/api/v1/url/versionbeta.json"),
    GetSettingResponse GetSetting(1:GetSettingRequest req)(api.get="/api/v1/url/settings.php"),
    GetTestResponse GetTest(1:GetSettingRequest req)(api.post="/api/v1/url/test"),
    GetCloudResponse GetCloud(1:GetCloudRequest req)(api.get="/api/v1/url/getcloud"),
    SetCloudResponse SetCloud(1:SetCloudRequest req)(api.get="/api/v1/url/setcloud"),
    GetDumpResponse GetDump(1:GetDumpRequest req)(api.get="/api/v1/url/dump"),
    GetCSSResponse GetCSS(1:GetCSSRequest req)(api.get="/api/v1/url/onekey/FZUHelper.css"),
    GetHtmlResponse GetHtml(1:GetHtmlRequest req)(api.get="/api/v1/url/onekey/FZUHelper.html"),
}
