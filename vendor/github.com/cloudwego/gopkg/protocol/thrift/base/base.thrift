namespace py base
namespace go base
namespace java com.bytedance.thrift.base

struct Base {
    1: string LogID = "",
    2: string Caller = "",
    3: string Addr = "",
    6: optional map<string, string> Extra,
}

struct BaseResp {
    1: string StatusMessage = "",
    2: i32 StatusCode = 0,
    3: optional map<string, string> Extra,
}
