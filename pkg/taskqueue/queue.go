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

package taskqueue

import (
	"sync"
	"time"

	"k8s.io/client-go/util/workqueue"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

type TaskQueue interface {
	Start()
	Add(key string, task QueueTask)
	AddSchedule(key string, task ScheduleQueueTask)
	worker()
}

// QueueTask 队列任务，使用指数退避和令牌桶限流
type QueueTask struct {
	Execute func() error
}

// ScheduleQueueTask 定时任务
type ScheduleQueueTask struct {
	Execute         func() error
	GetScheduleTime func() time.Duration
}

type BaseTaskQueue struct {
	workQueue workqueue.TypedRateLimitingInterface[string]
	taskMap   sync.Map
}

func NewBaseTaskQueue() *BaseTaskQueue {
	return &BaseTaskQueue{
		// 默认限流器
		// - 单任务重试采用指数退避策略：初始延迟为 5ms，最大延迟为 1000 秒。
		// - 整体速率限制：每秒最多 10 次请求，桶大小为 100 个令牌。
		workQueue: workqueue.NewTypedRateLimitingQueue[string](
			workqueue.DefaultTypedControllerRateLimiter[string](),
		),
	}
}

// Add 想task queue 中添加 task
// ScheduleQueueTask 也实现了 QueueTask 的接口，不需要显示声明
func (btq *BaseTaskQueue) Add(key string, task QueueTask) {
	btq.taskMap.Store(key, task)
	btq.workQueue.Add(key)
}

func (btq *BaseTaskQueue) AddSchedule(key string, task ScheduleQueueTask) {
	btq.taskMap.Store(key, task)
	btq.workQueue.Add(key)
}

func (btq *BaseTaskQueue) Start() {
	for i := 0; i < constants.WorkerNumber; i++ {
		go btq.worker()
	}
}

func (btq *BaseTaskQueue) worker() {
	for {
		key, shutdown := btq.workQueue.Get()
		if shutdown {
			logger.Info("BaseTaskQueue:worker shutdown")
			return
		}

		task, exists := btq.taskMap.Load(key)
		if !exists {
			btq.workQueue.Done(key)
			logger.Warnf("Task evaporated: %s", key)
			continue
		}
		switch task := task.(type) {
		case ScheduleQueueTask:
			if err := task.Execute(); err != nil {
				btq.workQueue.AddRateLimited(key)
				logger.Errorf("ScheduleQueueTask execute failed: %v", err)
			} else {
				btq.workQueue.AddAfter(key, task.GetScheduleTime())
				btq.workQueue.Forget(key)
			}
			btq.workQueue.Done(key)
		case QueueTask:
			if err := task.Execute(); err != nil {
				btq.workQueue.AddRateLimited(key)
				logger.Errorf("QueueTask execute failed: %v", err)
			} else {
				btq.taskMap.Delete(key)
			}
			btq.workQueue.Done(key)
		default:
			logger.Errorf("BaseTaskQueue:Unknown task type，key: %v", key)
		}
	}
}
