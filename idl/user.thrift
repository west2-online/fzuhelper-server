namespace go user

include "model.thrift"

// for backend testing
struct GetLoginDataRequest {
    1: required string id
    2: required string password
}

struct GetLoginDataResponse {
    1: required model.BaseResp base,
    2: required string id
    3: required string cookies
}

struct GetUserInfoRequest{
}

struct GetUserInfoResponse{
    1: required model.BaseResp base,
    2: optional model.UserInfo data,
}

struct GetLoginDataForYJSYRequest{
    1: required string id
    2: required string password
}

struct GetLoginDataForYJSYResponse{
    1: required model.BaseResp base,
    2: required string id
    3: required string cookies
}
struct GetInvitationCodeRequest{
    1: optional bool isRefresh
}
struct GetInvitationCodeResponse{
    1: required model.BaseResp base,
    2: required string invitation_code,
}
struct BindInvitationRequest{
        1: required string invitation_code
}
struct BindInvitationResponse{
        1: required model.BaseResp base,
}
struct GetFriendListRequest{

}
struct GetFriendListResponse{
    1: required model.BaseResp base,
    2: optional list<model.UserInfo> data
}
struct DeleteFriendRequest{
    1:required string id
}
struct DeleteFriendResponse{
         1: required model.BaseResp base,
}
struct VerifyFriendRequest{
    1: required string id,
    2: required string friend_id
}
struct VerifyFriendResponse{
     1: required model.BaseResp base,
     2: required bool friend_exist
}
service UserService {
    GetLoginDataResponse GetLoginData(1: GetLoginDataRequest req),
    GetUserInfoResponse GetUserInfo(1: GetUserInfoRequest request),
    GetLoginDataForYJSYResponse GetGetLoginDataForYJSY(1:GetLoginDataForYJSYRequest request),
    GetInvitationCodeResponse GetInvitationCode(1:GetInvitationCodeRequest request),
    BindInvitationResponse BindInvitation(1:BindInvitationRequest request),
    GetFriendListResponse GetFriendList(1:GetFriendListRequest request),
    DeleteFriendResponse DeleteFriend(1:DeleteFriendRequest request),
    VerifyFriendResponse VerifyFriend(1:VerifyFriendRequest request)
}
