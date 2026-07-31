/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package course

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBCourse_CreateExamOffering(t *testing.T) {
	testCases := []struct {
		name        string
		createError error
		expectError bool
		expectNil   bool
	}{
		{name: "success"},
		{name: "database error", createError: errors.New("insert failed"), expectError: true},
		{name: "duplicated", createError: gorm.ErrDuplicatedKey, expectNil: true},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			expected := &model.ExamOffering{ExamHash: "exam-hash", Tag: "exam-tag"}
			mockDB := new(gorm.DB)
			mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
			mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
			if tc.createError == nil {
				mockey.Mock((*gorm.DB).Create).Return(mockDB).Build()
			} else {
				mockey.Mock((*gorm.DB).Create).Return(&gorm.DB{Error: tc.createError}).Build()
			}

			result, err := NewDBCourse(mockDB, new(utils.Snowflake)).
				CreateExamOffering(context.Background(), expected)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "dal.CreateExamOffering error")
				return
			}
			assert.NoError(t, err)
			if tc.expectNil {
				assert.Nil(t, result)
				return
			}
			assert.Same(t, expected, result)
		})
	}
}
