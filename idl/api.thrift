namespace go api


// 获取用户信息
struct UserInfoRequest {

}

struct UserInfoResp {
    1: required string code,
    2: required string message,
    3: required UserInfoData data,
    
}

struct UserInfoData {
    1: required string id,
    2: required string sex,
    3: required string birthday,
    4: required string phone,
    5: required i64 grade,
    6: required string college,
    7: required string major,
    8: required string statusChange,
}


// 获取验证码结果
struct ValidateCodeRequest {
    1: optional string image,
}

struct ValidateCodeResp {
    1: required string code,
    2: required string message,
    3: required string data,
}


// 修改密码
struct ChangePasswordRequest {

    1: required string original,
    2: required string new,
}

struct ChangePasswordResp {
    1: required string code,
    2: required string message,
}

// 校历
struct SchoolCalendarRequest {
    1: optional string term,
}

struct SchoolCalendarResponse {
    1: required string code,
    2: required string message,
    3: required list<SchoolCalendarData> data,
}

struct SchoolCalendarData {
    1: required string dateBegin,
    2: required string dateEnd,
    3: required string name,
}


// 空教室
struct EmptyRoomRequest {
    1: required string date,
    2: required string campus,
    3: required string startTime,
    4: required string endTime,
}

struct EmptyRoomResponse{
    1: required string code,
    2: required string message,
    3: required list<EmptyRoomData> data,
}

struct EmptyRoomData{
    1: required string build,
    2: required string location,
    3: required string capacity,
    4: required string type,
}


// 考场查询
struct ExamRequest {
    1: required string term,
}

struct ExamResp {
    1: required string code,
    2: required string message,
    3: required list<ExamData> data,
}

struct ExamData{
    1: required string name,
    2: required string credit,
    3: required string teacher,
    4: required string location,
    5: required string date,
    6: required string time,
}

// 获取课程
struct CourseListRequest {
    1: optional string term,

}

struct CourseListResponse {
    1: required string code,
    2: required string message,
    3: required list<CourseListData> data,
}

struct CourseListData {
    1: required string name,
    2: required string location,
    3: required string startTime,
    4: required string endTime,
    5: required string startWeek,
    6: required string endWeek,
    7: required string double,
    8: required string single,
    9: required string weekday,
    10: required string year,
    11: required string term,
    12: required string note,
    13: required string plan,
    14: required string syllabus,
    15: required string teacher,
}

// 获取成绩详情
struct ScoresRequest {


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


service UserService {
    UserInfoResp GetUserInfo(1: UserInfoRequest req) (api.get = "/user/info"),
    ValidateCodeResp ValidateCode(1: ValidateCodeRequest req) (api.post = "user/validateCode"),
    ChangePasswordResp ChangePassword(1: ChangePasswordRequest req) (api.put = "/user/info"),
    SchoolCalendarResponse GetSchoolCalendar(1: SchoolCalendarRequest req) (api.get = "/user/schoolCalendar"),
}

service EmptyRoomService{
    EmptyRoomResponse GetEmptyRoom(1:EmptyRoomRequest req) (api.get = "/classroom/empty"),
    ExamResp GetExam(1:ExamRequest req) (api.get = "/classroom/exam"),
}

service CourseService {
    CourseListResponse GetCourseList(1: CourseListRequest req) (api.get = "/course/list")
}

service academicService {
    ScoresResponse GetScores(1: ScoresRequest req) (api.get = "/academic/scores"),
    GPAResp GetGPA(1: GPARequest req) (api.get = "/academic/gpa"),
    CreditResp GetCredit(1: CreditRequest req) (api.get = "/academic/credit"),
    UnifiedExamResp GetUnifiedExam(1: UnifiedExamRequest req) (api.get = "/academic/unifiedExam"),
    PlanResp GetPlan(1: PlanRequest req) (api.get = "/academic/plan"),
}