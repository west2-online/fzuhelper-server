namespace go classroom
struct BaseResp {
    1: i64 code,
    2: string msg,
}



struct Classroom {
    1: required string build
    2: required string location
    3: required string capacity
    4: required string type
}

struct LoginData{
    1: required string id
    2: required list<string> cookies
}

struct EmptyRoomRequest{
    1: required string date
    2: required string campus
    3: required string startTime;//节数
    4: required string endTime;
}

struct EmptyRoomResponse{
    1: required BaseResp base,
    2: required list<Classroom> rooms,
}

service ClassroomService {
    EmptyRoomResponse GetEmptyRoom(1:EmptyRoomRequest req),
}
