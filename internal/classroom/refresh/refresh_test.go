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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/jwch"
)

type sourceFixture struct {
	lists           map[string]int
	requests        []jwch.ClassroomUsageReq
	failList        bool
	listAlwaysFails bool
	alwaysFail      bool
	allFail         bool
	roomErrors      bool
}

func (s *sourceFixture) GetClassroomBuildings(_ context.Context, campus string) ([]string, error) {
	s.lists[campus]++
	if s.listAlwaysFails || (s.failList && s.lists[campus] == 1) {
		return nil, errors.New("list unavailable")
	}
	return []string{"楼A", "楼B", "楼C", "楼A", ""}, nil
}

func (s *sourceFixture) GetClassroomUsage(_ context.Context, req jwch.ClassroomUsageReq) (*jwch.ClassroomUsage, error) {
	s.requests = append(s.requests, req)
	count := 0
	for _, previous := range s.requests {
		if previous.Date == req.Date && previous.Campus == req.Campus && previous.Buildings[0] == req.Buildings[0] {
			count++
		}
	}
	if s.allFail || (req.Buildings[0] == "楼B" && (count == 1 || s.alwaysFail)) {
		if s.roomErrors {
			return &jwch.ClassroomUsage{
				Date: req.Date, Term: "202601", Rooms: []*jwch.ClassroomUsageRoom{{
					Building: req.Buildings[0], Error: errors.New("building unavailable"),
				}},
			}, nil
		}
		return nil, errors.New("building unavailable")
	}
	return &jwch.ClassroomUsage{
		Date: req.Date, Term: "202601", Rooms: []*jwch.ClassroomUsageRoom{
			{Building: req.Buildings[0], Location: req.Campus + req.Buildings[0] + "-101", Occupancies: []*jwch.ClassroomUsageOccupancy{}},
		},
	}, nil
}

type cacheFixture struct {
	days           map[string]*classroomCache.DayAvailability
	failures       map[string][]string
	writes         int
	failFirstWrite bool
	writeErr       error
	markErr        error
}

func (c *cacheFixture) SetAvailability(_ context.Context, campus string, day *classroomCache.DayAvailability, failed ...string) error {
	c.writes++
	if c.writeErr != nil {
		return c.writeErr
	}
	if c.failFirstWrite && c.writes == 1 {
		return errors.New("redis temporarily unavailable")
	}
	key := day.Date + campus
	c.days[key] = day
	c.failures[key] = append([]string(nil), failed...)
	return nil
}

func (c *cacheFixture) MarkAvailabilityFailure(_ context.Context, date, campus string, failed ...string) error {
	if c.markErr != nil {
		return c.markErr
	}
	c.failures[date+campus] = append([]string{}, failed...)
	return nil
}

func fixtures() (*sourceFixture, *cacheFixture, *Refresher) {
	source := &sourceFixture{lists: make(map[string]int)}
	cache := &cacheFixture{days: make(map[string]*classroomCache.DayAvailability), failures: make(map[string][]string)}
	return source, cache, &Refresher{Source: source, Cache: cache}
}

func TestRunRetriesOnlyFailedBuildings(t *testing.T) {
	source, cache, runner := fixtures()
	dates := []string{"2026-10-17", "2026-10-18"}
	campuses := []string{"旗山校区", "铜盘校区"}
	require.NoError(t, runner.Run(context.Background(), dates, campuses))
	require.Equal(t, map[string]int{"旗山校区": 1, "铜盘校区": 1}, source.lists)
	// 四个日期校区组合的首轮查询全部结束后，才补抓各自失败的楼B。
	require.Len(t, source.requests, 16)
	for _, req := range source.requests[:12] {
		require.NotEmpty(t, req.Buildings[0])
	}
	for _, req := range source.requests[12:] {
		require.Equal(t, "楼B", req.Buildings[0])
		require.Equal(t, "202601", req.Term)
	}
	require.Len(t, cache.days, 4)
	for key, day := range cache.days {
		require.Len(t, day.Rooms, 3)
		require.Empty(t, cache.failures[key])
	}
}

func TestRunPersistsSuccessfulBuildingsAndFailureNames(t *testing.T) {
	source, cache, runner := fixtures()
	source.alwaysFail = true
	err := runner.Run(context.Background(), []string{"2026-10-17"}, []string{"旗山校区"})
	require.NoError(t, err, "已记录的上游失败不应触发任务队列全量重试")
	require.Len(t, source.requests, 4)
	require.Len(t, cache.days["2026-10-17旗山校区"].Rooms, 2)
	require.Equal(t, []string{"楼B"}, cache.failures["2026-10-17旗山校区"])
}

func TestRunRetriesBuildingListOncePerCampus(t *testing.T) {
	source, cache, runner := fixtures()
	source.failList = true
	require.NoError(t, runner.Run(context.Background(), []string{"2026-10-17", "2026-10-18"}, []string{"旗山校区", "铜盘校区"}))
	require.Equal(t, map[string]int{"旗山校区": 2, "铜盘校区": 2}, source.lists)
	// 列表补抓成功后首次查询的楼也允许失败一次、补抓一次。
	require.Len(t, cache.days, 4)
}

func TestRunListFailureAndAllBuildingFailure(t *testing.T) {
	for _, listFailure := range []bool{true, false} {
		t.Run(fmt.Sprint(listFailure), func(t *testing.T) {
			source, cache, runner := fixtures()
			source.listAlwaysFails = listFailure
			source.allFail = !listFailure
			require.NoError(t, runner.Run(context.Background(), []string{"2026-10-17"}, []string{"旗山校区"}))
			require.Empty(t, cache.days)
			require.Contains(t, cache.failures, "2026-10-17旗山校区")
			if listFailure {
				require.Equal(t, 2, source.lists["旗山校区"])
				require.Empty(t, source.requests)
			} else {
				require.Len(t, source.requests, 6)
				require.Equal(t, []string{"楼A", "楼B", "楼C"}, cache.failures["2026-10-17旗山校区"])
			}
		})
	}
}

func TestRunRetriesCacheWithoutRefetchingSuccessfulBuildings(t *testing.T) {
	source, cache, runner := fixtures()
	// 预先满足楼B的一次失败，使本轮三个楼均在第一次抓取成功。
	source.requests = append(source.requests, jwch.ClassroomUsageReq{Date: "2026-10-17", Campus: "旗山校区", Buildings: []string{"楼B"}})
	cache.failFirstWrite = true
	require.NoError(t, runner.Run(context.Background(), []string{"2026-10-17"}, []string{"旗山校区"}))
	require.Len(t, source.requests, 4)
	require.Equal(t, 2, cache.writes)
}

func TestRunCancellation(t *testing.T) {
	source, cache, runner := fixtures()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, runner.Run(ctx, []string{"2026-10-17"}, []string{"旗山校区"}), context.Canceled)
	require.Empty(t, source.lists)
	require.Empty(t, source.requests)
	require.Empty(t, cache.days)
}

// 发布版 jwch 将单楼错误放在 Rooms 中，顶层 error 仍然可能为 nil。
func TestRunHandlesJwchRoomErrors(t *testing.T) {
	for _, persistent := range []bool{false, true} {
		t.Run(fmt.Sprint(persistent), func(t *testing.T) {
			source, cache, runner := fixtures()
			source.roomErrors = true
			source.alwaysFail = persistent
			err := runner.Run(context.Background(), []string{"2026-11-07"}, []string{"旗山校区"})
			if persistent {
				require.NoError(t, err, "已记录的上游失败不应触发任务队列全量重试")
				require.Equal(t, []string{"楼B"}, cache.failures["2026-11-07旗山校区"])
				require.Len(t, cache.days["2026-11-07旗山校区"].Rooms, 2)
			} else {
				require.NoError(t, err)
				require.Empty(t, cache.failures["2026-11-07旗山校区"])
				require.Len(t, cache.days["2026-11-07旗山校区"].Rooms, 3)
			}
			require.Len(t, source.requests, 4)
			for _, req := range source.requests {
				require.Len(t, req.Buildings, 1, "每次仅查询一栋楼")
			}
			require.Equal(t, []string{"楼B"}, source.requests[3].Buildings)
		})
	}
}

func TestValidateUsage(t *testing.T) {
	item := &target{date: "2026-11-07", campus: "旗山校区", term: "202601"}
	for _, tc := range []struct {
		name  string
		usage *jwch.ClassroomUsage
		valid bool
	}{
		{name: "nil response"},
		{name: "wrong date", usage: &jwch.ClassroomUsage{Date: "2026-11-08", Term: item.term, Rooms: []*jwch.ClassroomUsageRoom{}}},
		{name: "term changed", usage: &jwch.ClassroomUsage{Date: item.date, Term: "202602", Rooms: []*jwch.ClassroomUsageRoom{}}},
		{name: "nil rooms", usage: &jwch.ClassroomUsage{Date: item.date, Term: item.term}},
		{name: "empty building", usage: &jwch.ClassroomUsage{Date: item.date, Term: item.term, Rooms: []*jwch.ClassroomUsageRoom{}}, valid: true},
		{name: "nil room", usage: &jwch.ClassroomUsage{Date: item.date, Term: item.term, Rooms: []*jwch.ClassroomUsageRoom{nil}}},
		{name: "wrong building", usage: &jwch.ClassroomUsage{Date: item.date, Term: item.term, Rooms: []*jwch.ClassroomUsageRoom{
			{Building: "楼B", Location: "教室101"},
		}}},
		{name: "duplicate room", usage: &jwch.ClassroomUsage{Date: item.date, Term: item.term, Rooms: []*jwch.ClassroomUsageRoom{
			{Building: "楼A", Location: "教室101"}, {Building: "楼A", Location: "教室101"},
		}}},
		{name: "invalid occupancy", usage: &jwch.ClassroomUsage{Date: item.date, Term: item.term, Rooms: []*jwch.ClassroomUsageRoom{
			{Building: "楼A", Location: "教室101", Occupancies: []*jwch.ClassroomUsageOccupancy{{Start: 0, End: 12}}},
		}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateUsage(item, tc.usage, "楼A")
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestRunReturnsCacheFailures(t *testing.T) {
	redisErr := errors.New("redis unavailable")
	for _, name := range []string{"snapshot write", "partial snapshot write", "list failure marker", "all buildings failure marker"} {
		t.Run(name, func(t *testing.T) {
			source, cache, runner := fixtures()
			switch name {
			case "snapshot write":
				cache.writeErr = redisErr
			case "partial snapshot write":
				source.alwaysFail = true
				cache.writeErr = redisErr
			case "list failure marker":
				source.listAlwaysFails = true
				cache.markErr = redisErr
			case "all buildings failure marker":
				source.allFail = true
				cache.markErr = redisErr
			}
			err := runner.Run(context.Background(), []string{"2026-11-13"}, []string{"旗山校区"})
			require.ErrorIs(t, err, redisErr, "缓存失败必须继续传递给任务队列")
		})
	}
}

func TestRunPersistentFailureCompletesBatch(t *testing.T) {
	source, cache, runner := fixtures()
	source.alwaysFail = true
	dates := []string{"2026-11-13", "2026-11-14"}
	campuses := []string{"旗山校区", "铜盘校区"}
	require.NoError(t, runner.Run(context.Background(), dates, campuses))
	require.Len(t, source.requests, 16, "每个日期校区查询三个楼，再仅补抓失败楼一次")
	require.Len(t, cache.days, 4)
	for key, day := range cache.days {
		require.Len(t, day.Rooms, 2)
		require.Equal(t, []string{"楼B"}, cache.failures[key])
	}
}
