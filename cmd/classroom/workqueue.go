package main

import (
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"golang.org/x/time/rate"
	"k8s.io/client-go/util/workqueue"
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
	go worker()
}

// 不断接受来自队列的任务
func worker() {
	for {
		item, shutdown := WorkQueue.Get()
		if shutdown {
			return
		}
		if err := cache.ScheduledGetClassrooms(); err != nil {
			logger.LoggerObj.Errorf("Classroom.worker ScheduledGetClassrooms failed: %v", err)
			//如果失败，在使用该函数，采取避退策略
			WorkQueue.AddRateLimited(item)
		}
		//将signal重新放入队列，实现自驱动
		WorkQueue.AddAfter(item, constants.ScheduledTime)
		logger.LoggerObj.Info("Classroom.worker ScheduledGetClassrooms success")
	}
}
