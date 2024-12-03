package launch_screen

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"gorm.io/gorm"
)

func TestDBLaunchScreen_DeleteImage(t *testing.T) {
	type testCase struct {
		name            string
		inputID         int64
		mockErrorTake   error
		mockErrorDelete error
		expectedResult  *model.Picture
		expectingError  bool
	}
	picture := &model.Picture{
		ID:         2024,
		Url:        "newUrl",
		Href:       "href",
		Text:       "text",
		PicType:    3,
		ShowTimes:  0,
		PointTimes: 0,
		Duration:   0,
		StartAt:    time.Now().Add(-24 * time.Hour),
		EndAt:      time.Now().Add(24 * time.Hour),
		StartTime:  0,
		EndTime:    24,
		SType:      3,
		Frequency:  4,
		Regex:      "{\"device\": \"android,ios\", \"student_id\": \"102301517,102301544\"}",
	}
	testCases := []testCase{
		{
			name:            "DeleteImage_Success",
			inputID:         1001,
			mockErrorTake:   nil,
			mockErrorDelete: nil,
			expectedResult:  picture,
			expectingError:  false,
		},
		{
			name:            "DeleteImage_TakeError",
			inputID:         1001,
			mockErrorTake:   errors.New("record not found"),
			mockErrorDelete: nil,
			expectedResult:  nil,
			expectingError:  true,
		},
		{
			name:            "DeleteImage_DeleteError",
			inputID:         1001,
			mockErrorTake:   nil,
			mockErrorDelete: errors.New("delete failed"),
			expectedResult:  picture,
			expectingError:  true,
		},
	}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBLaunchScreen := NewDBLaunchScreen(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Take).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockErrorTake != nil {
					return &gorm.DB{Error: tc.mockErrorTake}
				}
				*dest.(*model.Picture) = *tc.expectedResult
				return &gorm.DB{Error: nil}
			}).Build()

			mockey.Mock((*gorm.DB).Delete).To(func(value interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockErrorDelete != nil {
					return &gorm.DB{Error: tc.mockErrorDelete}
				}
				return &gorm.DB{Error: nil}
			}).Build()

			result, err := mockDBLaunchScreen.DeleteImage(context.Background(), tc.inputID)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "dal.DeleteImage error")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
