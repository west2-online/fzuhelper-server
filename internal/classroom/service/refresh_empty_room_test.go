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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

type sourceState struct {
	lists           map[string]int
	requests        []jwch.ClassroomUsageReq
	calls           map[string]int
	failList        bool
	listAlwaysFails bool
	failBuilding    bool
	alwaysFail      bool
	allFail         bool
	topLevelError   bool
	empty           bool
}

type cacheState struct {
	days          map[string][]*model.Classroom
	failures      map[string]bool
	writes        int
	writeErr      error
	failureErr    error
	failureWrites int
}

func fixtures(t *testing.T) (*sourceState, *cacheState, *ClassroomService) {
	t.Helper()
	s := &sourceState{lists: make(map[string]int), calls: make(map[string]int), failBuilding: true}
	c := &cacheState{days: make(map[string][]*model.Classroom), failures: make(map[string]bool)}
	buildings := mockey.Mock((*jwch.Student).GetClassroomBuildings).To(func(_ *jwch.Student, _ context.Context, campus string) ([]string, error) {
		s.lists[campus]++
		if s.listAlwaysFails || (s.failList && s.lists[campus] == 1) {
			return nil, errors.New("list unavailable")
		}
		return []string{"楼A", "楼B", "楼C"}, nil
	}).Build()
	t.Cleanup(func() { buildings.UnPatch() })
	usage := mockey.Mock((*jwch.Student).GetClassroomUsage).To(func(_ *jwch.Student, _ context.Context, req jwch.ClassroomUsageReq) (*jwch.ClassroomUsage, error) {
		s.requests = append(s.requests, req)
		result := &jwch.ClassroomUsage{Date: req.Date, Term: "202601"}
		if s.empty {
			return result, nil
		}
		for _, building := range req.Buildings {
			key := req.Date + req.Campus + building
			s.calls[key]++
			if s.allFail || (s.failBuilding && building == "楼B" && (s.calls[key] == 1 || s.alwaysFail)) {
				if s.topLevelError {
					return nil, errors.New("page unavailable")
				}
				result.Rooms = append(result.Rooms, &jwch.ClassroomUsageRoom{Building: building, Error: errors.New("building unavailable")})
				continue
			}
			result.Rooms = append(result.Rooms, &jwch.ClassroomUsageRoom{Building: building, Location: req.Campus + building + "-101"})
		}
		return result, nil
	}).Build()
	t.Cleanup(func() { usage.UnPatch() })
	set := mockey.Mock((*classroomCache.CacheClassroom).SetAvailabilityCache).
		To(func(_ *classroomCache.CacheClassroom, _ context.Context, key string, rooms []*model.Classroom) error {
			if rooms == nil {
				c.failureWrites++
				if c.failureErr != nil {
					return c.failureErr
				}
				c.failures[key] = true
				return nil
			}
			c.writes++
			if c.writeErr != nil {
				return c.writeErr
			}
			c.days[key] = rooms
			return nil
		}).Build()
	t.Cleanup(func() { set.UnPatch() })
	sleep := mockey.Mock(time.Sleep).Return().Build()
	t.Cleanup(func() { sleep.UnPatch() })
	svc := NewClassroomService(context.Background(), &base.ClientSet{CacheClient: &cache.Cache{Classroom: new(classroomCache.CacheClassroom)}})
	return s, c, svc
}

func TestRefreshEmptyRoomDateRetriesOnlyFailedBuildings(t *testing.T) {
	source, cache, runner := fixtures(t)
	runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-10-17", []string{"旗山校区", "铜盘校区"})
	require.Equal(t, map[string]int{"旗山校区": 1, "铜盘校区": 1}, source.lists)
	// 两个校区先分别批量查询，全部完成后才补抓各自失败的楼B。
	require.Len(t, source.requests, 4)
	for _, req := range source.requests[:2] {
		require.Equal(t, "2026-10-17", req.Date)
		require.Equal(t, []string{"楼A", "楼B", "楼C"}, req.Buildings)
	}
	for _, req := range source.requests[2:] {
		require.Equal(t, []string{"楼B"}, req.Buildings)
		require.Equal(t, "202601", req.Term)
	}
	require.Len(t, cache.days, 2)
	for _, rooms := range cache.days {
		require.Len(t, rooms, 3)
	}
}

func TestRefreshEmptyRoomDateEmptyCampus(t *testing.T) {
	source, cache, runner := fixtures(t)
	source.empty = true
	runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-10-17", []string{"旗山校区"})
	require.Len(t, source.requests, 1)
	rooms := cache.days["classroom:availability:2026-10-17:旗山校区"]
	require.NotNil(t, rooms, "成功但无教室必须写入空数组，而不是失败标记")
	require.Empty(t, rooms)
	require.Empty(t, cache.failures)
}

func TestRefreshEmptyRoomDateMarksPersistentFailure(t *testing.T) {
	source, cache, runner := fixtures(t)
	source.alwaysFail = true
	runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-10-17", []string{"旗山校区", "铜盘校区"})
	require.Len(t, source.requests, 4)
	require.Empty(t, cache.days, "失败楼未补齐时不能发布部分数据")
	require.Equal(t, map[string]bool{
		"classroom:availability:2026-10-17:旗山校区": true,
		"classroom:availability:2026-10-17:铜盘校区": true,
	}, cache.failures)
}

func TestRefreshEmptyRoomDateRetriesBuildingList(t *testing.T) {
	source, cache, runner := fixtures(t)
	source.failList = true
	runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-10-17", []string{"旗山校区", "铜盘校区"})
	require.Equal(t, map[string]int{"旗山校区": 2, "铜盘校区": 2}, source.lists)
	require.Len(t, cache.days, 2)
}

func TestRefreshEmptyRoomDateUpstreamFailure(t *testing.T) {
	for _, listFailure := range []bool{true, false} {
		t.Run(fmt.Sprint(listFailure), func(t *testing.T) {
			source, cache, runner := fixtures(t)
			source.listAlwaysFails = listFailure
			source.allFail = !listFailure
			runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-10-17", []string{"旗山校区"})
			require.Empty(t, cache.days)
			require.True(t, cache.failures["classroom:availability:2026-10-17:旗山校区"])
			if listFailure {
				require.Equal(t, 2, source.lists["旗山校区"])
				require.Empty(t, source.requests)
			} else {
				require.Len(t, source.requests, 2)
				for _, req := range source.requests {
					require.Equal(t, []string{"楼A", "楼B", "楼C"}, req.Buildings)
				}
			}
		})
	}
}

func TestRefreshEmptyRoomDateRetriesPageFailure(t *testing.T) {
	for _, persistent := range []bool{false, true} {
		t.Run(fmt.Sprint(persistent), func(t *testing.T) {
			source, cache, runner := fixtures(t)
			source.topLevelError = true
			source.alwaysFail = persistent
			runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-10-17", []string{"旗山校区"})
			require.Len(t, source.requests, 2)
			for _, req := range source.requests {
				require.Equal(t, []string{"楼A", "楼B", "楼C"}, req.Buildings)
			}
			key := "classroom:availability:2026-10-17:旗山校区"
			if persistent {
				require.Empty(t, cache.days)
				require.True(t, cache.failures[key])
			} else {
				require.Len(t, cache.days[key], 3)
				require.Empty(t, cache.failures)
			}
		})
	}
}

func TestRefreshEmptyRoomDateCacheFailureDoesNotRefetch(t *testing.T) {
	source, cache, runner := fixtures(t)
	source.failBuilding = false
	cache.writeErr = assert.AnError
	runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-10-17", []string{"旗山校区", "铜盘校区"})
	require.Len(t, source.requests, 2, "缓存失败后继续下一校区，不重新爬取")
	require.Equal(t, 2, cache.writes)
	require.Empty(t, cache.days)
	require.Zero(t, cache.failureWrites, "缓存失败不能标记成上游抓取失败")
}

func TestRefreshEmptyRoomDateFailureMarkerWriteDoesNotStopBatch(t *testing.T) {
	for _, listFailure := range []bool{true, false} {
		t.Run(fmt.Sprint(listFailure), func(t *testing.T) {
			source, cache, runner := fixtures(t)
			source.listAlwaysFails = listFailure
			source.alwaysFail = !listFailure
			cache.failureErr = assert.AnError
			runner.refreshEmptyRoomDate(jwch.NewStudent(), "2026-11-13", []string{"旗山校区", "铜盘校区"})
			require.Equal(t, 2, cache.failureWrites)
			require.Contains(t, source.lists, "旗山校区")
			require.Contains(t, source.lists, "铜盘校区")
		})
	}
}

func TestSaveAvailability(t *testing.T) {
	for _, tc := range []struct {
		name     string
		campus   string
		rooms    []*model.Classroom
		writeErr error
	}{
		{name: "ordinary campus", campus: "旗山校区", rooms: []*model.Classroom{{Build: "东2", Location: "旗山东2-101"}}},
		{name: "Xiamen campuses", campus: "厦门工艺美院", rooms: []*model.Classroom{
			{Location: "鼓浪屿多媒体1"},
			{Location: "集美1-311", Occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 4, Event: "讲座"}}},
		}},
		{name: "empty campus", campus: "厦门工艺美院", rooms: []*model.Classroom{}},
		{name: "refresh failure", campus: "厦门工艺美院"},
		{name: "both writes attempted on failure", campus: "厦门工艺美院", writeErr: assert.AnError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			writes := make(map[string][]*model.Classroom)
			set := mockey.Mock((*classroomCache.CacheClassroom).SetAvailabilityCache).
				To(func(_ *classroomCache.CacheClassroom, _ context.Context, key string, rooms []*model.Classroom) error {
					writes[key] = rooms
					return tc.writeErr
				}).Build()
			t.Cleanup(func() { set.UnPatch() })
			svc := NewClassroomService(context.Background(), &base.ClientSet{
				CacheClient: &cache.Cache{Classroom: new(classroomCache.CacheClassroom)},
			})
			svc.saveAvailability("2026-10-11", tc.campus, tc.rooms)
			if tc.campus == "旗山校区" {
				require.Equal(t, map[string][]*model.Classroom{"classroom:availability:2026-10-11:旗山校区": tc.rooms}, writes)
				return
			}
			require.Len(t, writes, 2)
			guLangYu, ok := writes["classroom:availability:2026-10-11:鼓浪屿校区"]
			require.True(t, ok)
			jiMei, ok := writes["classroom:availability:2026-10-11:集美校区"]
			require.True(t, ok)
			switch {
			case tc.rooms == nil:
				require.Nil(t, guLangYu)
				require.Nil(t, jiMei)
			case len(tc.rooms) == 0:
				require.NotNil(t, guLangYu)
				require.NotNil(t, jiMei)
				require.Empty(t, guLangYu)
				require.Empty(t, jiMei)
			default:
				require.Len(t, guLangYu, 1)
				require.Len(t, jiMei, 1)
				assert.Equal(t, "鼓浪屿校区", guLangYu[0].Build)
				assert.Equal(t, "集美校区", jiMei[0].Build)
				assert.Equal(t, "讲座", jiMei[0].Occupancies[0].Event)
			}
		})
	}
}

func TestRefreshEmptyRoom(t *testing.T) {
	for _, tc := range []struct {
		name     string
		loginErr error
	}{
		{name: "refreshes supplied dates in order"},
		{name: "preserves login failure", loginErr: assert.AnError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			previousUser := config.DefaultUser
			config.DefaultUser = nil
			t.Cleanup(func() { config.DefaultUser = previousUser })
			require.NoError(t, json.Unmarshal([]byte(`{"Account":"test-account","Password":"test-password"}`), &config.DefaultUser))
			svc := NewClassroomService(context.Background(), &base.ClientSet{CacheClient: &cache.Cache{}})
			dates := []string{"2026-11-07", "2026-11-09", "2026-11-15"}
			loginCalled := false
			var loggedInStudent *jwch.Student
			login := mockey.Mock(utils.RetryLogin).To(func(stu *jwch.Student) error {
				loginCalled = true
				assert.NotNil(t, stu)
				loggedInStudent = stu
				return tc.loginErr
			}).Build()
			t.Cleanup(func() { login.UnPatch() })
			queried := make(map[string][]string)
			var queriedDates []string
			refresh := mockey.Mock((*ClassroomService).refreshEmptyRoomDate).To(func(s *ClassroomService, stu *jwch.Student, date string, campuses []string) {
				assert.Same(t, loggedInStudent, stu)
				queriedDates = append(queriedDates, date)
				queried[date] = append([]string(nil), campuses...)
			}).Build()
			t.Cleanup(func() { refresh.UnPatch() })
			err := svc.RefreshEmptyRoom(dates)
			switch {
			case tc.loginErr != nil:
				require.ErrorIs(t, err, tc.loginErr)
				require.True(t, loginCalled)
				require.Empty(t, queried)
			default:
				require.NoError(t, err)
				require.True(t, loginCalled)
				require.Len(t, queried, len(dates))
				for i, date := range dates {
					assert.Equal(t, date, queriedDates[i], "按传入顺序刷新")
					assert.Equal(t, constants.CampusArray, queried[date])
				}
			}
		})
	}
}
