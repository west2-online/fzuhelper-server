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

package refresh

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextDelay(t *testing.T) {
	loc := time.FixedZone("CST", 8*3600)
	for _, tc := range []struct {
		name string
		now  time.Time
		want time.Duration
	}{
		{name: "before first refresh", now: time.Date(2026, 11, 7, 4, 0, 0, 0, loc), want: time.Hour},
		{name: "at first refresh", now: time.Date(2026, 11, 7, 5, 0, 0, 0, loc), want: 8 * time.Hour},
		{name: "between refreshes", now: time.Date(2026, 11, 7, 12, 59, 0, 0, loc), want: time.Minute},
		{name: "at second refresh", now: time.Date(2026, 11, 7, 13, 0, 0, 0, loc), want: 16 * time.Hour},
		{name: "after second refresh", now: time.Date(2026, 11, 7, 23, 0, 0, 0, loc), want: 6 * time.Hour},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, NextDelay(tc.now))
		})
	}
}

func TestDates(t *testing.T) {
	loc := time.FixedZone("CST", 8*3600)
	for _, tc := range []struct {
		name  string
		now   time.Time
		first string
		last  string
		count int
	}{
		// 当周一到下周日共 14 天
		{name: "monday", now: time.Date(2026, 11, 2, 9, 30, 0, 0, loc), first: "2026-11-02", last: "2026-11-15", count: 14},
		{name: "saturday", now: time.Date(2026, 11, 7, 23, 0, 0, 0, loc), first: "2026-11-07", last: "2026-11-15", count: 9},
		{name: "sunday", now: time.Date(2026, 11, 8, 5, 0, 0, 0, loc), first: "2026-11-08", last: "2026-11-15", count: 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dates := Dates(tc.now)
			require.Len(t, dates, tc.count)
			assert.Equal(t, tc.first, dates[0].Format("2006-01-02"))
			assert.Equal(t, tc.last, dates[len(dates)-1].Format("2006-01-02"))
			for i, day := range dates {
				assert.Equal(t, 0, day.Hour(), "日期应归一到当天 0 点")
				if i > 0 {
					assert.Equal(t, 24*time.Hour, day.Sub(dates[i-1]), "日期应连续")
				}
			}
		})
	}
}
