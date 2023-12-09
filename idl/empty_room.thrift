namespace go empty_room

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct EmptyRoomRequest{
    1: required string token,
    2: required string time,
    3: required string start,
    4: required string end,
    5: required string building,
    6: optional string account,
    7: optional string password,
}

struct EmptyRoomResponse{
    1: required BaseResp base,
    2: required list<string> room_name,
}

service EmptyRoomService{
    EmptyRoomResponse GetEmptyRoom(1:EmptyRoomRequest req),
}