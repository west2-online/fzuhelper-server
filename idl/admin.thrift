namespace go admin

include "model.thrift"

// 生成飞书 OAuth 授权地址，并缓存 state 与登录完成后的返回地址。
struct LoginRequest {
    1: required string return_to,
}

struct LoginResponse {
    1: required model.BaseResp base,
    2: required string authorization_url,
}

// 处理飞书 OAuth 回调，校验 state、获取飞书 user_id 并验证管理员白名单。
struct CallbackRequest {
    1: required string state,
    2: required string code,
}

struct CallbackResponse {
    1: required model.BaseResp base,
    2: required string redirect_url,
}

// 消费一次性 ticket，并签发管理面板使用的 Admin Token。
struct ExchangeTicketRequest {
    1: required string ticket,
}

struct ExchangeTicketResponse {
    1: required model.BaseResp base,
    2: required string access_token,
}

service AdminService {
    LoginResponse Login(1: LoginRequest req),
    CallbackResponse Callback(1: CallbackRequest req),
    ExchangeTicketResponse ExchangeTicket(1: ExchangeTicketRequest req),
}
