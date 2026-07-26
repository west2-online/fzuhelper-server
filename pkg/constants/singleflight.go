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

// 此处为固定key
const (
	// 学期列表是全局校历数据，不依赖具体用户，可以使用固定 key 合并请求。
	SingleflightTermListKey = "term_list"

	// 下载地址按发布渠道拆 key，避免 release/beta 并发时复用到另一个渠道的地址。
	SingleflightDownloadReleaseKey = "download_release"
	SingleflightDownloadBetaKey    = "download_beta"

	// 版本信息按发布渠道拆 key，避免 release/beta 版本元数据互相复用。
	SingleflightReleaseVersionKey = "release_version"
	SingleflightBetaVersionKey    = "beta_version"

	SingleflightCloudKey = "cloud"

	// Android 接口一次返回 release 和 beta 两份数据，使用独立 key 避免和单独版本查询混用。
	SingleflightAndroidVersionKey = "android_version"
)

// 此处为动态key的 prefix，需调用sf.key函数生成完整 key
const (
	// 本科和研究生成绩来自不同上游，返回结构也不同，需要按身份隔离。
	SingleflightScoresPrefix = "academic:scores"

	// 考场结果按学期和身份区分，同一学生不同学期或不同身份不能复用结果。
	SingleflightExamRoomsPrefix = "classroom:exam_rooms"

	// 将刷新标记纳入 key，避免强刷请求复用普通请求的 singleflight 结果。
	SingleflightCourseListPrefix = "course:list"

	// 本科和研究生学期来源不同，同一学号也要按身份隔离。
	SingleflightCourseTermsPrefix = "course:terms"

	SingleflightTermPrefix       = "common:term"
	SingleflightNoticePrefix     = "common:notice"
	SingleflightPaperDirPrefix   = "paper:dir"
	SingleflightFriendListPrefix = "user:friend_list"

	// 本科和研究生用户信息来自不同上游，按身份隔离避免复用到错误来源的数据。
	SingleflightUserInfoPrefix = "user:info"
)
