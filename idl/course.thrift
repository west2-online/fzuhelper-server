namespace go course
include "model.thrift"

struct CourseListRequest {
    1: required model.LoginData loginData
    2: required string term
}

struct CourseListResponse {
    1: required model.BaseResp base
    2: required list<model.Course> data
}

service CourseService {
    CourseListResponse GetCourseList(1: CourseListRequest req)
}
