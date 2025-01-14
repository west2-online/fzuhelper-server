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

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const (
	baseDelay = 2
	maxDelay  = 1000
	maxTokens = 100
	maxRate   = rate.Limit(10)

	FailureRateLimiterBaseDelay = time.Minute
	FailureRateLimiterMaxDelay  = 30 * time.Minute
)

// NoticeSyncer 教务处通知同步
type NoticeSyncer struct {
	workQueue workqueue.TypedRateLimitingInterface[string]
	db        *db.Database
}

// InitNoticeSyncer will init a worker queue
func InitNoticeSyncer(db *db.Database) *NoticeSyncer {
	ns := &NoticeSyncer{
		db: db,
		workQueue: workqueue.NewTypedRateLimitingQueue[string](
			workqueue.NewTypedMaxOfRateLimiter(
				// For syncRec failures(i.e. doRecommend return err), the retry time is (2*minutes)*2^<num-failures>
				// The maximum retry time is 24 hours
				workqueue.NewTypedItemExponentialFailureRateLimiter[string](FailureRateLimiterBaseDelay, FailureRateLimiterMaxDelay),
				// 10 qps, 100 bucket size. This is only for retry speed, it's only the overall factor (not per item)
				// 每秒最多产生 10 个令牌（允许处理 10 个任务）。
				// 100：令牌桶最多存储 100 个令牌，允许积累的最大任务数量
				workqueue.NewTypedMaxOfRateLimiter[string](workqueue.NewTypedItemExponentialFailureRateLimiter[string](baseDelay, maxDelay),
					&workqueue.TypedBucketRateLimiter[string]{Limiter: rate.NewLimiter(maxRate, maxTokens)}),
			),
		),
	}
	return ns
}

func (ns *NoticeSyncer) Start() {
	ns.initNoticeSyncer()
	for i := 0; i < constants.NoticeWorker; i++ {
		go ns.worker()
	}
}

func (ns *NoticeSyncer) Add(key string) {
	ns.workQueue.Add(key)
}

// worker 将每 8 小时对 jwc 的通知进行同步存储    template from classroom.syncer
func (ns *NoticeSyncer) worker() {
	for {
		item, shutdown := ns.workQueue.Get()
		if shutdown {
			logger.Info("Notice.worker shutdown")
			return
		}
		// do update
		if err := ns.update(); err != nil {
			// 失败后采取重试策略
			logger.Errorf("Notice.worker: failed to update: %v", err)
			ns.workQueue.AddRateLimited(item)
		} else {
			ns.workQueue.AddAfter(item, constants.NoticeUpdateTime)
			ns.workQueue.Forget(item)
		}
		ns.workQueue.Done(item)
		logger.Info("Notice.worker update complete")
	}
}
