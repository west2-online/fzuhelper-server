namespace go version
include"model.thrift"

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
    1: model.BaseResp base,
    2: optional string policy,
    3: optional string authorization,
}

struct DownloadReleaseApkRequest{
}

struct DownloadReleaseApkResponse{
    1: model.BaseResp base,
    2: string redirect_url,
}

struct DownloadBetaApkRequest{
}

struct DownloadBetaApkResponse{
    1: model.BaseResp base,
    2: string redirect_url,
}

struct GetReleaseVersionRequest{
}

struct GetReleaseVersionResponse{
    1: model.BaseResp base,
    2: optional string code,
    3: optional string feature,
    4: optional string url,
    5: optional string version,
    6: optional bool force,
}

struct GetBetaVersionRequest{
}

struct GetBetaVersionResponse{
    1: model.BaseResp base,
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
    1: model.BaseResp base,
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
    1: model.BaseResp base,
    2: string data,
}

struct AndroidGetVersioneRequest{
}

struct AndroidGetVersionResponse{
    1: model.BaseResp base,
    2: optional model.Version release,
    3: optional model.Version beta,
}

service VersionService{
    LoginResponse Login(1:LoginRequest req)(api.post="/api/v1/url/login"),
    UploadResponse UploadVersion(1:UploadRequest req)(api.post="/api/v1/url/api/upload"),
    UploadParamsResponse UploadParams(1:UploadParamsRequest req)(api.post="/api/v1/url/api/uploadparams"),
    DownloadReleaseApkResponse DownloadReleaseApk(1:DownloadReleaseApkRequest req)(api.get="/api/v1/url/release.apk"),
    DownloadBetaApkResponse DownloadBetaApk(1:DownloadBetaApkRequest req)(api.get="/api/v1/url/beta.apk"),
    GetReleaseVersionResponse GetReleaseVersion(1:GetReleaseVersionRequest req)(api.get="/api/v1/url/version.json"),
    GetBetaVersionResponse GetBetaVersion(1:GetBetaVersionRequest req)(api.get="/api/v1/url/versionbeta.json"),
    GetSettingResponse GetSetting(1:GetSettingRequest req)(api.get="/api/v1/url/settings.php"),
    GetTestResponse GetTest(1:GetTestRequest req)(api.post="/api/v1/url/test"),
    GetCloudResponse GetCloud(1:GetCloudRequest req)(api.get="/api/v1/url/getcloud"),
    SetCloudResponse SetCloud(1:SetCloudRequest req)(api.post="/api/v1/url/setcloud"),
    GetDumpResponse GetDump(1:GetDumpRequest req)(api.get="/api/v1/url/dump"),
    AndroidGetVersionResponse AndroidGetVersion(1:AndroidGetVersioneRequest req),

}

