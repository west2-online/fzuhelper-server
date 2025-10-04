namespace go academic
include "model.thrift"

struct GetScoresRequest {
}

struct GetScoresResponse {
    1: required model.BaseResp base
    2: optional list<model.Score> scores
}

struct GetGPARequest {
}

struct GetGPAResponse {
    1: required model.BaseResp base
    2: optional model.GPABean gpa
}

struct GetCreditRequest {
}

struct GetCreditResponse {
    1: required model.BaseResp base
    2: optional list<model.Credit> major
}

struct GetUnifiedExamRequest {
}

struct GetUnifiedExamResponse {
    1: required model.BaseResp base
    2: optional list<model.UnifiedExam> unifiedExam
}

struct GetPlanRequest{
}

struct GetPlanResponse{
    1: required model.BaseResp base
    2: string url,
}

struct GetCreditV2Request {
}

struct GetCreditV2Response {
    1: required model.BaseResp base
    2: optional model.CreditResponse credit
}

service AcademicService {
    GetScoresResponse GetScores(1:GetScoresRequest req)
    GetGPAResponse GetGPA(1:GetGPARequest req)
    GetCreditResponse GetCredit(1:GetCreditRequest req)
    GetUnifiedExamResponse GetUnifiedExam(1:GetUnifiedExamRequest req)
    GetPlanResponse GetPlan(1:GetPlanRequest req)
    GetCreditV2Response GetCreditV2(1:GetCreditV2Request req)
}
