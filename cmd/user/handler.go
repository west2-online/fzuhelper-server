package main

import (
	"context"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetLoginData implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetLoginData(ctx context.Context, request *user.GetLoginDataRequest) (resp *user.GetLoginDataResponse, err error) {
	// TODO: Your code here...
	return
}
