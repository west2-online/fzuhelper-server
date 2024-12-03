package launch_screen

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

func TestDBLaunchScreen_AddImageListShowTime(t *testing.T) {
	type testCase struct {
		name                string
		mockError           error
		inputPictureList    *[]model.Picture
		expectedPictureList *[]model.Picture
		expectingError      bool
	}

	testCases := []testCase{
		{
			name:      "AddImageListShowTime_Success",
			mockError: nil,
			inputPictureList: &[]model.Picture{
				{ID: 1, Url: "https://example.com/image1.jpg", ShowTimes: 2},
				{ID: 2, Url: "https://example.com/image2.jpg", ShowTimes: 5},
			},
			expectedPictureList: &[]model.Picture{
				{ID: 1, Url: "https://example.com/image1.jpg", ShowTimes: 3},
				{ID: 2, Url: "https://example.com/image2.jpg", ShowTimes: 6},
			},
			expectingError: false,
		},
		{
			name:      "AddImageListShowTime_DBError",
			mockError: fmt.Errorf("db error"),
			inputPictureList: &[]model.Picture{
				{ID: 1, Url: "https://example.com/image1.jpg", ShowTimes: 2},
			},
			expectedPictureList: nil,
			expectingError:      true,
		},
		{
			name:                "AddImageListShowTime_EmptyList",
			mockError:           nil,
			inputPictureList:    &[]model.Picture{},
			expectedPictureList: &[]model.Picture{},
			expectingError:      false,
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

			mockey.Mock((*gorm.DB).Save).To(func(value interface{}) *gorm.DB {
				if tc.mockError != nil {
					return &gorm.DB{Error: tc.mockError}
				}
				return &gorm.DB{Error: nil}
			}).Build()

			err := mockDBLaunchScreen.AddImageListShowTime(context.Background(), tc.inputPictureList)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "dal.AddImageListShowTime error")

			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPictureList, tc.inputPictureList)
			}
		})
	}
}
