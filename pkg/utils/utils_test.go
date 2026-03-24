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

	"github.com/west2-online/jwch"
)

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

func TestFindTermByDate(t *testing.T) {
	terms := []jwch.CalTerm{
		{
			TermId:     "202401",
			SchoolYear: "2024-2025",
			Term:       "1",
			StartDate:  "2024-08-26",
			EndDate:    "2025-01-17",
		},
		{
			TermId:     "202402",
			SchoolYear: "2024-2025",
			Term:       "2",
			StartDate:  "2025-03-03",
			EndDate:    "2025-07-11",
		},
	}

	tests := []struct {
		name       string
		terms      []jwch.CalTerm
		date       time.Time
		wantTermId string
		wantFound  bool
	}{
		{
			name:       "DateInFirstTerm",
			terms:      terms,
			date:       time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			wantTermId: "202401",
			wantFound:  true,
		},
		{
			name:       "DateOnStartBoundary",
			terms:      terms,
			date:       time.Date(2024, 8, 26, 0, 0, 0, 0, time.UTC),
			wantTermId: "202401",
			wantFound:  true,
		},
		{
			name:       "DateOnEndBoundary",
			terms:      terms,
			date:       time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC),
			wantTermId: "202401",
			wantFound:  true,
		},
		{
			name:       "DateInSecondTerm",
			terms:      terms,
			date:       time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC),
			wantTermId: "202402",
			wantFound:  true,
		},
		{
			name:      "DateBetweenTerms",
			terms:     terms,
			date:      time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			wantFound: false,
		},
		{
			name:      "DateBeforeAllTerms",
			terms:     terms,
			date:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantFound: false,
		},
		{
			name:      "DateAfterAllTerms",
			terms:     terms,
			date:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			wantFound: false,
		},
		{
			name:      "EmptyTerms",
			terms:     []jwch.CalTerm{},
			date:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			wantFound: false,
		},
		{
			name: "InvalidStartDate",
			terms: []jwch.CalTerm{
				{TermId: "bad", StartDate: "not-a-date", EndDate: "2025-01-17"},
			},
			date:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			wantFound: false,
		},
		{
			name: "InvalidEndDate",
			terms: []jwch.CalTerm{
				{TermId: "bad", StartDate: "2024-08-26", EndDate: "not-a-date"},
			},
			date:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := FindTermByDate(tt.terms, tt.date)
			if found != tt.wantFound {
				t.Errorf("FindTermByDate() found = %v, want %v", found, tt.wantFound)
				return
			}
			if tt.wantFound && got.TermId != tt.wantTermId {
				t.Errorf("FindTermByDate() TermId = %v, want %v", got.TermId, tt.wantTermId)
			}
		})
	}
}
