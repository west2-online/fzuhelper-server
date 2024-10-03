//MVC--Model

package db

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/user/pack/pwd"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int64
	Account   string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

func Register(ctx context.Context, userModel *User) (*User, error) {
	userResp := new(User)
	//WithContext(ctx)是将一个context.Context对象和数据库连接绑定，以实现在数据库操作中使用context.Context上下文传递。
	if err := DB.WithContext(ctx).Where("account = ? OR name = ?", userModel.Account, userModel.Name).First(&userResp).Error; err == nil {
		return nil, errno.UserExistedError
	}

	if err := DB.WithContext(ctx).Create(userModel).Error; err != nil {
		return nil, err
	}
	return userModel, nil
}

func Login(ctx context.Context, userModel *User) (*User, error) {
	userResp := new(User)
	if err := DB.WithContext(ctx).Where("account = ?", userModel.Account).
		First(&userResp).Error; err != nil {
		return nil, errno.UserNonExistError
	}

	if !pwd.CheckPassword(userResp.Password, userModel.Password) {
		return nil, errno.AuthFailedError
	}

	return userResp, nil
}
