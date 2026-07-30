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
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// requestEvent 是滑动窗口计算保留的最小请求事件。
type requestEvent struct {
	route     string
	errorCode int64
	traceID   string
	timestamp time.Time
}

// routeStat 是按路由聚合后的报警判断数据。
type routeStat struct {
	requests  int64
	errors    int64
	errorRate float64
	traceID   string
	errorCode int64
}

// alertState 记录单个路由当前的报警状态。
type alertState struct {
	firing     bool
	lastAlert  time.Time
	lastTrace  string
	lastCode   int64
	lastErrors int64
}

// apiMonitor 维护滑动窗口事件和路由报警状态。
type apiMonitor struct {
	mu     sync.Mutex
	cfg    MonitorConfig
	events []requestEvent
	alerts map[string]alertState
}

func newAPIMonitor(cfg MonitorConfig) *apiMonitor {
	if cfg.Blacklist == nil {
		cfg.Blacklist = make(map[string]struct{})
	}
	return &apiMonitor{
		cfg:    cfg,
		events: make([]requestEvent, 0),
		alerts: make(map[string]alertState),
	}
}

func (m *apiMonitor) enabled() bool {
	return m.cfg.Enabled
}

func (m *apiMonitor) checkInterval() time.Duration {
	return m.cfg.CheckInterval
}

func (m *apiMonitor) shouldIgnore(route string) bool {
	_, ok := m.cfg.Blacklist[route]
	return ok
}

// 追加到滑动窗口维护中枢
func (m *apiMonitor) record(event requestEvent) {
	if !m.cfg.Enabled || m.shouldIgnore(event.route) {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

// 删除过期事件,并且判断是否报错
func (m *apiMonitor) check() {
	if !m.cfg.Enabled {
		return
	}

	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := now.Add(-m.cfg.Window)
	m.events = compactWindow(m.events, cutoff)
	stats := aggregateRouteStats(m.events)

	for route, stat := range stats {
		m.checkRoute(now, route, stat)
	}
	for route, alert := range m.alerts {
		if _, ok := stats[route]; !ok && alert.firing {
			m.logRecovered(route, routeStat{})
			delete(m.alerts, route)
		}
	}
}

// 判断是否报错
func (m *apiMonitor) checkRoute(now time.Time, route string, stat routeStat) {
	if stat.requests < m.cfg.MinRequests {
		return
	}

	alert := m.alerts[route]
	if stat.errorRate >= m.cfg.Threshold {
		if !alert.firing || now.Sub(alert.lastAlert) >= m.cfg.Cooldown {
			alert.firing = true
			alert.lastAlert = now
			alert.lastTrace = stat.traceID
			alert.lastCode = stat.errorCode
			alert.lastErrors = stat.errors
			m.alerts[route] = alert
			m.logAlert(route, stat)
		}
		return
	}

	if alert.firing {
		m.logRecovered(route, stat)
		delete(m.alerts, route)
	}
}

// 生成事件
func buildRequestEvent(route string, responseBody []byte, traceID string, now time.Time) requestEvent {
	code, ok := responseCode(responseBody)
	if !ok || isSuccessCode(code) {
		code = 0
	}

	return requestEvent{
		route:     route,
		errorCode: code,
		traceID:   traceID,
		timestamp: now,
	}
}

// 删除过期时间的事件
func compactWindow(events []requestEvent, cutoff time.Time) []requestEvent {
	kept := events[:0]
	for _, event := range events {
		if !event.timestamp.Before(cutoff) {
			kept = append(kept, event)
		}
	}
	return kept
}

// 计算错误率并生成报错这个事件
func aggregateRouteStats(events []requestEvent) map[string]routeStat {
	stats := make(map[string]routeStat)
	for _, event := range events {
		stat := stats[event.route]
		stat.requests++
		if isMonitorError(event.errorCode) {
			stat.errors++
			if stat.traceID == "" {
				stat.traceID = event.traceID
				stat.errorCode = event.errorCode
			}
		}
		stats[event.route] = stat
	}

	for route, stat := range stats {
		stat.errorRate = float64(stat.errors) / float64(stat.requests)
		stats[route] = stat
	}
	return stats
}

func isMonitorError(code int64) bool {
	return code != 0 && !isSuccessCode(code)
}

func isSuccessCode(code int64) bool {
	switch code {
	case errno.SuccessCode, errno.SuccessCodePaper:
		return true
	default:
		return false
	}
}

func responseCode(responseBody []byte) (int64, bool) {
	if len(responseBody) == 0 {
		return 0, false
	}

	var payload struct {
		Code any `json:"code"`
	}
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return 0, false
	}

	switch code := payload.Code.(type) {
	case string:
		value, err := strconv.ParseInt(code, 10, 64)
		if err != nil {
			return 0, false
		}
		return value, true
	case float64:
		return int64(code), true
	default:
		return 0, false
	}
}

func (m *apiMonitor) logAlert(route string, stat routeStat) {
	logger.Error("api service anomaly detected",
		m.logFields("api_service_anomaly", route, stat)...,
	)
}

func (m *apiMonitor) logRecovered(route string, stat routeStat) {
	logger.Info("api service anomaly recovered",
		m.logFields("api_service_anomaly_recovered", route, stat)...,
	)
}

func (m *apiMonitor) logFields(event string, route string, stat routeStat) []zap.Field {
	return []zap.Field{
		zap.String("event", event),
		zap.String("route", route),
		zap.String("traceid", stat.traceID),
		zap.Int64("error_code", stat.errorCode),
		zap.Int64("requests", stat.requests),
		zap.Int64("errors", stat.errors),
		zap.Float64("error_rate", stat.errorRate),
		zap.Float64("threshold", m.cfg.Threshold),
		zap.Int64("min_requests", m.cfg.MinRequests),
	}
}
