namespace go course
include "model.thrift"

struct TermListRequest {}

struct TermListResponse {
    1: required model.BaseResp base
    2: required list<string> data
}

struct CourseListRequest {
    1: required string term
    2: optional bool isRefresh
}

struct CourseListResponse {
    1: required model.BaseResp base
    2: required list<model.Course> data
    3: optional list<CustomCourseItem> customCourses
}

struct GetCalendarRequest {
    1: required string stu_id
}

struct GetCalendarResponse {
    1: required model.BaseResp base
    2: binary ics
}

struct GetLocateDateRequest{}

struct GetLocateDateResponse{
    1: required model.BaseResp base
    2: optional model.LocateDate locateDate
}
struct GetFriendCourseRequest {
    1: required string term
    2: required string id
}

struct GetFriendCourseResponse {
    1: required model.BaseResp base
    2: required list<model.Course> data
}

struct GetAutoAdjustCourseListRequest {
    1: required string term
}

struct GetAutoAdjustCourseListResponse {
    1: required model.BaseResp base
    2: required list<model.AdjustCourse> data
}

struct UpdateAdjustCourseRequest {
    1: required i64 id
    2: required string secret
    3: optional bool enabled
    4: optional string from_date
    5: optional string to_date
}

struct UpdateAdjustCourseResponse {
    1: required model.BaseResp base
}

// ====== 自定义课程 ======

// 自定义课程项（用于传输）
struct CustomCourseItem {
    1: optional string id
    2: required string name
    3: optional string teacher
    4: required string location
    5: required i32 startClass
    6: required i32 endClass
    7: required i32 startWeek
    8: required i32 endWeek
    9: required i32 weekday
    10: optional bool single
    11: optional bool double_
    12: optional string color
    13: optional string remark
}

// Upsert 自定义课程请求（新增或更新）
struct UpsertCustomCourseRequest {
    1: required string term
    2: required CustomCourseItem course
}

// Upsert 自定义课程响应
struct UpsertCustomCourseResponse {
    1: required model.BaseResp base
    2: optional string courseId
}

// 删除自定义课程请求
struct DeleteCustomCourseRequest {
    1: required string term
    2: required string courseId
}

// 删除自定义课程响应
struct DeleteCustomCourseResponse {
    1: required model.BaseResp base
}

// ====== END 自定义课程 ======

service CourseService {
    CourseListResponse GetCourseList(1: CourseListRequest req)
    TermListResponse GetTermList(1: TermListRequest req)
    GetCalendarResponse GetCalendar(1: GetCalendarRequest req)
    GetLocateDateResponse GetLocateDate(1: GetLocateDateRequest req)
    GetFriendCourseResponse GetFriendCourse(1: GetFriendCourseRequest req)
    GetAutoAdjustCourseListResponse GetAutoAdjustCourseList(1: GetAutoAdjustCourseListRequest req)
    UpdateAdjustCourseResponse UpdateAdjustCourse(1: UpdateAdjustCourseRequest req)
    
    // 自定义课程接口
    UpsertCustomCourseResponse UpsertCustomCourse(1: UpsertCustomCourseRequest req)
    DeleteCustomCourseResponse DeleteCustomCourse(1: DeleteCustomCourseRequest req)
}
