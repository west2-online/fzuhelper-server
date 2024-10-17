namespace go model

struct BaseResp {
    1: i64 code,
    2: string msg,
}

//由前端给的登陆信息，包括id和cookies, 这个struct仅用于测试返回数据，因为登录实现在前端完成，不会在实际项目中使用
struct LoginData {
    1: required string id
    2: required list<string> cookies
}

//Classroom 前端想要返回的fields
struct Classroom {
    1: required string build
    2: required string location
    3: required string capacity
    4: required string type
}

// === Course ===
// CourseScheduleRule 课程安排，详见 apifox
struct CourseScheduleRule {
    1: required string location
    2: required i32 startClass
    3: required i32 endClass
    4: required i32 startWeek
    5: required i32 endWeek
    6: required i32 weekday
    7: required bool single
    8: required bool double
    9: required bool adjust
}

// Course 课程信息，详见 apifox
struct Course {
    1: required string name
    2: required string teacher
    3: required list<CourseScheduleRule> scheduleRules
    4: required string remark
    5: required string lessonplan
    6: required string syllabus
    7: required string rawScheduleRules
    8: required string rawAdjust
}

// === END Course ===
