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

type AcademicSyncer struct {
	workQueue workqueue.TypedRateLimitingInterface[QueueTask]
}

// InitAcademicSyncer will init a worker queue
func InitAcademicSyncer() *AcademicSyncer {
	ers := &AcademicSyncer{
		workQueue: workqueue.NewTypedRateLimitingQueue[QueueTask](
			// 默认限流器
			// - 单任务重试采用指数退避策略：初始延迟为 5ms，最大延迟为 1000 秒。
			// - 整体速率限制：每秒最多 10 次请求，桶大小为 100 个令牌。
			workqueue.DefaultTypedControllerRateLimiter[QueueTask](),
		),
	}
	return ers
}

// Start will launch the worker
func (ers *AcademicSyncer) Start() {
	for i := 0; i < constants.AcademicWorker; i++ {
		go ers.worker()
	}
}

func (ers *AcademicSyncer) Add(key QueueTask) {
	ers.workQueue.Add(key)
}

// 不断接受来自队列的任务
func (ers *AcademicSyncer) worker() {
	var err error
	for {
		task, shutdown := ers.workQueue.Get()
		if shutdown {
			logger.Info("AcademicSyncer.worker shutdown")
			return
		}
		switch task := task.(type) {
		case *SetScoresCacheTask:
			if err = task.Execute(); err != nil {
				ers.workQueue.AddRateLimited(task)
				logger.Warnf("AcademicSyncer.worker: SetScoresCacheTask failed: %v", err)
			} else {
				ers.workQueue.Done(task)
			}
		default:
			logger.Errorf("Unknown task type: %T", task)
		}
	}
}
