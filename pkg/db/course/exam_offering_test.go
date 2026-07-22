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
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBCourse_GetExamOfferingByHash(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer mockey.UnPatchAll()

		expected := &model.ExamOffering{
			ID:       1,
			ExamHash: "exam-hash",
			Tag:      "exam-tag",
		}
		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Where).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).First).To(func(dest interface{}, _ ...interface{}) *gorm.DB {
			offering, ok := dest.(*model.ExamOffering)
			if !ok {
				return &gorm.DB{Error: errors.New("unexpected destination type")}
			}
			*offering = *expected
			return mockDB
		}).Build()

		result, err := NewDBCourse(mockDB, new(utils.Snowflake)).
			GetExamOfferingByHash(context.Background(), expected.ExamHash)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		defer mockey.UnPatchAll()

		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Where).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).First).Return(&gorm.DB{Error: gorm.ErrRecordNotFound}).Build()

		result, err := NewDBCourse(mockDB, new(utils.Snowflake)).
			GetExamOfferingByHash(context.Background(), "missing")

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("database error", func(t *testing.T) {
		defer mockey.UnPatchAll()

		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Where).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).First).Return(&gorm.DB{Error: errors.New("query failed")}).Build()

		result, err := NewDBCourse(mockDB, new(utils.Snowflake)).
			GetExamOfferingByHash(context.Background(), "broken")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "dal.GetExamOfferingByHash error")
	})
}

func TestDBCourse_CreateExamOffering(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer mockey.UnPatchAll()

		expected := &model.ExamOffering{ExamHash: "exam-hash", Tag: "exam-tag"}
		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Create).Return(mockDB).Build()

		result, err := NewDBCourse(mockDB, new(utils.Snowflake)).
			CreateExamOffering(context.Background(), expected)

		assert.NoError(t, err)
		assert.Same(t, expected, result)
	})

	t.Run("database error", func(t *testing.T) {
		defer mockey.UnPatchAll()

		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Create).Return(&gorm.DB{Error: fmt.Errorf("insert failed")}).Build()

		result, err := NewDBCourse(mockDB, new(utils.Snowflake)).
			CreateExamOffering(context.Background(), &model.ExamOffering{})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "dal.CreateExamOffering error")
	})

	t.Run("duplicated", func(t *testing.T) {
		defer mockey.UnPatchAll()

		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Create).Return(&gorm.DB{Error: gorm.ErrDuplicatedKey}).Build()

		result, err := NewDBCourse(mockDB, new(utils.Snowflake)).
			CreateExamOffering(context.Background(), &model.ExamOffering{})

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestDBCourse_DeleteExamOfferingByHash(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer mockey.UnPatchAll()

		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Where).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Unscoped).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Delete).Return(mockDB).Build()

		err := NewDBCourse(mockDB, new(utils.Snowflake)).
			DeleteExamOfferingByHash(context.Background(), "exam-hash")

		assert.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		defer mockey.UnPatchAll()

		mockDB := new(gorm.DB)
		mockey.Mock((*gorm.DB).WithContext).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Table).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Where).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Unscoped).Return(mockDB).Build()
		mockey.Mock((*gorm.DB).Delete).Return(&gorm.DB{Error: errors.New("delete failed")}).Build()

		err := NewDBCourse(mockDB, new(utils.Snowflake)).
			DeleteExamOfferingByHash(context.Background(), "exam-hash")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dal.DeleteExamOfferingByHash error")
	})
}
