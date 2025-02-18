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
	CheckFileTypeBufferSize = 512 // 适用于判断文件类型，需要读取前512个字节

	SnowflakeWorkerID     = 0
	SnowflakeDatacenterID = 0

	ClassroomWorker        = 1              // (class_room) 同时启用的 goroutine 数量
	ClassroomScheduledTime = 24 * time.Hour // (class_room) 空教室非当天同步时间
	ClassroomUpdatedTime   = 6 * time.Hour  // (class_room) 当天空教室更新间隔

	NoticeWorker     = 1
	NoticeUpdateTime = 8 * time.Hour // (notice) 通知更新间隔
	NoticePageSize   = 20            // 教务处教学通知一页大小固定 20

	AcademicWorker = 1 // (academic)同时启动的 goroutine 数量

	CacheFileDir = "UssFileDir" // (paper) 文件缓存目录

	// ValidateCodeURL 获取验证码结果的本地python服务url，需要保证 login-verify 和 api 处于同一个 dokcer 网络中
	ValidateCodeURL = "http://login-verify:8081/api/v1/jwch/user/validateCode"
	// UmengURL 友盟推送 API
	UmengURL = "https://msgapi.umeng.com/api/send"
)

const (
	OneDay = 24 * time.Hour // 一天的时间，适用于 classroom、paper 等服务计算时间
)

// CampusArray 校区数组
var CampusArray = []string{"旗山校区", "厦门工艺美院", "铜盘校区", "怡山校区", "晋江校区", "泉港校区"}
