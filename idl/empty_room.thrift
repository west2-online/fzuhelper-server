namespace go empty_room



// 空教室
struct EmptyRoomRequest {
    1: required string date,
    2: required string campus,
    3: required string startTime,
    4: required string endTime,
    5: required string id,
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
    2: required string id,
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

service EmptyRoomService{
    EmptyRoomResponse GetEmptyRoom(1:EmptyRoomRequest req),
    ExamResp GetExam(1:ExamRequest req)
}