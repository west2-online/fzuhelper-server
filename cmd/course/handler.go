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

// GetTermList implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetTermList(ctx context.Context, req *course.TermListRequest) (resp *course.TermListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCalendar implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetCalendar(ctx context.Context, req *course.GetCalendarRequest) (resp *course.GetCalendarResponse, err error) {
	// TODO: Your code here...
	return
}

// GetLocateDate implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetLocateDate(ctx context.Context, req *course.GetLocateDateRequest) (resp *course.GetLocateDateResponse, err error) {
	// TODO: Your code here...
	return
}
