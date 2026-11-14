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
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
)

func TestGetEmptyRoom(t *testing.T) {
	for _, tc := range []struct {
		name            string
		includeOccupied bool
		cacheErr        error
		wantLocations   []string
	}{
		{name: "legacy request filters occupied rooms", wantLocations: []string{"旗山东2-102"}},
		{name: "opt in includes occupied rooms", includeOccupied: true, wantLocations: []string{"旗山东2-101", "旗山东2-102"}},
		{name: "cache miss remains an error", cacheErr: redis.Nil},
		{name: "redis failure is preserved", cacheErr: assert.AnError},
	} {
		mockey.PatchConvey(tc.name, t, func() {
			rooms := []*model.Classroom{
				{Location: "旗山东2-101",
					FreeRanges:  []*model.ClassroomSectionRange{{Start: 1, End: 3}, {Start: 5, End: 11}},
					Occupancies: []*model.ClassroomOccupancy{{Start: 4, End: 4, Event: "补课"}}},
				{Location: "旗山东2-102", FreeRanges: []*model.ClassroomSectionRange{{Start: 1, End: 11}}},
			}
			mockey.Mock((*classroomCache.CacheClassroom).GetAvailability).
				To(func(_ *classroomCache.CacheClassroom, _ context.Context, date, campus string) (*classroomCache.DayAvailability, error) {
					assert.Equal(t, "2026-10-11", date)
					assert.Equal(t, "旗山校区", campus)
					return &classroomCache.DayAvailability{Date: date, Term: "202601", Rooms: rooms}, tc.cacheErr
				}).Build()
			svc := NewClassroomService(context.Background(), &base.ClientSet{
				CacheClient: &cache.Cache{Classroom: new(classroomCache.CacheClassroom)},
			})
			req := &classroom.EmptyRoomRequest{Date: "2026-10-11", Campus: "旗山校区", StartTime: "3", EndTime: "5"}
			if tc.includeOccupied {
				req.IncludeOccupied = &tc.includeOccupied
			}
			result, err := svc.GetEmptyRoom(req)
			if tc.cacheErr != nil {
				require.Error(t, err)
				if tc.cacheErr != redis.Nil {
					assert.True(t, errors.Is(err, tc.cacheErr))
				}
				return
			}
			require.NoError(t, err)
			locations := make([]string, 0, len(result))
			for _, room := range result {
				locations = append(locations, room.Location)
			}
			assert.Equal(t, tc.wantLocations, locations)
			if tc.includeOccupied {
				assert.Equal(t, []*model.ClassroomSectionRange{{Start: 1, End: 3}, {Start: 5, End: 11}}, result[0].FreeRanges)
				assert.Equal(t, "补课", result[0].Occupancies[0].Event)
				assert.Equal(t, []*model.ClassroomSectionRange{{Start: 1, End: 11}}, result[len(result)-1].FreeRanges)
			} else {
				// 字段裁剪由 api pack 层按 includeOccupied 完成，service 层原样返回快照。
				require.Len(t, result, 1)
				assert.Same(t, rooms[1], result[0])
				assert.Equal(t, []*model.ClassroomSectionRange{{Start: 1, End: 11}}, result[0].FreeRanges)
			}
			assert.Equal(t, []*model.ClassroomSectionRange{{Start: 1, End: 3}, {Start: 5, End: 11}}, rooms[0].FreeRanges)
			assert.NotEmpty(t, rooms[0].Occupancies)
		})
	}
}

func TestFilterRoomsInclusiveRange(t *testing.T) {
	for _, tc := range []struct {
		name        string
		start, end  int32
		occupancies []*model.ClassroomOccupancy
		wantFree    bool
	}{
		{name: "gap between occupancies", start: 2, end: 4,
			occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 1}, {Start: 5, End: 11}}, wantFree: true},
		{name: "overlap at start", start: 1, end: 2,
			occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 1}, {Start: 5, End: 11}}},
		{name: "overlap at end", start: 4, end: 5,
			occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 1}, {Start: 5, End: 11}}},
		{name: "free all day", start: 1, end: 11, wantFree: true},
		{name: "overlapping occupancies", start: 1, end: 11,
			occupancies: []*model.ClassroomOccupancy{{Start: 8, End: 9}, {Start: 3, End: 4}, {Start: 4, End: 5}}},
		{name: "occupied single section", start: 4, end: 4,
			occupancies: []*model.ClassroomOccupancy{{Start: 8, End: 9}, {Start: 3, End: 4}, {Start: 4, End: 5}}},
		{name: "outside occupancies", start: 6, end: 7,
			occupancies: []*model.ClassroomOccupancy{{Start: 8, End: 9}, {Start: 3, End: 4}, {Start: 4, End: 5}}, wantFree: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			room := &model.Classroom{Occupancies: tc.occupancies}
			result := filterRooms([]*model.Classroom{room}, tc.start, tc.end, false)
			if tc.wantFree {
				require.Len(t, result, 1)
				assert.Same(t, room, result[0])
			} else {
				assert.Empty(t, result)
			}
		})
	}
	room := &model.Classroom{Occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 11}}}
	result := filterRooms([]*model.Classroom{room}, 1, 11, true)
	require.Len(t, result, 1)
	assert.Same(t, room, result[0])
	assert.NotEmpty(t, room.Occupancies)
}

func TestParseSections(t *testing.T) {
	for _, values := range [][2]string{{"0", "1"}, {"1", "12"}, {"5", "3"}, {"x", "3"}, {"1", ""}} {
		_, _, err := parseSections(values[0], values[1])
		assert.Error(t, err)
	}
	start, end, err := parseSections("01", "11")
	require.NoError(t, err)
	assert.EqualValues(t, 1, start)
	assert.EqualValues(t, 11, end)
}
