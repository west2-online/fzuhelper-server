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

package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
)

func TestBuildUsageRooms(t *testing.T) {
	rooms := BuildUsageRooms([]*jwch.ClassroomUsageRoom{
		{Location: "旗山东2-100", Capacity: "318", Type: "智慧教室报告厅型",
			Occupancies: []*jwch.ClassroomUsageOccupancy{
				{Start: 1, End: 4, StartWeek: 6, EndWeek: 6, Event: "讲座", Occupant: "先造", RawText: "讲座()(06-06)"},
				{Start: 3, End: 5, Event: "补课"},
				{Start: 9, End: 11, Event: "课设"},
			}},
		{Location: "旗山东2-101", Capacity: "123", Type: "智慧教室普通型"},
	}, "旗山校区")
	require.Len(t, rooms, 2)
	assert.Equal(t, "东2", rooms[0].Build)
	assert.Equal(t, "318", rooms[0].Capacity)
	assert.Equal(t, "智慧教室报告厅型", rooms[0].Type)
	assert.Len(t, rooms[0].Occupancies, 3)
	assert.Equal(t, []*model.ClassroomSectionRange{{Start: 6, End: 8}}, rooms[0].FreeRanges)
	assert.EqualValues(t, 6, rooms[0].Occupancies[0].StartWeek)
	assert.EqualValues(t, 6, rooms[0].Occupancies[0].EndWeek)
	assert.EqualValues(t, 0, rooms[0].Occupancies[1].StartWeek)
	assert.EqualValues(t, 0, rooms[0].Occupancies[1].EndWeek)
	assert.Equal(t, "讲座", rooms[0].Occupancies[0].Event)
	assert.Equal(t, "先造", rooms[0].Occupancies[0].Occupant)
	assert.Equal(t, "讲座()(06-06)", rooms[0].Occupancies[0].RawText)
	assert.Empty(t, rooms[1].Occupancies)
	assert.Equal(t, []*model.ClassroomSectionRange{{Start: 1, End: 11}}, rooms[1].FreeRanges)
	assert.Empty(t, BuildFreeRanges([]*model.ClassroomOccupancy{{Start: 1, End: 11}}))
}

func TestLegacyClassroomMetadata(t *testing.T) {
	for _, tc := range []struct{ text, campus, build string }{
		{"旗山东2-100 318(159) 智慧教室报告厅型", "旗山校区", "东2"},
		{"旗山公语-101 123(60) 多媒体", "旗山校区", "东1"},
		{"旗山中-101 123(60) 多媒体", "旗山校区", "中楼"},
		{"晋江A102 150(75) 多媒体", "晋江校区", "晋江校区"},
	} {
		room := BuildClassroom(tc.text, tc.campus)
		require.NotNil(t, room)
		assert.Equal(t, tc.build, room.Build)
	}
	assert.Nil(t, BuildClassroom("invalid", "旗山校区"))
}
