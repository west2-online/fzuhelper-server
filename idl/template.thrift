namespace go template

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct PingRequest {
    1: optional string text,
}

struct PingResponse {
    1: BaseResp base,
    2: string pong,
}

service TemplateService {
    PingResponse Ping(1: PingRequest req),
}

