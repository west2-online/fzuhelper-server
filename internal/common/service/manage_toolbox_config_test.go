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

package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/toolbox"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func newToolboxTestService() *CommonService {
	return NewCommonService(context.Background(), &base.ClientSet{
		DBClient: new(db.Database),
	}, new(taskqueue.BaseTaskQueue))
}

func validToolboxConfig() *model.ToolboxConfig {
	return &model.ToolboxConfig{
		ToolID:    1,
		Visible:   true,
		Name:      new("Tool"),
		Icon:      new("icon"),
		Type:      new("web"),
		Message:   new("message"),
		Extra:     new("extra"),
		StudentID: new("102300217"),
		Platform:  new("android"),
		Version:   new(int64(1)),
	}
}

func TestCreateToolboxConfig(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("success", t, func() {
		mockey.Mock(utils.CheckPwd).Return(true).Build()
		mockey.Mock((*toolbox.DBToolbox).CreateToolboxConfig).To(
			func(_ context.Context, config *model.ToolboxConfig) error {
				config.Id = 123
				return nil
			},
		).Build()

		result, err := newToolboxTestService().CreateToolboxConfig(context.Background(), "secret", validToolboxConfig())
		assert.NoError(t, err)
		assert.Equal(t, int64(123), result.Id)
	})

	mockey.PatchConvey("validation and database errors", t, func() {
		mockey.Mock(utils.CheckPwd).Return(true).Build()
		service := newToolboxTestService()

		_, err := service.CreateToolboxConfig(context.Background(), "secret", &model.ToolboxConfig{})
		assert.ErrorContains(t, err, "tool_id must be positive")

		config := validToolboxConfig()
		config.Version = new(int64(MaxVersionNumber + 1))
		_, err = service.CreateToolboxConfig(context.Background(), "secret", config)
		assert.ErrorContains(t, err, "version cannot exceed")

		mockey.Mock((*toolbox.DBToolbox).CreateToolboxConfig).Return(assert.AnError).Build()
		_, err = service.CreateToolboxConfig(context.Background(), "secret", validToolboxConfig())
		assert.ErrorContains(t, err, "service.CreateToolboxConfig")
	})

	mockey.PatchConvey("invalid secret", t, func() {
		mockey.Mock(utils.CheckPwd).Return(false).Build()
		_, err := newToolboxTestService().CreateToolboxConfig(context.Background(), "wrong", validToolboxConfig())
		assert.ErrorContains(t, err, "invalid admin secret")
	})
}

func TestGetUpdateDeleteToolboxConfigByID(t *testing.T) {
	defer mockey.UnPatchAll()

	mockey.PatchConvey("get success", t, func() {
		expected := validToolboxConfig()
		expected.Id = 123
		mockey.Mock(utils.CheckPwd).Return(true).Build()
		mockey.Mock((*toolbox.DBToolbox).GetToolboxConfigByID).Return(expected, nil).Build()

		result, err := newToolboxTestService().GetToolboxConfigByID(context.Background(), "secret", 123)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	mockey.PatchConvey("update preserves explicit zero values", t, func() {
		expected := validToolboxConfig()
		expected.Id = 123
		expected.Visible = false
		expected.Name = nil
		expected.Version = nil
		mockey.Mock(utils.CheckPwd).Return(true).Build()
		mockey.Mock((*toolbox.DBToolbox).UpdateToolboxConfig).To(
			func(_ context.Context, id int64, config *model.ToolboxConfig) (*model.ToolboxConfig, error) {
				assert.Equal(t, int64(123), id)
				assert.False(t, config.Visible)
				assert.Nil(t, config.Name)
				assert.Nil(t, config.Version)
				return expected, nil
			},
		).Build()

		result, err := newToolboxTestService().UpdateToolboxConfig(context.Background(), "secret", 123, expected)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	mockey.PatchConvey("delete success", t, func() {
		mockey.Mock(utils.CheckPwd).Return(true).Build()
		mockey.Mock((*toolbox.DBToolbox).DeleteToolboxConfig).Return(nil).Build()
		assert.NoError(t, newToolboxTestService().DeleteToolboxConfig(context.Background(), "secret", 123))
	})

	mockey.PatchConvey("invalid id", t, func() {
		mockey.Mock(utils.CheckPwd).Return(true).Build()
		service := newToolboxTestService()

		_, err := service.GetToolboxConfigByID(context.Background(), "secret", 0)
		assert.ErrorContains(t, err, "config_id must be positive")
		_, err = service.UpdateToolboxConfig(context.Background(), "secret", -1, validToolboxConfig())
		assert.ErrorContains(t, err, "config_id must be positive")
		err = service.DeleteToolboxConfig(context.Background(), "secret", 0)
		assert.ErrorContains(t, err, "config_id must be positive")
	})

	mockey.PatchConvey("not found remains BizNotExist", t, func() {
		notFound := errno.NewErrNo(errno.BizNotExist, "toolbox config not found")
		mockey.Mock(utils.CheckPwd).Return(true).Build()
		mockey.Mock((*toolbox.DBToolbox).GetToolboxConfigByID).Return(nil, notFound).Build()

		_, err := newToolboxTestService().GetToolboxConfigByID(context.Background(), "secret", 123)
		converted := errno.ConvertErr(err)
		assert.Equal(t, int64(errno.BizNotExist), converted.ErrorCode)
	})
}
