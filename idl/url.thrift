namespace go url
include"model.thrift"

struct LoginRequest{
    1: optional string password,
}

struct LoginResponse{
    1: model.BaseResp base,
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
    1: model.BaseResp base,

}

struct UploadParamsRequest{
    1: optional string password,
}

struct UploadParamsResponse{
    1: model.BaseResp base,
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
    1: model.BaseResp base,
    2: optional string feature,
    3: optional string url,
    4: optional string version,
}

struct GetBetaVersionRequest{
}

struct GetBetaVersionResponse{
    1: model.BaseResp base,
    2: optional string feature,
    3: optional string url,
    4: optional string version,
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
    1: model.BaseResp base,
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
    1: model.BaseResp base,
    //todo:补全字段
}

struct GetCloudRequest{
}

struct GetCloudResponse{
    1: model.BaseResp base,
    2: string data,
}

struct SetCloudRequest{
    1: optional string password,
    2: optional string setting,
}

struct SetCloudResponse{
    1: model.BaseResp base,
}

struct GetDumpRequest{
}

struct GetDumpResponse{
    1: model.BaseResp base,
    2: list<i64> data_list,
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
    UploadResponse UploadVersion(1:UploadRequest req)(api.post="/api/v1/url/api/upload"),
    UploadResponse UploadParams(1:UploadParamsRequest req)(api.post="/api/v1/url/api/uploadparams"),
    DownloadReleaseApkResponse DownloadReleaseApk(1:DownloadReleaseApkRequest req)(api.get="/api/v1/url/release.apk"),
    DownloadBetaApkResponse DownloadBetaApk(1:DownloadBetaApkRequest req)(api.get="/api/v1/url/beta.apk"),
    GetReleaseVersionResponse GetReleaseVersion(1:GetReleaseVersionRequest req)(api.get="/api/v1/url/version.json"),
    GetBetaVersionResponse GetBetaVersion(1:GetBetaVersionRequest req)(api.get="/api/v1/url/versionbeta.json"),
    GetSettingResponse GetSetting(1:GetSettingRequest req)(api.get="/api/v1/url/settings.php"),
    GetTestResponse GetTest(1:GetSettingRequest req)(api.post="/api/v1/url/test"),
    GetCloudResponse GetCloud(1:GetCloudRequest req)(api.get="/api/v1/url/getcloud"),
    SetCloudResponse SetCloud(1:SetCloudRequest req)(api.get="/api/v1/url/setcloud"),
    GetDumpResponse GetDump(1:GetDumpRequest req)(api.get="/api/v1/url/dump"),
    GetCSSResponse GetCSS(1:GetCSSRequest req)(api.get="/api/v1/url/onekey/FZUHelper.css"),
    GetHtmlResponse GetHtml(1:GetHtmlRequest req)(api.get="/api/v1/url/onekey/FZUHelper.html"),
    GetUserAgreementResponse GetUserAgreement(1: GetUserAgreementRequest req) (api.get="/api/v1/url/onekey/UserAgreement.html")
}

