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
}
