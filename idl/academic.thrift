namespace go academic
include "model.thrift"

struct GetScoresRequest {
    1: required string id
    2: required list<string> cookies
}

struct GetScoresResponse {
    1: required model.BaseResp base
    2: required list<model.Score> scores
}

struct GetGPARequest {
    1: required string id
    2: required list<string> cookies
}

struct GetGPAResponse {
    1: required model.BaseResp base
    2: required model.GPABean gpa
}

struct GetCreditRequest {
    1: required string id
    2: required list<string> cookies
}

struct GetCreditResponse {
    1: required model.BaseResp base
    2: required list<model.Credit> major
}

struct GetUnifiedExamRequest {
    1: required string id
    2: required list<string> cookies
}

struct GetUnifiedExamResponse {
    1: required model.BaseResp base
    2: required list<model.UnifiedExam> unifiedExam
}

service AcademicService {
    GetScoresResponse GetScores(1:GetScoresRequest req)
    GetGPAResponse GetGPA(1:GetGPARequest req)
    GetCreditResponse GetCredit(1:GetCreditRequest req)
    GetUnifiedExamResponse GetUnifiedExam(1:GetUnifiedExamRequest req)
}
