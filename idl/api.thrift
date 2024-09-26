namespace go api

//由前端给的登陆信息，包括id和cookies
struct LoginData {
    1: required string id
    2: required list<string> cookies
}

struct Classroom {
    1: required string build
    2: required string location
    3: required string capacity
    4: required string type
}

struct EmptyClassroomRequest {
    1: required string date
    2: required string campus
    3: required string startTime;//节数
    4: required string endTime;
}

struct EmptyClassroomResponse {
    1: required list<Classroom> classrooms
}


service ClassRoomService {
    EmptyClassroomResponse GetEmptyClassrooms(1: EmptyClassroomRequest request)(api.get="/api/v1/common/classroom/empty")
}
