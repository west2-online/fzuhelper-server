namespace go template

include "model.thrift"

struct PingRequest {
    1: optional string text,
}

struct PingResponse {
    1: model.BaseResp base,
    2: string pong,
}

service TemplateService {
    PingResponse Ping(1: PingRequest req),
}
