namespace go model

struct BaseResp {
    1: i64 code,
    2: string msg,
}

//由前端给的登陆信息，包括id和cookies, 这个struct仅用于测试返回数据，因为登录实现在前端完成，不会在实际项目中使用
struct LoginData {
    1: required string id               //教务处给出的标识，它的组成是时间+学号
    2: required list<string> cookies    //登录凭证，访问资源的时候应该必须携带cookies
}

//Classroom 前端想要返回的fields
struct Classroom {
    1: required string build            //空教室所在楼，例 西三
    2: required string location         //空教室，例 旗山西3-104
    3: required string capacity         //可容纳人数，例 153人
    4: required string type             //教师类型，例 智慧教室普通型
}

