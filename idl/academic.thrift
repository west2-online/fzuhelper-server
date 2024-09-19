namespace go academic

// 获取成绩详情
struct ScoresRequest {
    1: required string id,
}

struct ScoresResponse {
    1: required string code,
    2: required string message,
    3: required list<ScoresData> data,
}

struct ScoresData {
    1: required string name,
    2: required string gpa,
    3: required string credit,
    4: required string score,
    5: required string teacher,
    6: required string year,
    7: required string term,
}


// 获取绩点排名
struct GPARequest {
    1: required string id,
}

struct GPAResp {
    1: required string code,
    2: required string message,
    3: required GPAData data,
}

struct GPAData {
    1: required string time,
    2: required list<GPAArrayData> data,
}

struct GPAArrayData {
    1: required string type,
    2: required string value,
}

// 获取学分统计
struct CreditRequest {
    1: required string id,
}

struct CreditResp {
    1: required string code,
    2: required string message,
    3: required list<CreditData> data,
}

struct CreditData {
    1: required string type,
    2: required CreditArrayData data,

}

struct CreditArrayData {
    1: required string type,
    2: required string gain,
    3: required string total,
}


// 获取统考成绩
struct UnifiedExamRequest {
    1: required string id,
}

struct UnifiedExamResp {
    1: required string code,
    2: required string message,
    3: required list<UnifiedExamData> data,
}

struct UnifiedExamData {
    1: required string name,
    2: required string score,
    3: required string term,
}


// 获取专业培养计划
struct PlanRequest {

}

struct PlanResp {
    1: required string code,
    2: required string message,
    3: required string data,
}



service academicService {
    ScoresResponse GetScores(1: ScoresRequest req),
    GPAResp GetGPA(1: GPARequest req),
    CreditResp GetCredit(1: CreditRequest req),
    UnifiedExamResp GetUnifiedExam(1: UnifiedExamRequest req),
    PlanResp GetPlan(1: PlanRequest req),
}