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

func TestDBLaunchScreen_AddPointTime(t *testing.T) {
	type testCase struct {
		name            string
		mockErrorFirst  error
		mockErrorSave   error
		inputId         int64
		initialPicture  *model.Picture
		expectedPicture *model.Picture
		expectingError  bool
	}

	testCases := []testCase{
		{
			name:            "AddPointTime_Success",
			mockErrorFirst:  nil,
			mockErrorSave:   nil,
			inputId:         1,
			initialPicture:  &model.Picture{ID: 1, Url: "https://example.com/image.jpg", PointTimes: 2},
			expectedPicture: &model.Picture{ID: 1, Url: "https://example.com/image.jpg", PointTimes: 3},
			expectingError:  false,
		},
		{
			name:            "AddPointTime_RecordNotFound",
			mockErrorFirst:  gorm.ErrRecordNotFound,
			mockErrorSave:   nil,
			inputId:         2,
			initialPicture:  nil,
			expectedPicture: nil,
			expectingError:  true,
		},
		{
			name:            "AddPointTime_DBErrorOnSave",
			mockErrorFirst:  nil,
			mockErrorSave:   fmt.Errorf("db save error"),
			inputId:         3,
			initialPicture:  &model.Picture{ID: 3, Url: "https://example.com/image2.jpg", PointTimes: 5},
			expectedPicture: nil,
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

			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockErrorFirst != nil {
					mockGormDB.Error = tc.mockErrorFirst
					return mockGormDB
				}
				if tc.initialPicture != nil {
					*dest.(*model.Picture) = *tc.initialPicture
				}
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Save).To(func(value interface{}) *gorm.DB {
				if tc.mockErrorSave != nil {
					mockGormDB.Error = tc.mockErrorSave
					return mockGormDB
				}
				if picture, ok := value.(*model.Picture); ok {
					// 直接 tc.initialPicture.PointTimes++ 是无效的,原因未知
					tc.initialPicture.PointTimes = picture.PointTimes
				}
				return mockGormDB
			}).Build()

			err := mockDBLaunchScreen.AddPointTime(context.Background(), tc.inputId)

			if tc.expectingError {
				assert.Error(t, err)
				if tc.mockErrorFirst != nil || tc.mockErrorSave != nil {
					assert.Contains(t, err.Error(), "dal.AddPointTime error")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPicture.PointTimes, tc.initialPicture.PointTimes)
			}
		})
	}
}
