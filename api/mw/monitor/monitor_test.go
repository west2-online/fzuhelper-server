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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildRequestEvent(t *testing.T) {
	now := time.Now()

	event := buildRequestEvent("/api/foo", []byte(`{"code":"10000","message":"ok"}`), "trace-ok", now)
	assert.Equal(t, "/api/foo", event.route)
	assert.Zero(t, event.errorCode)
	assert.Equal(t, "trace-ok", event.traceID)
	assert.Equal(t, now, event.timestamp)

	event = buildRequestEvent("/api/foo", []byte(`{"code":"10000","message":"ok"}`), "trace-5xx", now)
	assert.Zero(t, event.errorCode)

	event = buildRequestEvent("/api/foo", []byte(`{"code":50001,"message":"panic recovered"}`), "trace-panic", now)
	assert.Equal(t, int64(50001), event.errorCode)
	assert.Equal(t, "trace-panic", event.traceID)

	event = buildRequestEvent("/api/foo", []byte(`{"code":"50001","message":"internal error"}`), "trace-biz", now)
	assert.Equal(t, int64(50001), event.errorCode)
	assert.Equal(t, "trace-biz", event.traceID)

	event = buildRequestEvent("/api/foo", []byte(`{"code":"30002","message":"auth invalid"}`), "trace-auth", now)
	assert.Equal(t, int64(30002), event.errorCode)

	event = buildRequestEvent("/api/foo", []byte(`{"code":"20001","message":"param error"}`), "trace-param", now)
	assert.Equal(t, int64(20001), event.errorCode)

	event = buildRequestEvent("/api/foo", []byte(`{"code":"40001","message":"biz error"}`), "trace-biz-4xx", now)
	assert.Equal(t, int64(40001), event.errorCode)

	event = buildRequestEvent("/api/foo", []byte(`{"code":50001,"message":"internal error"}`), "trace-biz-int", now)
	assert.Equal(t, int64(50001), event.errorCode)

	event = buildRequestEvent("/api/foo", []byte(`{"code":2000,"msg":"Success"}`), "trace-paper", now)
	assert.Zero(t, event.errorCode)

	event = buildRequestEvent("/api/foo", []byte(`{"code":200,"message":"ok"}`), "trace-custom", now)
	assert.Zero(t, event.errorCode)
}

func TestCompactWindow(t *testing.T) {
	now := time.Now()
	events := []requestEvent{
		{route: "/expired", timestamp: now.Add(-2 * time.Minute)},
		{route: "/kept", timestamp: now.Add(-30 * time.Second)},
	}

	kept := compactWindow(events, now.Add(-time.Minute))
	assert.Len(t, kept, 1)
	assert.Equal(t, "/kept", kept[0].route)
}

func TestAggregateRouteStats(t *testing.T) {
	events := []requestEvent{
		{route: "/api/foo", traceID: "trace-1"},
		{route: "/api/foo", errorCode: 50001, traceID: "trace-2"},
		{route: "/api/foo", errorCode: 30002, traceID: "trace-ignored"},
		{route: "/api/bar", errorCode: 40001, traceID: "trace-3"},
	}

	stats := aggregateRouteStats(events)

	assert.Equal(t, int64(3), stats["/api/foo"].requests)
	assert.Equal(t, int64(2), stats["/api/foo"].errors)
	assert.InDelta(t, 0.6667, stats["/api/foo"].errorRate, 0.0001)
	assert.Equal(t, "trace-2", stats["/api/foo"].traceID)
	assert.Equal(t, int64(50001), stats["/api/foo"].errorCode)
	assert.Equal(t, int64(1), stats["/api/bar"].requests)
	assert.Equal(t, int64(1), stats["/api/bar"].errors)
	assert.Equal(t, "trace-3", stats["/api/bar"].traceID)
	assert.Equal(t, int64(40001), stats["/api/bar"].errorCode)
}

func TestMonitorRecordSkipsDisabledAndBlacklisted(t *testing.T) {
	disabled := newAPIMonitor(MonitorConfig{})
	disabled.record(requestEvent{route: "/api/foo"})
	assert.Empty(t, disabled.events)

	enabled := newAPIMonitor(MonitorConfig{
		Enabled:   true,
		Blacklist: map[string]struct{}{"/api/foo": {}},
	})
	enabled.record(requestEvent{route: "/api/foo"})
	enabled.record(requestEvent{route: "/api/bar"})

	assert.Len(t, enabled.events, 1)
	assert.Equal(t, "/api/bar", enabled.events[0].route)
}

func TestMonitorAlertCooldownAndRecover(t *testing.T) {
	monitor := newAPIMonitor(MonitorConfig{
		Enabled:       true,
		Window:        time.Minute,
		CheckInterval: time.Second,
		Threshold:     0.5,
		MinRequests:   2,
		Cooldown:      10 * time.Minute,
	})
	now := time.Now()
	stat := routeStat{
		requests:  2,
		errors:    1,
		errorRate: 0.5,
		traceID:   "trace-alert",
		errorCode: 50001,
	}

	monitor.checkRoute(now, "/api/foo", stat)
	firstAlert := monitor.alerts["/api/foo"]
	assert.True(t, firstAlert.firing)
	assert.Equal(t, "trace-alert", firstAlert.lastTrace)
	assert.Equal(t, int64(50001), firstAlert.lastCode)

	stat.traceID = "trace-cooldown"
	monitor.checkRoute(now.Add(time.Minute), "/api/foo", stat)
	assert.Equal(t, firstAlert.lastAlert, monitor.alerts["/api/foo"].lastAlert)
	assert.Equal(t, "trace-alert", monitor.alerts["/api/foo"].lastTrace)

	stat.errorRate = 0.1
	monitor.checkRoute(now.Add(2*time.Minute), "/api/foo", stat)
	_, ok := monitor.alerts["/api/foo"]
	assert.False(t, ok)
}
