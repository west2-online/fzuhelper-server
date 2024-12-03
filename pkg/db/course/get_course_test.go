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

func TestDBCourse_GetUserTermCourseByStuIdAndTerm(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		mockRows       *gorm.DB
		stuId          string
		term           string
		expectedResult *model.UserCourse
		expectingError bool
	}
	testCases := []testCase{
		{
			name:      "GetUserTermCourseByStuIdAndTerm_Success",
			mockError: nil,
			mockRows:  &gorm.DB{Error: nil},
			stuId:     "222200311",
			term:      "202401",
			expectedResult: &model.UserCourse{
				Id:                1001,
				StuId:             "222200311",
				Term:              "202401",
				TermCourses:       `[{"courseId":"C123","courseName":"Math"}]`,
				TermCoursesSha256: "abc123def456"},
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseByStuIdAndTerm_RecordNotFound",
			mockError:      nil,
			mockRows:       &gorm.DB{Error: gorm.ErrRecordNotFound},
			stuId:          "222200311",
			term:           "202401",
			expectedResult: nil,
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseByStuIdAndTerm_DBError",
			mockError:      fmt.Errorf("db error"),
			mockRows:       &gorm.DB{Error: fmt.Errorf("db error")},
			stuId:          "222200311",
			term:           "202401",
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

			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				} else if tc.mockRows != nil && tc.mockRows.Error == gorm.ErrRecordNotFound {
					mockGormDB.Error = gorm.ErrRecordNotFound
					return mockGormDB
				} else {
					*dest.(*model.UserCourse) = *tc.expectedResult
					return mockGormDB
				}

			}).Build()

			result, err := mockDBCourse.GetUserTermCourseByStuIdAndTerm(context.Background(), tc.stuId, tc.term)

			if tc.mockRows != nil && tc.mockRows.Error == gorm.ErrRecordNotFound {
				assert.NoError(t, err)
				assert.Nil(t, result)
			} else if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestDBCourse_GetUserTermCourseSha256ByStuIdAndTerm(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		mockRows       *gorm.DB
		stuId          string
		term           string
		expectedResult *model.UserCourse
		expectingError bool
	}

	testCases := []testCase{
		{
			name:      "GetUserTermCourseSha256ByStuIdAndTerm_Success",
			mockError: nil,
			mockRows:  &gorm.DB{Error: nil}, // No error
			stuId:     "222200311",
			term:      "202401",
			expectedResult: &model.UserCourse{
				Id:                1001,
				StuId:             "222200311",
				Term:              "202401",
				TermCourses:       `[{"courseId":"C123","courseName":"Math"}]`,
				TermCoursesSha256: "abc123def456"},
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseSha256ByStuIdAndTerm_RecordNotFound",
			mockError:      nil,
			mockRows:       &gorm.DB{Error: gorm.ErrRecordNotFound}, // Simulate not found
			stuId:          "222200311",
			term:           "202401",
			expectedResult: nil,
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseSha256ByStuIdAndTerm_DBError",
			mockError:      fmt.Errorf("db error"),
			mockRows:       &gorm.DB{Error: fmt.Errorf("db error")}, // Simulate db error
			stuId:          "222200311",
			term:           "202401",
			expectedResult: nil,
			expectingError: true,
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBCourse := NewDBCourse(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Select).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					return &gorm.DB{Error: tc.mockError}
				} else if tc.mockRows != nil && tc.mockRows.Error == gorm.ErrRecordNotFound {
					return &gorm.DB{Error: gorm.ErrRecordNotFound}
				} else {
					*dest.(*model.UserCourse) = *tc.expectedResult
					return &gorm.DB{Error: nil}
				}
			}).Build()

			result, err := mockDBCourse.GetUserTermCourseSha256ByStuIdAndTerm(context.Background(), tc.stuId, tc.term)

			if tc.mockRows != nil && tc.mockRows.Error == gorm.ErrRecordNotFound {
				assert.NoError(t, err)
				assert.Nil(t, result)
			} else if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})

	}
}
