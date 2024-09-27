package service

import (
	"context"
	"fmt"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"testing"
)

func TestGetLoginData(t *testing.T) {
	ctx := context.Background()
	s := NewUserService(ctx, "", nil)
	id, cookies, err := s.GetLoginData(&user.GetLoginDataRequest{
		Id:       "082100170",
		Password: "Zhuyinfan815",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)
	fmt.Println(cookies)
}
