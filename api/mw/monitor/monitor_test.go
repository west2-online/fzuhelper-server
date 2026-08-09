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

package monitor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildRequestEvent(t *testing.T) {
	type testCase struct {
		name      string
		body      []byte
		traceID   string
		errorCode int64
	}

	testCases := []testCase{
		{
			name:    "success string code",
			body:    []byte(`{"code":"10000","message":"ok"}`),
			traceID: "trace-ok",
		},
		{
			name:    "success response with 5xx trace",
			body:    []byte(`{"code":"10000","message":"ok"}`),
			traceID: "trace-5xx",
		},
		{
			name:      "panic recovered",
			body:      []byte(`{"code":50001,"message":"panic recovered"}`),
			traceID:   "trace-panic",
			errorCode: 50001,
		},
		{
			name:      "internal error string code",
			body:      []byte(`{"code":"50001","message":"internal error"}`),
			traceID:   "trace-biz",
			errorCode: 50001,
		},
		{
			name:      "auth error",
			body:      []byte(`{"code":"30002","message":"auth invalid"}`),
			traceID:   "trace-auth",
			errorCode: 30002,
		},
		{
			name:      "parameter error",
			body:      []byte(`{"code":"20001","message":"param error"}`),
			traceID:   "trace-param",
			errorCode: 20001,
		},
		{
			name:      "business error",
			body:      []byte(`{"code":"40001","message":"biz error"}`),
			traceID:   "trace-biz-4xx",
			errorCode: 40001,
		},
		{
			name:      "internal error numeric code",
			body:      []byte(`{"code":50001,"message":"internal error"}`),
			traceID:   "trace-biz-int",
			errorCode: 50001,
		},
		{
			name:    "paper response",
			body:    []byte(`{"code":2000,"msg":"Success"}`),
			traceID: "trace-paper",
		},
	}

	now := time.Now()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := buildRequestEvent("/api/foo", tc.body, tc.traceID, now)

			assert.Equal(t, "/api/foo", event.route)
			assert.Equal(t, tc.errorCode, event.errorCode)
			assert.Equal(t, tc.traceID, event.traceID)
			assert.Equal(t, now, event.timestamp)
		})
	}
}

func TestCompactWindow(t *testing.T) {
	type testCase struct {
		name          string
		events        []requestEvent
		cutoff        time.Time
		expectedRoute string
	}

	now := time.Date(2026, time.July, 30, 12, 0, 0, 0, time.UTC)
	testCases := []testCase{
		{
			name: "remove events before cutoff",
			events: []requestEvent{
				{route: "/expired", timestamp: now.Add(-2 * time.Minute)},
				{route: "/kept", timestamp: now.Add(-30 * time.Second)},
			},
			cutoff:        now.Add(-time.Minute),
			expectedRoute: "/kept",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kept := compactWindow(tc.events, tc.cutoff)

			if assert.Len(t, kept, 1) {
				assert.Equal(t, tc.expectedRoute, kept[0].route)
			}
		})
	}
}

func TestAggregateRouteStats(t *testing.T) {
	testCases := []struct {
		name     string
		events   []requestEvent
		expected map[string]routeStat
	}{
		{
			name: "aggregate requests and reportable errors",
			events: []requestEvent{
				{route: "/api/foo", traceID: "trace-1"},
				{route: "/api/foo", errorCode: 50001, traceID: "trace-2"},
				{route: "/api/foo", errorCode: 30002, traceID: "trace-ignored"},
				{route: "/api/bar", errorCode: 40001, traceID: "trace-3"},
			},
			expected: map[string]routeStat{
				"/api/foo": {
					requests:  3,
					errors:    2,
					errorRate: 0.6667,
					traceID:   "trace-2",
					errorCode: 50001,
				},
				"/api/bar": {
					requests:  1,
					errors:    1,
					errorRate: 1,
					traceID:   "trace-3",
					errorCode: 40001,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stats := aggregateRouteStats(tc.events)

			for route, expected := range tc.expected {
				actual, ok := stats[route]
				if assert.True(t, ok) {
					assert.Equal(t, expected.requests, actual.requests)
					assert.Equal(t, expected.errors, actual.errors)
					assert.InDelta(t, expected.errorRate, actual.errorRate, 0.0001)
					assert.Equal(t, expected.traceID, actual.traceID)
					assert.Equal(t, expected.errorCode, actual.errorCode)
				}
			}
		})
	}
}

func TestMonitorRecordSkipsDisabledAndBlacklisted(t *testing.T) {
	testCases := []struct {
		name          string
		config        MonitorConfig
		events        []requestEvent
		expectedRoute string
		expectedCount int
	}{
		{
			name: "disabled monitor skips events",
			events: []requestEvent{
				{route: "/api/foo"},
			},
			expectedCount: 0,
		},
		{
			name: "blacklisted route is skipped",
			config: MonitorConfig{
				Enabled:   true,
				Blacklist: map[string]struct{}{"/api/foo": {}},
			},
			events: []requestEvent{
				{route: "/api/foo"},
				{route: "/api/bar"},
			},
			expectedRoute: "/api/bar",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			monitor := newAPIMonitor(tc.config)
			for _, event := range tc.events {
				monitor.record(event)
			}

			assert.Len(t, monitor.events, tc.expectedCount)
			if tc.expectedCount > 0 {
				assert.Equal(t, tc.expectedRoute, monitor.events[0].route)
			}
		})
	}
}

func TestStartAPIMonitorStopReturnsBeforeNextCheck(t *testing.T) {
	testCases := []struct {
		name          string
		checkInterval time.Duration
	}{
		{
			name:          "stop returns before next check",
			checkInterval: time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiMonitorInstance = nil
			apiMonitorStartOnce = sync.Once{}
			apiMonitorStop = func(context.Context) {}
			defer func() {
				apiMonitorInstance = nil
				apiMonitorStartOnce = sync.Once{}
				apiMonitorStop = func(context.Context) {}
			}()

			stop := StartAPIMonitor(MonitorConfig{
				Enabled:       true,
				CheckInterval: tc.checkInterval,
			})

			done := make(chan struct{})
			go func() {
				stop(context.Background())
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("api monitor stop should not wait for the next check interval")
			}
		})
	}
}

func TestMonitorAlertCooldownAndRecover(t *testing.T) {
	testCases := []struct {
		name       string
		config     MonitorConfig
		stat       routeStat
		firstCheck time.Time
	}{
		{
			name: "cooldown prevents duplicate alert and recovery clears state",
			config: MonitorConfig{
				Enabled:       true,
				Window:        time.Minute,
				CheckInterval: time.Second,
				Threshold:     0.5,
				MinRequests:   2,
				Cooldown:      10 * time.Minute,
			},
			stat: routeStat{
				requests:  2,
				errors:    1,
				errorRate: 0.5,
				traceID:   "trace-alert",
				errorCode: 50001,
			},
			firstCheck: time.Date(2026, time.July, 30, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			monitor := newAPIMonitor(tc.config)

			monitor.checkRoute(tc.firstCheck, "/api/foo", tc.stat)
			firstAlert := monitor.alerts["/api/foo"]
			assert.True(t, firstAlert.firing)
			assert.Equal(t, "trace-alert", firstAlert.lastTrace)
			assert.Equal(t, int64(50001), firstAlert.lastCode)

			tc.stat.traceID = "trace-cooldown"
			monitor.checkRoute(tc.firstCheck.Add(time.Minute), "/api/foo", tc.stat)
			assert.Equal(t, firstAlert.lastAlert, monitor.alerts["/api/foo"].lastAlert)
			assert.Equal(t, "trace-alert", monitor.alerts["/api/foo"].lastTrace)

			tc.stat.errorRate = 0.1
			monitor.checkRoute(tc.firstCheck.Add(2*time.Minute), "/api/foo", tc.stat)
			_, ok := monitor.alerts["/api/foo"]
			assert.False(t, ok)
		})
	}
}
