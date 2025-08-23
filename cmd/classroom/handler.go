package main

import (
	"context"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
)

// ClassroomServiceImpl implements the last service interface defined in the IDL.
type ClassroomServiceImpl struct{}

// GetEmptyRoom implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetEmptyRoom(ctx context.Context, req *classroom.EmptyRoomRequest) (resp *classroom.EmptyRoomResponse, err error) {
	// TODO: Your code here...
	return
}

// GetExamRoomInfo implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetExamRoomInfo(ctx context.Context, req *classroom.ExamRoomInfoRequest) (resp *classroom.ExamRoomInfoResponse, err error) {
	// TODO: Your code here...
	return
}
