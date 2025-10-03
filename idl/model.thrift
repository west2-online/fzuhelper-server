namespace go model

struct BaseResp {
    1: i64 code,
    2: string msg,
}

// 由前端给的登陆信息，包括id和cookies, 这个struct仅用于测试返回数据，因为登录实现在前端完成，不会在实际项目中使用
struct LoginData {
    1: required string id         // 教务处给出的标识，它的组成是时间+学号
    2: required string cookies    // 登录凭证，访问资源的时候应该必须携带cookies
}

struct UserInfo{
    1: required string stu_id,
    2: required string name,
    3: required string birthday,
    4: required string sex,
    5: required string college,
    6: required string grade,
    7: required string major,
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
    1: required string location         // 定制
    2: required i64 startClass          // 开始节数
    3: required i64 endClass            // 结束节数
    4: required i64 startWeek           // 起始周
    5: required i64 endWeek             // 结束周
    6: required i64 weekday             // 星期几
    7: required bool single             // 单周
    8: required bool double             // 双周
    9: required bool adjust             // 是否是调课
}

// 课程信息
struct Course {
    1: required string name                             // 课程名称
    2: required string teacher                          // 教师
    3: required list<CourseScheduleRule> scheduleRules  // 排课规则
    4: required string remark                           // 备注
    5: required string lessonplan                       // 授课计划
    6: required string syllabus                         // 教学大纲
    7: required string rawScheduleRules                 // (原始数据) 排课规则
    8: required string rawAdjust                        // (原始数据) 调课规则
    9: required string examType                        // 考试类型(用于查看是否免听)
}

// 当前周数、学期、学年
struct LocateDate {
    1: required string week
    2: required string year
    3: required string term
    4: required string date
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


// 又拍云文件目录结构
struct UpYunFileDir {
    1: optional string basePath,        // 当前所在路径
    2: required list<string> files,     // 当前所在目录文件
    3: required list<string> folders,   // 当前所在目录下的文件夹
}

// 课程成绩
struct Score {
    1: required string credit           // 学分
    2: required string gpa              // 绩点
    3: required string name             // 课程名
    4: required string score            // 得分
    5: required string teacher          // 授课教师
    6: required string term             // 学期
    7: required string exam_type        // 考试类型
    8: required string elective_type    // 选修类型
    9: required string classroom        // 上课地点
}

// 绩点排名
struct GPABean {
    1: required string time             // 更新时间
    2: required list<GPAData> data      // 数据
}

// 绩点信息
struct GPAData {
    1: required string type             // 类型（如修读类别或总学分）
    2: required string value            // 信息（对应的信息）
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

// 学分详细数据项
struct CreditDetail {
    1: required string key
    2: required string value
}

// 学分分类
struct CreditCategory {
    1: required string type
    2: required list<CreditDetail> data
}

// 学分响应
typedef list<CreditCategory> CreditResponse

// 又拍云文件目录结构,兼容旧版安卓
struct PaperData {
    1: optional string base_path,       // 当前所在路径
    2: required list<string> files,     // 当前所在目录文件，使用required保证files不为nil
    3: required list<string> folders,   // 当前所在目录下的文件夹，使用required保证folders不为nil
}

struct PaperUrlData {
    1: required string url,
}

// ====== Common ======
// 校历
struct Term {
    1: optional string term_id
    2: optional string school_year
    3: optional string term
    4: optional string start_date
    5: optional string end_date
}

struct TermEvent {
    1: optional string name
    2: optional string start_date
    3: optional string end_date
}

struct TermList {
    1: optional string current_term
    2: optional list<Term> terms
}

struct TermInfo {
    1: optional string term_id
    2: optional string term
    3: optional string school_year
    4: optional list<TermEvent> events
}

struct NoticeInfo {
    1: optional string title
    2: optional string url
    3: optional string date
}

struct Contributor {
  1: string name
  2: string avatar_url
  3: string url
  4: i64 contributions
}

struct ToolboxConfig {
    1: required i64 tool_id
    2: optional bool visible
    3: optional string name
    4: optional string icon
    5: optional string type
    6: optional string message
    7: optional string extra
    8: optional string platform
    9: optional i64 version
}

// ====== END Common ======

// version
struct Version{
    1: optional string version_code
    2: optional string version_name
    3: optional bool force
    4: optional string changelog
    5: optional string url
}
