namespace go course

// 获取课程
struct CourseListRequest {
    1: optional string term,
    2: required string id,

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

service CourseService {
    CourseListResponse GetCourseList(1: CourseListRequest req)
}