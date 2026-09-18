namespace go oa

include "model.thrift"

// for backend testing
struct CreateFeedbackRequest {
    1: required string stu_id,
    2: required string name,
    3: required string college,
    // 联系方式至少填写一项；由 service 校验
    4: optional string contact_phone,
    5: optional string contact_qq,
    6: optional string contact_email,

    7:  required string network_env,    // "2G"/"3G"/"4G"/"5G"/"wifi"/"unknown"
    8:  required bool   is_on_campus,    // true/false
    9:  required string os_name,
    10: required string os_version,
    11: required string manufacturer,
    12: required string device_model,

    13: required string problem_desc,

    14: required string screenshots,     // JSON URL 数组，最多 9 项，如 "[]"
    15: required string app_version,
    16: required string version_history,  // JSON，建议 "[]"

    17: required string network_traces,   // JSON 对象，缺省值 "{}"
    18: required string events,           // JSON 数组，缺省值 "[]"
    19: required string user_settings     // JSON 对象，缺省值 "{}"
}

struct CreateFeedbackResponse {
    1: required model.BaseResp base,
    2: required i64 report_id
}

struct UploadFeedbackScreenshotRequest {
    1: required binary file
}

struct UploadFeedbackScreenshotResponse {
    1: required model.BaseResp base,
    2: required string url
}

struct UploadFeedbackLogRequest {
    1: required binary file
}

struct UploadFeedbackLogResponse {
    1: required model.BaseResp base,
    2: required string url
}

struct GetFeedbackByIDRequest {
    1: required i64   report_id,
    2: required string stu_id
}

struct GetFeedbackByIDResponse {
    1: required model.BaseResp base,
    2: optional model.Feedback data,
}


struct GetListFeedbackRequest {
    1: required string stu_id,

    // fields 2-9 原管理员筛选条件，不再使用
    10: optional i64 limit
    11: optional i64 page_token
    12: optional bool order_desc
}

struct GetListFeedbackResponse {
    1: required model.BaseResp base,
    2: optional list<model.FeedbackListItem> data,
    3: optional i64 page_token
}


service OAService {
    CreateFeedbackResponse CreateFeedback(
        1: CreateFeedbackRequest request
    )

    UploadFeedbackScreenshotResponse UploadFeedbackScreenshot(
        1: UploadFeedbackScreenshotRequest request
    )

    UploadFeedbackLogResponse UploadFeedbackLog(
        1: UploadFeedbackLogRequest request
    )

    GetFeedbackByIDResponse GetFeedbackById(
        1: GetFeedbackByIDRequest request
    )

    GetListFeedbackResponse GetFeedbackList(
        1: GetListFeedbackRequest request
    )
}
