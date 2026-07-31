namespace go common
include "model.thrift"

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

// 学期信息
struct TermRequest {
    1: required string term
}

struct TermResponse {
    1: required model.BaseResp base
    2: required model.TermInfo term_info
}

// 学期列表
struct TermListRequest {
}

struct TermListResponse {
    1: required model.BaseResp base
    2: required model.TermList term_lists
}

// 教务处教学通知
struct NoticeRequest {
    1: required i64 pageNum
}

struct NoticeResponse {
    1: required model.BaseResp base
    2: optional list<model.NoticeInfo> notices
    3: required i64 total
}

// 获取贡献者列表
struct GetContributorInfoRequest {
}

struct GetContributorInfoResponse {
    1: required model.BaseResp base
    2: required list<model.Contributor> fzuhelper_app
    3: required list<model.Contributor> fzuhelper_server
    4: required list<model.Contributor> jwch
    5: required list<model.Contributor> yjsy
}

struct GetToolboxConfigRequest {
    1: optional i64 version
    2: optional string student_id
    3: optional string platform
}

struct GetToolboxConfigResponse {
    1: required model.BaseResp base
    2: required list<model.ToolboxConfig> config
}

struct CreateToolboxConfigRequest {
    1: required string secret
    2: required i64 tool_id
    3: required bool visible
    4: required string name
    5: required string icon
    6: required string type
    7: required string message
    8: required string extra
    9: required string student_id
    10: required string platform
    11: required i64 version
}

struct CreateToolboxConfigResponse {
    1: required model.BaseResp base
    2: optional model.ToolboxConfigDetail config
}

struct ListToolboxConfigsRequest {
    1: required string secret
    2: optional i64 page_num
    3: optional i64 page_size
}

struct ListToolboxConfigsResponse {
    1: required model.BaseResp base
    2: required list<model.ToolboxConfigDetail> config
    3: required i64 total
}

struct GetToolboxConfigByIDRequest {
    1: required string secret
    2: required i64 config_id
}

struct GetToolboxConfigByIDResponse {
    1: required model.BaseResp base
    2: optional model.ToolboxConfigDetail config
}

struct UpdateToolboxConfigRequest {
    1: required string secret
    2: required i64 config_id
    3: required i64 tool_id
    4: required bool visible
    5: required string name
    6: required string icon
    7: required string type
    8: required string message
    9: required string extra
    10: required string student_id
    11: required string platform
    12: required i64 version
}

struct UpdateToolboxConfigResponse {
    1: required model.BaseResp base
    2: optional model.ToolboxConfigDetail config
}

struct DeleteToolboxConfigRequest {
    1: required string secret
    2: required i64 config_id
}

struct DeleteToolboxConfigResponse {
    1: required model.BaseResp base
}

struct TracePingRequest {
}

struct TracePingResponse {
    1: required model.BaseResp base
    2: required string message
}

// 获取查询地理位置所需的签名 URL 和 Headers
struct GetSignedLocationApiUrlRequest{
    1: required string location
}

struct GetSignedLocationApiUrlResponse{
    1: required model.BaseResp base
    2: required string signed_url   // 签好名的完整请求 URL
    3: required map<string, string> headers  // 客户端发请求时需携带的 Headers
}

service CommonService {
    GetCSSResponse GetCSS(1:GetCSSRequest req)(api.get="/api/v1/url/onekey/FZUHelper.css"),
    GetHtmlResponse GetHtml(1:GetHtmlRequest req)(api.get="/api/v1/url/onekey/FZUHelper.html"),
    GetUserAgreementResponse GetUserAgreement(1: GetUserAgreementRequest req) (api.get="/api/v1/url/onekey/UserAgreement.html")
    // 学期信息：学期列表
    TermListResponse GetTermsList(1: TermListRequest req)
    // 学期信息：学期详情
    TermResponse GetTerm(1: TermRequest req)
    // 教务处教学通知
    NoticeResponse GetNotices(1: NoticeRequest req)
    // 获取贡献者列表
    GetContributorInfoResponse GetContributorInfo(1: GetContributorInfoRequest req)
    // 获取工具箱配置
    GetToolboxConfigResponse GetToolboxConfig(1:GetToolboxConfigRequest req)
    // 创建工具箱配置
    CreateToolboxConfigResponse CreateToolboxConfig(1:CreateToolboxConfigRequest req)
    // 获取工具箱配置列表
    ListToolboxConfigsResponse ListToolboxConfigs(1:ListToolboxConfigsRequest req)
    // 按 ID 获取工具箱配置
    GetToolboxConfigByIDResponse GetToolboxConfigByID(1:GetToolboxConfigByIDRequest req)
    // 按 ID 更新工具箱配置
    UpdateToolboxConfigResponse UpdateToolboxConfig(1:UpdateToolboxConfigRequest req)
    // 按 ID 删除工具箱配置
    DeleteToolboxConfigResponse DeleteToolboxConfig(1:DeleteToolboxConfigRequest req)
    // 链路追踪探针
    TracePingResponse TracePing(1:TracePingRequest req)
    // 获取查询地理位置所需的签名 URL 和 Headers
    GetSignedLocationApiUrlResponse GetSignedLocationApiUrl(1:GetSignedLocationApiUrlRequest req)
}
