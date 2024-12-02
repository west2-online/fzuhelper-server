namespace go model

struct BaseResp {
    1: i64 code,
    2: string msg,
}

// 由前端给的登陆信息，包括id和cookies, 这个struct仅用于测试返回数据，因为登录实现在前端完成，不会在实际项目中使用
struct LoginData {
    1: required string id               // 教务处给出的标识，它的组成是时间+学号
    2: required list<string> cookies    // 登录凭证，访问资源的时候应该必须携带cookies
}

// 空教室
struct Classroom {
    1: required string build            // 空教室所在楼，例 西三
    2: required string location         // 空教室，例 旗山西3-104
    3: required string capacity         // 可容纳人数，例 153人
    4: required string type             // 教师类型，例 智慧教室普通型
}

// 考场信息
struct ExamRoomInfo {
    1: required string name            // 课程名
    2: required string credit          // 学分
    3: required string teacher         // 任课教师
    4: required string location        // 考场
    5: required string time            // 时间
    6: required string date            // 日期
}

// 课程安排
struct CourseScheduleRule {
    1: required string location
    2: required i64 startClass
    3: required i64 endClass
    4: required i64 startWeek
    5: required i64 endWeek
    6: required i64 weekday
    7: required bool single
    8: required bool double
    9: required bool adjust
}

// 课程信息
struct Course {
    1: required string name
    2: required string teacher
    3: required list<CourseScheduleRule> scheduleRules
    4: required string remark
    5: required string lessonplan
    6: required string syllabus
    7: required string rawScheduleRules
    8: required string rawAdjust
}

// 开屏页
struct Picture{
    1:i64 id,                           // sf自动生成的id
    3:string url,                       // 图片地址
    4:string href,                      // type字段的网址/uri
    5:string text,                      // 开屏页点击区域/工具箱图片下方文字区域的文字
    6:i64 type,                         // 1为空，2为页面跳转，3为app跳转
    7:optional i64 show_times,          // 开屏页被推送展示的次数
    8:optional i64 point_times,         // 点击查看开屏页的次数
    9:i64 duration,                     // 开屏时长（秒）
    10:optional i64 s_type,             // s_type,1为开屏页，2为轮播图，3为生日当天的开屏页
    11:i64 frequency,                   // 一天内的展示次数
    12:i64 start_at,                    // 开始推送的时间戳
    13:i64 end_at,                      // 结束推送的时间戳
    14:i64 start_time,                  // 比如6表示6点
    15:i64 end_time,                    // 比如24 这样就表示6-24点期间会推送该图片
    16:string regex,                    // 推送对象，通过正则里是否有学号来判断是否为推送目标
}


/*
* @Description 又拍云文件目录结构
* @Param basePath 当前所在路径
* @Param files 当前所在目录文件
* @Param folders 当前所在目录下的文件夹
*/
struct UpYunFileDir {
    1: required string basePath,
    2: required list<string> files,
    3: required list<string> folders,
}

// 课程成绩
struct Score {
    1: required string credit      
    2: required string gpa         
    3: required string name         
    4: required string score        
    5: required string teacher     
    6: required string term         
    7: required string year            
}
// 绩点排名
struct GPABean {
    1: required string time
    2: required list<GPAData> data
}

// 绩点信息
struct GPAData {
    1: required string type
    2: required string value
}

// 学分统计字段
struct Credit {
    1: required string type
    2: required string gain
    3: required string total
}
// 统考成绩字段
struct UnifiedExam {
    1: required string name
    2: required string score
    3: required string term
}

// ====== Common ======
// 校历
struct Term {
    1: required string term_id
    2: required string school_year
    3: required string term
    4: required string start_date
    5: required string end_date
}

struct TermEvent {
    1: required string name
    2: required string start_date
    3: required string end_date
}
// ====== END Common ======

struct PaperData {
    1: required string base_path,
    2: required list<string> files,
    3: required list<string> folders,
}

struct PaperUrlData {
    1: required string url,
}
