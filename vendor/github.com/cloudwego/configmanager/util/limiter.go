// Copyright 2023 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"sync/atomic"
	"time"

	"github.com/cloudwego/configmanager/iface"
)

var _ iface.Limiter = (*FrequencyLimiter)(nil)

// FrequencyLimiter is used to limit the frequency of any action.
type FrequencyLimiter struct {
	allowed  int32
	remain   int32
	interval time.Duration
	exit     int32
}

// NewFrequencyLimiter creates a new FrequencyLimiter with the given interval and allowed count.
func NewFrequencyLimiter(interval time.Duration, allowed int32) *FrequencyLimiter {
	if interval <= 0 {
		panic("interval must be positive")
	}
	l := &FrequencyLimiter{
		allowed:  allowed,
		remain:   allowed,
		interval: interval,
	}
	go l.refresh()
	return l
}

// Stop stops the refresh of this limiter.
func (l *FrequencyLimiter) Stop() {
	atomic.StoreInt32(&l.exit, 1)
}

// refresh refreshes remain count atomically for each interval passed.
func (l *FrequencyLimiter) refresh() {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()
	for range ticker.C {
		atomic.StoreInt32(&l.remain, l.allowed)
		if atomic.LoadInt32(&l.exit) == 1 {
			return
		}
	}
}

// Allow returns true if the action is allowed.
func (l *FrequencyLimiter) Allow() bool {
	return atomic.AddInt32(&l.remain, -1) >= 0
}
