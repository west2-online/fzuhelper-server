package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/jwch"
)

func TestGetExamRoomInfo(t *testing.T) {
	type testCase struct {
		name           string
		mockReturn     interface{}
		expectedResult interface{}
		expectingError bool
	}

	tests := []testCase{
		{
			name: "GetExamRoomInfo",
			mockReturn: []*jwch.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectedResult: []*jwch.ExamRoomInfo{
				{Location: "旗山东1"},
			},
			expectingError: true,
		},
	}

	req := &classroom.ExamRoomInfoRequest{
		Term:      "202401",
		LoginData: new(model.LoginData),
	}

	defer mockey.UnPatchAll()

	// 运行所有测试用例
	for _, tc := range tests {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockey.Mock((*jwch.Student).WithLoginData).Return(jwch.NewStudent()).Build()
			mockey.Mock((*jwch.Student).GetExamRoom).Return(tc.mockReturn, nil).Build()

			classroomService := NewClassroomService(context.Background(), mockClientSet)
			result, err := classroomService.GetExamRoomInfo(req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
