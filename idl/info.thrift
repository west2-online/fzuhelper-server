namespace go info

include "model.thrift"


struct HelloWorldRequest{

}

struct HelloWroldResponse{

}


struct APILoginRequest {
    1: required string password
}

struct APILoginResponse {
    1: required string status
}

struct APIUploadRequest {
    1: required string password
    2: required string type
    3: required string version
    4: required string code
    5: required string feature
    6: required string url
}

struct APIUploadResponse {
    1: required string status
}

service InfoService {
    HelloWroldResponse HelloWorld(1:HelloWorldRequest req)
    APILoginResponse APILogin(1:APILoginRequest req)
}