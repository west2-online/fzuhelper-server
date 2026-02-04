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

package umeng

import (
	"sync"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// asyncDispatcher 负责异步消费 Umeng 发送任务并执行限流。
// 设计目标：
// 1) 异步：发送端仅入队，不阻塞业务线程或主任务队列。
// 2) 限流：支持最小间隔控制 + 每日配额控制。
// 3) 单例：进程内仅启动一个 dispatcher，确保全局限流一致性。
type asyncDispatcher struct {
	// ch 为任务队列通道，元素是“发送函数”。
	// 发送函数返回 error，便于统一记录失败日志。
	ch chan func() error
	// interval 为相邻两次发送的最小间隔。
	interval time.Duration
	// dailyLimit 为每日最大允许发送次数。
	dailyLimit int
	// dailyCount 为当天已发送次数。
	dailyCount int
	// lastResetDate 记录上次重置日期，用于跨日清零计数。
	lastResetDate time.Time
	// lastRequestTime 记录上一次实际发送的时间，用于间隔限流。
	lastRequestTime time.Time
}

var (
	// dispatcherOnce 确保 dispatcher 只被初始化一次。
	dispatcherOnce sync.Once
	// dispatcher 为全局单例实例。
	dispatcher *asyncDispatcher
)

// getDispatcher 获取全局 dispatcher 单例。
// 首次调用会完成初始化并启动后台消费协程。
// 该方法为内部使用，外部通过 EnqueueAsync 入队即可。
func getDispatcher() *asyncDispatcher {
	dispatcherOnce.Do(func() {
		dispatcher = newAsyncDispatcher(constants.UmengAsyncQueueSize, constants.UmengRateLimitDelay, constants.UmengDailyLimit)
		// 后台消费协程：串行处理任务，确保限流语义正确。
		go dispatcher.run()
	})
	return dispatcher
}

// newAsyncDispatcher 创建一个新的 dispatcher 实例。
// 参数：
// - queueSize：队列缓冲长度，<=0 时使用默认值。
// - interval：发送最小间隔。
// - dailyLimit：每日最大发送次数，<=0 时使用默认值。
// 返回值仅在 getDispatcher 中使用，避免重复创建。
func newAsyncDispatcher(queueSize int, interval time.Duration, dailyLimit int) *asyncDispatcher {
	if queueSize <= 0 {
		queueSize = constants.UmengAsyncQueueSize
	}
	if dailyLimit <= 0 {
		dailyLimit = constants.UmengDailyLimit
	}
	return &asyncDispatcher{
		ch:              make(chan func() error, queueSize),
		interval:        interval,
		dailyLimit:      dailyLimit,
		dailyCount:      0,
		lastResetDate:   time.Now(),
		lastRequestTime: time.Now().Add(-interval),
	}
}

// EnqueueAsync 将 Umeng 发送任务放入异步队列。
// 特性：
// - 非阻塞：队列满时立即返回 false，不阻塞业务线程。
// - 安全：task 为空时直接返回 false。
// - 单例：内部确保 dispatcher 只初始化一次。
// 用法：将实际发送逻辑包装成闭包传入，例如：
//
//	umeng.EnqueueAsync(func() error { return umeng.SendAndroidGroupcastWithGoApp(...) })
func EnqueueAsync(task func() error) bool {
	if task == nil {
		return false
	}
	d := getDispatcher()
	select {
	case d.ch <- task:
		return true
	default:
		return false
	}
}

// run 后台消费循环。
// 该循环串行读取队列并执行任务，先限流再发送，保证顺序与配额一致。
// 注意：此方法应仅在后台协程中运行。
func (d *asyncDispatcher) run() {
	for task := range d.ch {
		d.wait()
		if err := task(); err != nil {
			logger.Errorf("umeng async task failed: %v", err)
		}
	}
}

// wait 执行限流等待。
// 逻辑顺序：
// 1) 跨日判断：新的一天重置 dailyCount。
// 2) 每日配额：超过 dailyLimit 则等待到次日零点并重置。
// 3) 间隔限制：确保相邻发送间隔不小于 interval。
// 该方法在后台消费协程中调用，因此可以阻塞而不影响业务线程。
func (d *asyncDispatcher) wait() {
	now := time.Now()
	if !sameDay(now, d.lastResetDate) {
		d.dailyCount = 0
		d.lastResetDate = now
	}

	if d.dailyCount >= d.dailyLimit {
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		time.Sleep(time.Until(nextDay))
		d.dailyCount = 0
		d.lastResetDate = time.Now()
		d.lastRequestTime = time.Now().Add(-d.interval)
	}

	now = time.Now()
	elapsed := now.Sub(d.lastRequestTime)
	if elapsed < d.interval {
		time.Sleep(d.interval - elapsed)
		now = time.Now()
	}

	d.lastRequestTime = now
	d.dailyCount++
}

// sameDay 判断两个时间是否在同一天（按本地时区）。
// 用于每日配额的跨日判断。
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
