package main

import (
	"context"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResp, err error) {
	// TODO: Your code here...
	return
}

// ValidateCode implements the UserServiceImpl interface.
func (s *UserServiceImpl) ValidateCode(ctx context.Context, req *user.ValidateCodeRequest) (resp *user.ValidateCodeResp, err error) {
	// TODO: Your code here...
	return
}

// ChangePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangePassword(ctx context.Context, req *user.ChangePasswordRequest) (resp *user.ChangePasswordResp, err error) {
	// TODO: Your code here...
	return
}

// GetSchoolCalendar implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetSchoolCalendar(ctx context.Context, req *user.SchoolCalendarRequest) (resp *user.SchoolCalendarResponse, err error) {
	// TODO: Your code here...
	return
}
