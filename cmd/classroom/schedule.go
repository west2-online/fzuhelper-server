package main

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/service"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
	"golang.org/x/time/rate"
	"k8s.io/client-go/util/workqueue"
	"strconv"
	"time"
)

var WorkQueue workqueue.RateLimitingInterface

func InitWorkerQueue() {
	WorkQueue = workqueue.NewNamedRateLimitingQueue(
		workqueue.NewMaxOfRateLimiter(
			// For syncRec failures(i.e. doRecommend return err), the retry time is (2*minutes)*2^<num-failures>
			// The maximum retry time is 24 hours
			workqueue.NewItemExponentialFailureRateLimiter(constants.FailureRateLimiterBaseDelay, constants.FailureRateLimiterMaxDelay),
			// 10 qps, 100 bucket size. This is only for retry speed, it's only the overall factor (not per item)
			//每秒最多产生 10 个令牌（允许处理 10 个任务）。
			//100：令牌桶最多存储 100 个令牌，允许积累的最大任务数量
			&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
		),
		constants.ClassroomService)
	for i := 0; i < constants.NumWorkers; i++ {
		go worker()
	}
}

// CacheEmptyRooms 缓存所有空教室的数据
// 缓存一周的所有信息，每两天更新一次
func CacheEmptyRooms() {
	for {
		var dates []string
		currentTime := time.Now()
		//设定一周时间
		for i := 0; i < 7; i++ {
			date := currentTime.AddDate(0, 0, i).Format("2006-01-02")
			dates = append(dates, date)
		}
		for _, date := range dates {
			for _, campus := range constants.CampusArray {
				for startTime := 1; startTime <= 11; startTime++ {
					for endTime := startTime; endTime <= 11; endTime++ {
						args := &classroom.EmptyRoomRequest{
							Date:      date,
							Campus:    campus,
							StartTime: strconv.Itoa(startTime),
							EndTime:   strconv.Itoa(endTime),
						}
						WorkQueue.Add(args)
						logger.LoggerObj.Debugf("classroom.service.CacheEmptyRooms add task %v", args)
					}
				}
				logger.LoggerObj.Infof("classroom.service.CacheEmptyRooms add all tasks of campus %v in the day %v", campus, date)
			}
		}
		time.Sleep(constants.ScheduledTime)
	}
}

// 从工作队列取出task并处理
func worker() {
	for {
		ctx := context.Background()
		id, cookies := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password).GetIdentifierAndCookies()
		l := service.NewClassroomService(ctx, id, cookies)
		task, shutDown := WorkQueue.Get()
		if shutDown {
			logger.LoggerObj.Debug("classroom.service.worker worker shutDown")
			return
		}
		func(task any) {
			defer WorkQueue.Done(task)
			args, ok := task.(*classroom.EmptyRoomRequest)
			if !ok {
				logger.LoggerObj.Errorf("classroom.service.worker task type error: %T", task)
				return
			}
			_, err := l.GetEmptyRooms(args)
			if err != nil {
				logger.LoggerObj.Errorf("classroom.service.worker GetEmptyRooms failed, args %v: %v", err, args)
				return
			}
			//将任务标记为完成
			WorkQueue.Forget(task)
			logger.LoggerObj.Debug("classroom.service.worker task %v done", args)
		}(task)
	}
}
