namespace go oa

include "model.thrift"

// for backend testing
struct CreateFeedbackRequest {
    1: required i64   report_id,
    2: required string stu_id,
    3: required string name,
    4: required string college,
    5: required string contact_phone,
    6: required string contact_qq,
    7: required string contact_email,

    8:  required string network_env,    // "2G"/"3G"/"4G"/"5G"/"wifi"/"unknown"
    9:  required bool   is_on_campus,    // true/false
    10: required string os_name,
    11: required string os_version,
    12: required string manufacturer,
    13: required string device_model,

    14: required string problem_desc,

    15: required string screenshots,     // JSON 字符串文本，如 "[]"
    16: required string app_version,
    17: required string version_history,  // JSON，建议 "[]"

    18: required string network_traces,   // JSON，允许对象或数组，建议 "[]"
    19: required string events,          // JSON，建议 "[]"
    20: required string user_settings     // JSON，建议 "{}"
}

struct CreateFeedbackResponse {
    1: required model.BaseResp base,
}

struct GetFeedbackRequest{
    1: required i64   report_id,
}

struct GetFeedbackResponse {
    1: required model.BaseResp base,
    2: optional model.Feedback data,
}

service OAService {
    CreateFeedbackResponse CreateFeedback(1: CreateFeedbackRequest request);
    GetFeedbackResponse GetFeedback(1: GetFeedbackRequest request);
}
