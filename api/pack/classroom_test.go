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

	"github.com/cloudwego/frugal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestClassroomAvailabilityAcrossRPC(t *testing.T) {
	room := &model.Classroom{
		Build: "东2", Location: "旗山东2-100", Capacity: "318", Type: "智慧教室报告厅型",
		FreeRanges:  []*model.ClassroomSectionRange{{Start: 5, End: 11}},
		Occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 4, StartWeek: 6, EndWeek: 6, Event: "讲座", Occupant: "先造", RawText: "讲座()(06-06)"}},
	}
	response := &classroom.EmptyRoomResponse{Base: &model.BaseResp{Code: 10000}, Rooms: []*model.Classroom{room}}
	buffer := make([]byte, frugal.EncodedSize(response))
	_, err := frugal.EncodeObject(buffer, nil, response)
	require.NoError(t, err)
	var decoded classroom.EmptyRoomResponse
	_, err = frugal.DecodeObject(buffer, &decoded)
	require.NoError(t, err)
	require.Len(t, decoded.Rooms, 1)
	assert.Equal(t, room, decoded.Rooms[0])
	apiRoom := BuildClassroom(decoded.Rooms[0], true)
	assert.Equal(t, "318", apiRoom.Capacity)
	require.Len(t, apiRoom.FreeRanges, 1)
	assert.EqualValues(t, 5, apiRoom.FreeRanges[0].Start)
	assert.EqualValues(t, 11, apiRoom.FreeRanges[0].End)
	assert.Equal(t, "先造", apiRoom.Occupancies[0].Occupant)
	assert.EqualValues(t, 6, apiRoom.Occupancies[0].StartWeek)
	assert.EqualValues(t, 6, apiRoom.Occupancies[0].EndWeek)
	include := true
	request := &classroom.EmptyRoomRequest{Date: "2026-10-11", Campus: "旗山校区", StartTime: "1", EndTime: "4", IncludeOccupied: &include}
	buffer = make([]byte, frugal.EncodedSize(request))
	_, err = frugal.EncodeObject(buffer, nil, request)
	require.NoError(t, err)
	var decodedRequest classroom.EmptyRoomRequest
	_, err = frugal.DecodeObject(buffer, &decodedRequest)
	require.NoError(t, err)
	assert.True(t, decodedRequest.GetIncludeOccupied())
}
