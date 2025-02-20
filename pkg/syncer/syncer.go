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
	"k8s.io/client-go/util/workqueue"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

type Syncer interface {
	Start()
	Add(task QueueTask)
	worker()
}

type QueueTask interface {
	Execute() error
}

type BaseSyncer struct {
	workQueue workqueue.TypedRateLimitingInterface[QueueTask]
}

func NewBaseSyncer() *BaseSyncer {
	return &BaseSyncer{
		// 默认限流器
		// - 单任务重试采用指数退避策略：初始延迟为 5ms，最大延迟为 1000 秒。
		// - 整体速率限制：每秒最多 10 次请求，桶大小为 100 个令牌。
		workQueue: workqueue.NewTypedRateLimitingQueue[QueueTask](
			workqueue.DefaultTypedControllerRateLimiter[QueueTask](),
		),
	}
}

func (bs *BaseSyncer) Add(task QueueTask) {
	bs.workQueue.Add(task)
}

func (bs *BaseSyncer) Start() {
	for i := 0; i < constants.AcademicWorker; i++ {
		go bs.worker()
	}
}

func (bs *BaseSyncer) worker() {
	for {
		task, shutdown := bs.workQueue.Get()
		if shutdown {
			logger.Info("BaseSyncer:worker shutdown")
			return
		}
		switch task := task.(type) {
		case QueueTask:
			if err := task.Execute(); err != nil {
				bs.workQueue.AddRateLimited(task)
				logger.Warnf("BaseSyncer:task failed: %v", err)
			} else {
				bs.workQueue.Done(task)
			}
		default:
			logger.Errorf("BaseSyncer:Unknown task type: %T", task)
		}
	}
}
