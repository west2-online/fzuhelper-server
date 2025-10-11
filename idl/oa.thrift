namespace go oa

include "model.thrift"

// for backend testing
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


service OAService {
    CreateFeedbackResponse CreateFeedback(1: CreateFeedbackRequest request);
    FeedbackDetailResponse GetFeedbackById(1: GetFeedbackByIDRequest request);
    GetListFeedbackResponse GetFeedbackList(1: GetListFeedbackRequest request);
}
