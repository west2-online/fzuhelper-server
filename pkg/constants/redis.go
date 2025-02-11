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

import "time"

const (
	RedisSlowQuery = 10 // ms redis默认的慢查询时间，适用于 logger
)

// Redis Key and Expire Time
const (
	ClassroomKeyExpire    = 2 * 24 * time.Hour
	LaunchScreenKeyExpire = 2 * 24 * time.Hour
	// UserInfoKeyExpire 用户信息过期时间
	UserInfoKeyExpire = 7 * 24 * time.Hour
	// CommonTermListKeyExpire common模块中的学期列表过期时间
	CommonTermListKeyExpire = 7 * 24 * time.Hour
	// CourseTermsKeyExpire course模块中学期列表过期时间
	CourseTermsKeyExpire  = 7 * 24 * time.Hour
	LastLaunchScreenIdKey = "last_launch_screen_id"
	// TermListKey CommonTermListKey
	TermListKey = "term_list"
	// CourseCacheMaxNum CourseList储存最新TopN个学期
	CourseCacheMaxNum = 2
)

// Redis DB Name
const (
	RedisDBEmptyRoom    = 0
	RedisDBLaunchScreen = 1
	RedisDBPaper        = 2
	RedisDBUser         = 3
	RedisDBCommon       = 4
	RedisDBCourse       = 5
)
