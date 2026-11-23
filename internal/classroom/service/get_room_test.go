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
	"time"

	"github.com/bytedance/mockey"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	classroomCache "github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func TestGetEmptyRoomVersions(t *testing.T) {
	for _, version := range []string{"v1", "v2"} {
		for _, tc := range []struct {
			name     string
			cacheErr error
			empty    bool
			failed   bool
		}{
			{name: "classrooms"},
			{name: "empty campus", empty: true},
			{name: "refresh failure", failed: true},
			{name: "cache miss", cacheErr: redis.Nil},
			{name: "redis failure", cacheErr: assert.AnError},
		} {
			mockey.PatchConvey(version+" "+tc.name, t, func() {
				mockey.Mock(time.Now).Return(time.Date(2026, 10, 11, 8, 0, 0, 0, constants.ChinaTZ)).Build()
				rooms := []*model.Classroom{
					{
						Build: "东2", Location: "旗山东2-101", Capacity: "100", Type: "多媒体",
						FreeRanges:  []*model.ClassroomSectionRange{{Start: 1, End: 3}, {Start: 5, End: 11}},
						Occupancies: []*model.ClassroomOccupancy{{Start: 4, End: 4, Event: "补课"}},
					},
					{
						Build: "东2", Location: "旗山东2-102", Capacity: "80", Type: "普通教室",
						FreeRanges: []*model.ClassroomSectionRange{{Start: 1, End: 11}},
					},
				}
				cachedRooms := rooms
				if tc.empty {
					cachedRooms = make([]*model.Classroom, 0)
				}
				if tc.failed {
					cachedRooms = nil
				}
				mockey.Mock((*classroomCache.CacheClassroom).GetAvailabilityCache).
					To(func(_ *classroomCache.CacheClassroom, _ context.Context, key string) ([]*model.Classroom, error) {
						assert.Equal(t, "classroom:availability:2026-10-11:旗山校区", key)
						return cachedRooms, tc.cacheErr
					}).Build()
				svc := NewClassroomService(context.Background(), &base.ClientSet{CacheClient: &cache.Cache{Classroom: new(classroomCache.CacheClassroom)}})
				var result []*model.Classroom
				var err error
				if version == "v1" {
					result, err = svc.GetEmptyRoom(&classroom.EmptyRoomRequest{Date: "2026-10-11", Campus: "旗山校区", StartTime: "3", EndTime: "5"})
				} else {
					result, err = svc.GetEmptyRoomV2(&classroom.EmptyRoomV2Request{Date: "2026-10-11", Campus: "旗山校区"})
				}
				if tc.failed {
					require.ErrorContains(t, err, "教室数据拉取失败")
					return
				}
				if tc.cacheErr != nil {
					require.Error(t, err)
					if !errors.Is(tc.cacheErr, redis.Nil) {
						require.ErrorIs(t, err, tc.cacheErr)
					}
					return
				}
				require.NoError(t, err)
				if tc.empty {
					require.NotNil(t, result)
					require.Empty(t, result)
					return
				}
				if version == "v1" {
					require.Equal(t, []*model.Classroom{{Build: "东2", Location: "旗山东2-102", Capacity: "80", Type: "普通教室"}}, result)
				} else {
					require.Equal(t, rooms, result)
					require.Len(t, result, 2)
					require.Equal(t, "补课", result[0].Occupancies[0].Event)
				}
				require.NotEmpty(t, rooms[0].Occupancies, "V1 过滤不能修改共享日缓存")
			})
		}
	}
}

func TestGetEmptyRoomFiltersFullRequestedRange(t *testing.T) {
	for _, tc := range []struct {
		start, end string
		want       int
	}{
		{start: "1", end: "3", want: 1},
		{start: "4", end: "4", want: 0},
		{start: "3", end: "5", want: 0},
		{start: "5", end: "11", want: 1},
		{start: "1", end: "11", want: 0},
	} {
		mockey.PatchConvey(tc.start+"-"+tc.end, t, func() {
			mockey.Mock(time.Now).Return(time.Date(2026, 10, 11, 8, 0, 0, 0, constants.ChinaTZ)).Build()
			mockey.Mock((*classroomCache.CacheClassroom).GetAvailabilityCache).Return([]*model.Classroom{{
				Location: "旗山东2-101", FreeRanges: []*model.ClassroomSectionRange{{Start: 1, End: 3}, {Start: 5, End: 11}},
			}}, nil).Build()
			svc := NewClassroomService(context.Background(), &base.ClientSet{CacheClient: &cache.Cache{Classroom: new(classroomCache.CacheClassroom)}})
			rooms, err := svc.GetEmptyRoom(&classroom.EmptyRoomRequest{Date: "2026-10-11", Campus: "旗山校区", StartTime: tc.start, EndTime: tc.end})
			require.NoError(t, err)
			require.Len(t, rooms, tc.want)
		})
	}
}

func TestClassroomDateValidation(t *testing.T) {
	for _, date := range []string{"invalid", "2026-10-10", "2026-10-19"} {
		mockey.PatchConvey(date, t, func() {
			mockey.Mock(time.Now).Return(time.Date(2026, 10, 11, 8, 0, 0, 0, constants.ChinaTZ)).Build()
			svc := NewClassroomService(context.Background(), &base.ClientSet{})
			_, err := svc.GetEmptyRoom(&classroom.EmptyRoomRequest{Date: date, Campus: "旗山校区", StartTime: "1", EndTime: "11"})
			require.Error(t, err)
			_, err = svc.GetEmptyRoomV2(&classroom.EmptyRoomV2Request{Date: date, Campus: "旗山校区"})
			require.Error(t, err)
		})
	}
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
