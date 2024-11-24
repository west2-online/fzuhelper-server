namespace go common
include "model.thrift"

// 学期列表
struct TermListRequest {
}

struct TermListResponse {
    1: required model.BaseResp base
    2: required string current_term
    3: required list<model.Term> terms
}

// 学期信息
struct TermRequest {
    1: required string term
}

struct TermResponse {
    1: required model.BaseResp base
    2: required string term_id
    3: required string term
    4: required string school_year
    5: required list<model.TermEvent> events
}

service CommonService {
    // 学期信息：学期列表
    TermListResponse GetTermsList(1: TermListRequest req)
    // 学期信息：学期详情
    TermResponse GetTerm(1: TermRequest req)
}
