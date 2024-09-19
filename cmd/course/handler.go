package main

import (
	"context"
	course "github.com/west2-online/fzuhelper-server/kitex_gen/course"
)

// CourseServiceImpl implements the last service interface defined in the IDL.
type CourseServiceImpl struct{}

// GetCourseList implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetCourseList(ctx context.Context, req *course.CourseListRequest) (resp *course.CourseListResponse, err error) {
	// TODO: Your code here...
	return
}
