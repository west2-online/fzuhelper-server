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

package toolbox

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func mockToolboxGormChain(db *gorm.DB) {
	mockey.Mock((*gorm.DB).WithContext).Return(db).Build()
	mockey.Mock((*gorm.DB).Model).Return(db).Build()
	mockey.Mock((*gorm.DB).Table).Return(db).Build()
	mockey.Mock((*gorm.DB).Where).Return(db).Build()
}

func TestDBToolbox_CreateToolboxConfig(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("success", t, func() {
		db := new(gorm.DB)
		mockToolboxGormChain(db)
		mockey.Mock((*gorm.DB).Create).To(func(value interface{}) *gorm.DB {
			value.(*model.ToolboxConfig).Id = 123
			return db
		}).Build()

		config := &model.ToolboxConfig{ToolID: 1}
		err := NewDBToolbox(db, nil).CreateToolboxConfig(context.Background(), config)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), config.Id)
	})

	mockey.PatchConvey("duplicate", t, func() {
		db := &gorm.DB{Error: gorm.ErrDuplicatedKey}
		mockToolboxGormChain(db)
		mockey.Mock((*gorm.DB).Create).Return(db).Build()

		err := NewDBToolbox(db, nil).CreateToolboxConfig(context.Background(), &model.ToolboxConfig{ToolID: 1})
		assert.Equal(t, int64(errno.BizLogicCode), errno.ConvertErr(err).ErrorCode)
	})
}

func TestDBToolbox_GetToolboxConfigByID(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("success", t, func() {
		db := new(gorm.DB)
		mockToolboxGormChain(db)
		mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
			config := dest.(*model.ToolboxConfig)
			config.Id = 123
			config.ToolID = 1
			return db
		}).Build()

		config, err := NewDBToolbox(db, nil).GetToolboxConfigByID(context.Background(), 123)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), config.Id)
	})

	mockey.PatchConvey("not found", t, func() {
		db := &gorm.DB{Error: gorm.ErrRecordNotFound}
		mockToolboxGormChain(db)
		mockey.Mock((*gorm.DB).First).Return(db).Build()

		_, err := NewDBToolbox(db, nil).GetToolboxConfigByID(context.Background(), 123)
		assert.Equal(t, int64(errno.BizNotExist), errno.ConvertErr(err).ErrorCode)
	})
}

func TestDBToolbox_UpdateToolboxConfig(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("zero values are included and updated row is returned", t, func() {
		db := new(gorm.DB)
		mockToolboxGormChain(db)
		mockey.Mock((*gorm.DB).Updates).To(func(values interface{}) *gorm.DB {
			updates := values.(map[string]any)
			assert.Equal(t, false, updates["visible"])
			assert.Equal(t, "", updates["name"])
			assert.Equal(t, int64(0), updates["version"])
			return &gorm.DB{RowsAffected: 1}
		}).Build()
		expected := &model.ToolboxConfig{Id: 123, ToolID: 1}
		mockey.Mock((*DBToolbox).GetToolboxConfigByID).Return(expected, nil).Build()

		config, err := NewDBToolbox(db, nil).UpdateToolboxConfig(
			context.Background(),
			123,
			&model.ToolboxConfig{ToolID: 1},
		)
		assert.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}

func TestDBToolbox_DeleteToolboxConfig(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("success", t, func() {
		db := new(gorm.DB)
		mockToolboxGormChain(db)
		mockey.Mock((*gorm.DB).Delete).Return(&gorm.DB{RowsAffected: 1}).Build()
		assert.NoError(t, NewDBToolbox(db, nil).DeleteToolboxConfig(context.Background(), 123))
	})

	mockey.PatchConvey("not found", t, func() {
		db := new(gorm.DB)
		mockToolboxGormChain(db)
		mockey.Mock((*gorm.DB).Delete).Return(&gorm.DB{}).Build()
		err := NewDBToolbox(db, nil).DeleteToolboxConfig(context.Background(), 123)
		assert.Equal(t, int64(errno.BizNotExist), errno.ConvertErr(err).ErrorCode)
	})
}
