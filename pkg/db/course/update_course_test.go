package course

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"gorm.io/gorm"
)

func TestDBCourse_UpdateUserTermCourse(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		input          *model.UserCourse
		expectedResult *model.UserCourse
		expectingError bool
	}

	testCases := []testCase{
		{
			name:      "UpdateUserTermCourse_Success",
			mockError: nil,
			input: &model.UserCourse{
				Id:                1001,
				StuId:             "222200311",
				Term:              "202401",
				TermCourses:       `[{"courseId":"C123","courseName":"Math"}]`,
				TermCoursesSha256: "abc123def456",
			},
			expectedResult: &model.UserCourse{
				Id:                1001,
				StuId:             "222200311",
				Term:              "202401",
				TermCourses:       `[{"courseId":"C123","courseName":"Math"}]`,
				TermCoursesSha256: "abc123def456",
			},
			expectingError: false,
		},
		{
			name:           "UpdateUserTermCourse_DBError",
			mockError:      fmt.Errorf("db error"),
			input:          &model.UserCourse{Id: 1002, StuId: "222200311"},
			expectedResult: nil,
			expectingError: true,
		},
	}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBCourse := NewDBCourse(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Model).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Updates).To(func(values interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				return &gorm.DB{Error: nil}
			}).Build()

			result, err := mockDBCourse.UpdateUserTermCourse(context.Background(), tc.input)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "dal.UpdateUserTermCourse error")

			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
