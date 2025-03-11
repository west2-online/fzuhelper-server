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

// Classroom 空教室
const (
	ClassroomWorker        = 1            //  同时启用的 goroutine 数量
	ClassroomScheduledTime = ONE_DAY      // 空教室非当天同步时间
	ClassroomUpdatedTime   = 6 * ONE_HOUR // 当天空教室更新间隔
)

// notice 教务处教学通知
const (
	NoticeTaskKey    = "notice"
	NoticeUpdateTime = 1 * time.Hour // (notice) 通知更新间隔
	NoticePageSize   = 20            // 教务处教学通知一页大小固定 20
)

// course 课程信息
const (
	LocateDateTaskKey = "locateDate"
)

// version
const (
	VersionVisitedTaskKey = "versionVisited"
)

// contributor 贡献者信息
const (
	ContributorTaskKey        = "contributor"
	ContributorInfoUpdateTime = 24 * 7 * time.Hour // 贡献者信息更新间隔

	// 下方的服务器由 @renbaoshuo 维护，如果挂了请通过工作室渠道联系。注释是反代的源 URL
	// "https://api.github.com/repos/west2-online/jwch/contributors"
	ContributorJwch = "https://fuu.api.baoshuo.dev/contributors/jwch"
	// "https://api.github.com/repos/west2-online/yjsy/contributors"
	ContributorYJSY = "https://fuu.api.baoshuo.dev/contributors/yjsy"
	// "https://api.github.com/repos/west2-online/fzuhelper-app/contributors"
	ContributorFzuhelperApp = "https://fuu.api.baoshuo.dev/contributors/fzuhelper-app"
	// "https://api.github.com/repos/west2-online/fzuhelper-server/contributors"
	ContributorFzuhelperServer = "https://fuu.api.baoshuo.dev/contributors/fzuhelper-server"
	// "https://avatars.githubusercontent.com/u/%s
	AvatarProxy = "https://fuu.api.baoshuo.dev/avatar/%s"
)
