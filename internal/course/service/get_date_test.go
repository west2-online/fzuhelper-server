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

package service

import (
	"testing"
	"time"
)

func TestISOWeekCrossWeek(t *testing.T) {
	// 创建一个周日的时间（2025-03-23）
	sunday := time.Date(2025, 3, 23, 0, 0, 0, 0, time.UTC)
	sundayYear, sundayWeek := sunday.ISOWeek()

	// 创建一个周一的时间（2025-03-25）
	monday := time.Date(2025, 3, 24, 0, 0, 0, 0, time.UTC)
	mondayYear, mondayWeek := monday.ISOWeek()

	// 验证周日和周一的周数是否不同
	if sundayWeek == mondayWeek {
		t.Errorf("Expected different weeks for Sunday and Monday, but got same week: Sunday(Week %d) Monday(Week %d)",
			sundayWeek, mondayWeek)
	}

	t.Logf("Sunday: Year %d Week %d", sundayYear, sundayWeek)
	t.Logf("Monday: Year %d Week %d", mondayYear, mondayWeek)
}
