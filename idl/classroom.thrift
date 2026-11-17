namespace go classroom
include "model.thrift"

struct EmptyRoomRequest{
    1: required string date
    2: required string campus
    3: required string startTime;
    4: required string endTime;
}

struct EmptyRoomV2Request {
    1: required string date
    2: required string campus
}

struct EmptyRoomResponse{
    1: required model.BaseResp base,
    2: optional list<model.Classroom> rooms,
}

struct ExamRoomInfoRequest {
    1: required string term
}

struct ExamRoomInfoResponse {
    1: required model.BaseResp base,
    2: optional list<model.ExamRoomInfo> rooms,
}

service ClassroomService {
    EmptyRoomResponse GetEmptyRoom(1:EmptyRoomRequest req),
    EmptyRoomResponse GetEmptyRoomV2(1:EmptyRoomV2Request req),
    ExamRoomInfoResponse GetExamRoomInfo(1:ExamRoomInfoRequest req),
}
