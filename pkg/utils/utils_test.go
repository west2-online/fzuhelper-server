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

package utils

import (
	"testing"
	"time"
)

func TestTimeParse(t *testing.T) {
	tests := []struct {
		name     string
		time     string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "Success",
			time:     "2026-03-02",
			wantTime: time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:    "Error",
			time:    "2024~03~05",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultTime, err := TimeParse(tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resultTime != tt.wantTime {
				t.Errorf("TimeParse() time = %v, want %v", resultTime, tt.wantTime)
			}
		})
	}
}

func TestGetWeekdayByDate(t *testing.T) {
	tests := []struct {
		name          string
		termStartDate string
		date          string
		wantWeek      int
		wantWeekday   int
		wantErr       bool
	}{
		{
			name:          "Normal",
			termStartDate: "2026-03-02",
			date:          "2026-03-17",
			wantWeek:      3,
			wantWeekday:   2,
			wantErr:       false,
		},
		{
			name:          "StartDateIsNotMonday",
			termStartDate: "2026-03-08",
			date:          "2026-03-17",
			wantWeek:      3,
			wantWeekday:   2,
			wantErr:       false,
		},
		{
			name:          "DateBeforeTermStart",
			termStartDate: "2024-03-05",
			date:          "2024-03-01",
			wantWeekday:   0,
			wantWeek:      0,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			week, weekday, err := GetWeekdayByDate(tt.termStartDate, tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeekdayByDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if week != tt.wantWeek {
				t.Errorf("GetWeekdayByDate() week = %v, want %v", week, tt.wantWeek)
			}
			if weekday != tt.wantWeekday {
				t.Errorf("GetWeekdayByDate() weekday = %v, want %v", weekday, tt.wantWeekday)
			}
		})
	}
}
