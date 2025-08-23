package main

import (
	"context"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetLoginData implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetLoginData(ctx context.Context, req *user.GetLoginDataRequest) (resp *user.GetLoginDataResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, request *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// GetGetLoginDataForYJSY implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetGetLoginDataForYJSY(ctx context.Context, request *user.GetLoginDataForYJSYRequest) (resp *user.GetLoginDataForYJSYResponse, err error) {
	// TODO: Your code here...
	return
}
