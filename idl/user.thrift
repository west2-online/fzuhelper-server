namespace go user




// 获取用户信息
struct UserInfoRequest {

}

struct UserInfoResp {
    1: required string code,
    2: required string message,
    3: required UserInfoData data,
    
}

struct UserInfoData {
    1: required string id,
    2: required string sex,
    3: required string birthday,
    4: required string phone,
    5: required i64 grade,
    6: required string college,
    7: required string major,
    8: required string statusChange,
}


// 获取验证码结果
struct ValidateCodeRequest {
    1: optional string image,
}

struct ValidateCodeResp {
    1: required string code,
    2: required string message,
    3: required string data,
}


// 修改密码
struct ChangePasswordRequest {

    1: required string original,
    2: required string new,
}

struct ChangePasswordResp {
    1: required string code,
    2: required string message,
}

// 校历
struct SchoolCalendarRequest {
    1: optional string term,
}

struct SchoolCalendarResponse {
    1: required string code,
    2: required string message,
    3: required list<SchoolCalendarData> data,
}

struct SchoolCalendarData {
    1: required string dateBegin,
    2: required string dateEnd,
    3: required string name,
}



service UserService {
    UserInfoResp GetUserInfo(1: UserInfoRequest req),
    ValidateCodeResp ValidateCode(1: ValidateCodeRequest req),
    ChangePasswordResp ChangePassword(1: ChangePasswordRequest req),
    SchoolCalendarResponse GetSchoolCalendar(1: SchoolCalendarRequest req),
}