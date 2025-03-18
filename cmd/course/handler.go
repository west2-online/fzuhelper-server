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
