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

package syncer

import (
	"time"

	"golang.org/x/time/rate"
	"k8s.io/client-go/util/workqueue"

	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

type EmptyRoomSyncer struct {
	workQueue workqueue.TypedRateLimitingInterface[string]
	cache     *cache.Cache
}

// InitEmptyRoomSyncer will init a worker queue
func InitEmptyRoomSyncer(cache *cache.Cache) *EmptyRoomSyncer {
	ers := &EmptyRoomSyncer{
		cache: cache,
		workQueue: workqueue.NewTypedRateLimitingQueue[string](
			workqueue.NewTypedMaxOfRateLimiter(
				// For syncRec failures(i.e. doRecommend return err), the retry time is (2*minutes)*2^<num-failures>
				// The maximum retry time is 24 hours
				workqueue.NewTypedItemExponentialFailureRateLimiter[string](constants.FailureRateLimiterBaseDelay, constants.FailureRateLimiterMaxDelay),
				// 10 qps, 100 bucket size. This is only for retry speed, it's only the overall factor (not per item)
				// 每秒最多产生 10 个令牌（允许处理 10 个任务）。
				// 100：令牌桶最多存储 100 个令牌，允许积累的最大任务数量
				workqueue.NewTypedMaxOfRateLimiter[string](workqueue.NewTypedItemExponentialFailureRateLimiter[string](5, 1000), &workqueue.TypedBucketRateLimiter[string]{Limiter: rate.NewLimiter(rate.Limit(10), 100)}),
			),
		),
	}
	return ers
}

// Start will launch the worker
func (ers *EmptyRoomSyncer) Start() {
	for i := 0; i < constants.ClassroomWorker; i++ {
		go ers.worker()
	}
}

func (ers *EmptyRoomSyncer) Add(key string) {
	ers.workQueue.Add(key)
}

// 不断接受来自队列的任务
func (ers *EmptyRoomSyncer) worker() {
	for {
		item, shutdown := ers.workQueue.Get()
		if shutdown {
			logger.Info("Classroom.worker shutdown")
			return
		}
		var err error
		switch item {
		case "update":
			if err = ers.apply(constants.ClassroomUpdatedTime, item, ers.UpdateClassroomsInfo); err != nil {
				logger.Infof("Classroom.worker: update failed: %v", err)
			}
		case "schedule":
			if err = ers.apply(constants.ClassroomScheduledTime, item, ers.ScheduledUpdateClassroomsInfo); err != nil {
				logger.Infof("Classroom.worker: schedule failed: %v", err)
			}
		}
	}
}

// apply will trigger func and select queue strategy based the error
func (ers *EmptyRoomSyncer) apply(timeDuration time.Duration, item string, scheduledFunc func(date time.Time) error) (err error) {
	if err = scheduledFunc(time.Now()); err != nil {
		// 如果失败，在使用该函数，采取避退策略
		ers.workQueue.AddRateLimited(item)
	} else {
		// 将signal重新放入队列，实现自驱动
		ers.workQueue.AddAfter(item, timeDuration)
		// 任务成功，清除其失败记录
		ers.workQueue.Forget(item)
	}
	// 任务完成, 释放资源, 防止队列阻塞
	ers.workQueue.Done(item)
	return err
}
