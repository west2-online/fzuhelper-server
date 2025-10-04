namespace go oa

include "model.thrift"

// for backend testing
struct CreateFeedbackRequest {
    1: required i64   reportId,
    2: required string stuId,
    3: required string name,
    4: required string college,
    5: required string contactPhone,
    6: required string contactQQ,
    7: required string contactEmail,

    8:  required string networkEnv,    // "2G"/"3G"/"4G"/"5G"/"wifi"/"unknown"
    9:  required bool   isOnCampus,    // true/false
    10: required string osName,
    11: required string osVersion,
    12: required string manufacturer,
    13: required string deviceModel,

    14: required string problemDesc,

    15: required string screenshots,     // JSON 字符串文本，如 "[]"
    16: required string appVersion,
    17: required string versionHistory,  // JSON，建议 "[]"

    18: required string networkTraces,   // JSON，允许对象或数组，建议 "[]"
    19: required string events,          // JSON，建议 "[]"
    20: required string userSettings     // JSON，建议 "{}"
}

struct CreateFeedbackResponse {
    1: required model.BaseResp base,
}

struct GetFeedbackRequest{
    1: required i64   reportId,
}

struct GetFeedbackResponse {
    1: required model.BaseResp base,
    2: optional model.Feedback data,
}

service OaService {
    CreateFeedbackResponse CreateFeedback(1: CreateFeedbackRequest request);
    GetFeedbackResponse GetFeedback(1: GetFeedbackRequest request);
}
