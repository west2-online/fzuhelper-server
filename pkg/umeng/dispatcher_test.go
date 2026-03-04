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
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func setMockDispatcherForTest(d *asyncDispatcher) {
	dispatcher = d
	dispatcherOnce = sync.Once{}
	dispatcherOnce.Do(func() {})
}

func resetDispatcherForTest() {
	dispatcher = nil
	dispatcherOnce = sync.Once{}
}

func TestNewAsyncDispatcher(t *testing.T) {
	type testCase struct {
		name            string
		expectQueueSize int
		expectInterval  time.Duration
		expectDaily     int
	}

	testCases := []testCase{
		{
			name:            "UseConstants",
			expectQueueSize: constants.UmengAsyncQueueSize,
			expectInterval:  constants.UmengRateLimitDelay,
			expectDaily:     constants.UmengDailyLimit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := newAsyncDispatcher()

			assert.NotNil(t, d)
			assert.Equal(t, tc.expectQueueSize, cap(d.ch))
			assert.Equal(t, tc.expectInterval, d.interval)
			assert.Equal(t, tc.expectDaily, d.dailyLimit)
			assert.Equal(t, 0, d.dailyCount)
		})
	}
}

func TestGetDispatcher(t *testing.T) {
	resetDispatcherForTest()
	t.Cleanup(resetDispatcherForTest)

	d1 := getDispatcher()
	d2 := getDispatcher()

	assert.NotNil(t, d1)
	assert.Same(t, d1, d2)
}

func TestEnqueueAsync(t *testing.T) {
	type testCase struct {
		name       string
		task       func() error
		prepare    func()
		expectEnq  bool
		expectRead bool
	}

	testCases := []testCase{
		{
			name:      "NilTask",
			task:      nil,
			prepare:   resetDispatcherForTest,
			expectEnq: false,
		},
		{
			name: "EnqueueSuccess",
			task: func() error { return nil },
			prepare: func() {
				setMockDispatcherForTest(&asyncDispatcher{ch: make(chan func() error, 1)})
			},
			expectEnq:  true,
			expectRead: true,
		},
		{
			name: "QueueFull",
			task: func() error { return nil },
			prepare: func() {
				d := &asyncDispatcher{ch: make(chan func() error, 1)}
				d.ch <- func() error { return nil }
				setMockDispatcherForTest(d)
			},
			expectEnq: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetDispatcherForTest()
			t.Cleanup(resetDispatcherForTest)
			if tc.prepare != nil {
				tc.prepare()
			}

			ok := EnqueueAsync(tc.task)
			assert.Equal(t, tc.expectEnq, ok)

			if tc.expectRead {
				select {
				case <-dispatcher.ch:
				default:
					t.Fatalf("expected task enqueued but channel is empty")
				}
			}
		})
	}
}

func TestSameDay(t *testing.T) {
	type testCase struct {
		name   string
		a      time.Time
		b      time.Time
		expect bool
	}

	now := time.Now()
	testCases := []testCase{
		{
			name:   "SameDay",
			a:      now,
			b:      now.Add(2 * time.Hour),
			expect: true,
		},
		{
			name:   "DifferentDay",
			a:      now,
			b:      now.Add(24 * time.Hour),
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, sameDay(tc.a, tc.b))
		})
	}
}

func TestAsyncDispatcherWait(t *testing.T) {
	type testCase struct {
		name        string
		dispatcher  *asyncDispatcher
		beforeWait  func(t *testing.T)
		afterAssert func(t *testing.T, d *asyncDispatcher, start time.Time)
	}

	testCases := []testCase{
		{
			name: "ResetOnNextDay",
			dispatcher: &asyncDispatcher{
				ch:              make(chan func() error, 1),
				interval:        0,
				dailyLimit:      10,
				dailyCount:      5,
				lastResetDate:   time.Now().AddDate(0, 0, -1),
				lastRequestTime: time.Now(),
			},
			afterAssert: func(t *testing.T, d *asyncDispatcher, _ time.Time) {
				assert.Equal(t, 1, d.dailyCount)
				assert.True(t, sameDay(time.Now(), d.lastResetDate))
				assert.False(t, d.lastRequestTime.IsZero())
			},
		},
		{
			name: "RespectInterval",
			dispatcher: &asyncDispatcher{
				ch:              make(chan func() error, 1),
				interval:        20 * time.Millisecond,
				dailyLimit:      10,
				dailyCount:      0,
				lastResetDate:   time.Now(),
				lastRequestTime: time.Now(),
			},
			afterAssert: func(t *testing.T, d *asyncDispatcher, start time.Time) {
				elapsed := time.Since(start)
				assert.True(t, elapsed >= 15*time.Millisecond)
				assert.Equal(t, 1, d.dailyCount)
			},
		},
		{
			name: "DailyLimitReached",
			dispatcher: &asyncDispatcher{
				ch:              make(chan func() error, 1),
				interval:        0,
				dailyLimit:      1,
				dailyCount:      1,
				lastResetDate:   time.Now(),
				lastRequestTime: time.Now(),
			},
			beforeWait: func(t *testing.T) {
				mockey.Mock(time.Sleep).To(func(d time.Duration) {
					assert.GreaterOrEqual(t, d, time.Millisecond)
				}).Build()
			},
			afterAssert: func(t *testing.T, d *asyncDispatcher, _ time.Time) {
				assert.Equal(t, 1, d.dailyCount)
				assert.True(t, sameDay(time.Now(), d.lastResetDate))
				assert.False(t, d.lastRequestTime.IsZero())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer mockey.UnPatchAll()
			if tc.beforeWait != nil {
				tc.beforeWait(t)
			}
			start := time.Now()
			tc.dispatcher.wait()
			tc.afterAssert(t, tc.dispatcher, start)
		})
	}
}

func TestAsyncDispatcherRun(t *testing.T) {
	d := &asyncDispatcher{
		ch:              make(chan func() error, 2),
		interval:        0,
		dailyLimit:      10,
		dailyCount:      0,
		lastResetDate:   time.Now(),
		lastRequestTime: time.Now(),
	}

	var mu sync.Mutex
	called := make([]int, 0, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go d.run()

	d.ch <- func() error {
		mu.Lock()
		called = append(called, 1)
		mu.Unlock()
		wg.Done()
		return nil
	}
	d.ch <- func() error {
		mu.Lock()
		called = append(called, 2)
		mu.Unlock()
		wg.Done()
		return errors.New("mock error")
	}
	close(d.ch)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		mu.Lock()
		defer mu.Unlock()
		assert.Equal(t, []int{1, 2}, called)
		assert.Equal(t, 2, d.dailyCount)
	case <-time.After(time.Second):
		t.Fatal("dispatcher run did not finish in time")
	}
}
