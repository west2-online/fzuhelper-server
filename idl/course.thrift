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
    1: required string term
}

struct GetCalendarResponse {
    1: required model.BaseResp base
    2: required string content
}

struct GetLocateDateRequest{}

struct GetLocateDateResponse{
    1: required model.BaseResp base
    2: optional model.LocateDate locateDate
}

service CourseService {
    CourseListResponse GetCourseList(1: CourseListRequest req)
    TermListResponse GetTermList(1: TermListRequest req)
    GetCalendarResponse GetCalendar(1: GetCalendarRequest req)
    GetLocateDateResponse GetLocateDate(1:GetLocateDateRequest req)
}
