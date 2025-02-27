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
	NoticeWorker     = 1
	NoticeUpdateTime = 1 * time.Hour // (notice) 通知更新间隔
	NoticePageSize   = 20            // 教务处教学通知一页大小固定 20
)
