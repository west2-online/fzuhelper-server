package service

import (
	"context"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"testing"
)

func TestAddPointTime(t *testing.T) {
	config.InitForTest()
	dal.Init()
	type testCase struct {
		name       string
		mockReturn interface{}
	}
	testCases := []testCase{
		{
			name:       "AddPointTime",
			mockReturn: nil,
		},
	}
	req := &launch_screen.AddImagePointTimeRequest{
		PictureId: 634858211152756736, //请确保该id对应picture存在
	}
	defer mockey.UnPatchAll() //撤销所有mock操作，不会影响其他测试

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background())

			mockey.Mock(db.AddPointTime).Return(tc.mockReturn).Build()

			err := launchScreenService.AddPointTime(req.PictureId)

			assert.NoError(t, err)
		})

	}
}
