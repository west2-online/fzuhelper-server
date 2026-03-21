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

service CourseService {
    CourseListResponse GetCourseList(1: CourseListRequest req)
    TermListResponse GetTermList(1: TermListRequest req)
    GetCalendarResponse GetCalendar(1: GetCalendarRequest req)
    GetLocateDateResponse GetLocateDate(1: GetLocateDateRequest req)
    GetFriendCourseResponse GetFriendCourse(1: GetFriendCourseRequest req)
    GetAutoAdjustCourseListResponse GetAutoAdjustCourseList(1: GetAutoAdjustCourseListRequest req)
    UpdateAdjustCourseResponse UpdateAdjustCourse(1: UpdateAdjustCourseRequest req)
}
