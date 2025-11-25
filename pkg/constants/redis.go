/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package constants

const (
	RedisSlowQuery = 10 // ms redis默认的慢查询时间，适用于 logger
)

const (
	// CourseCacheMaxNum CourseList储存最新TopN个学期
	CourseCacheMaxNum = 2
	// KeyNeverExpire key 永不过期
	KeyNeverExpire = 0
)

// Expire Time
const (
	ClassroomKeyExpire          = 2 * ONE_DAY     // [classroom] 空教室
	LaunchScreenKeyExpire       = 2 * ONE_DAY     // [launch_screen] 开屏页
	UserInfoKeyExpire           = 1 * ONE_WEEK    // [user] 用户信息
	CommonTermListKeyExpire     = 1 * ONE_WEEK    // [common] 学期列表
	CourseTermsKeyExpire        = 3 * ONE_DAY     // [course] 学期列表
	TermInfoKeyExpire           = 7 * ONE_DAY     // [common] 学期详细信息
	ExamRoomKeyExpire           = 10 * ONE_MINUTE // [classroom] 考场信息
	PaperFileDirKeyExpire       = 2 * ONE_DAY     // [paper] 历年卷文件目录
	AcademicScoresExpire        = 5 * ONE_MINUTE  // [academic] 成绩信息
	VisitExpire                 = 1 * ONE_DAY     // [version]访问统计
	LocateDateExpire            = 1 * ONE_HOUR    // [course] 定位日期
	UserInvitationCodeKeyExpire = 1 * ONE_DAY
	UserFriendKeyExpire         = 3 * ONE_DAY
)

// Key Name
const (
	TermListKey                   = "term_list"                    // [common]
	ContributorJwchKey            = "contributor:jwch"             // [common]
	ContributorYJSYKey            = "contributor:yjsy"             // [common]
	ContributorFzuhelperAppKey    = "contributor:fzuhelper-app"    // [common]
	ContributorFzuhelperServerKey = "contributor:fzuhelper-server" // [common]
	LastLaunchScreenIdKey         = "last_launch_screen_id"        // [launch_screen]
	LocateDateKey                 = "locateDate"                   // [course]
)

// DB Name
const (
	RedisDBEmptyRoom    = 0
	RedisDBLaunchScreen = 1
	RedisDBPaper        = 2
	RedisDBUser         = 3
	RedisDBCommon       = 4
	RedisDBCourse       = 5
	RedisDBAcademic     = 6
	RedisDBVersion      = 7
	RedisDBOA           = 8
)
