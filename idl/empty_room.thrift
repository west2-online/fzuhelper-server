namespace go empty_room

// struct BaseResp {
//     1: i64 code,
//     2: string msg,
// }

// struct EmptyRoomRequest{
//     1: required string token,
//     2: required string time,
//     3: required string start,
//     4: required string end,
//     5: required string building,
//     6: optional string account,
//     7: optional string password,
//     8: optional string id,
//     9: required list<string>cookies,
// }


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

service EmptyRoomService{
    EmptyRoomResponse GetEmptyRoom(1:EmptyRoomRequest req),
    ExamResp GetExam(1:ExamRequest req)
}