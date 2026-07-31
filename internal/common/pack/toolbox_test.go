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

package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func TestBuildToolboxConfigPreservesNullableValues(t *testing.T) {
	empty := ""
	zero := int64(0)

	config := BuildToolboxConfig(&dbmodel.ToolboxConfig{
		ToolID:  1,
		Visible: false,
		Name:    &empty,
		Version: &zero,
	})

	assert.NotNil(t, config.Visible)
	assert.False(t, *config.Visible)
	assert.NotNil(t, config.Name)
	assert.Empty(t, *config.Name)
	assert.Nil(t, config.Icon)
	assert.NotNil(t, config.Version)
	assert.Zero(t, *config.Version)
}

func TestBuildToolboxConfigDetailPreservesSQLNull(t *testing.T) {
	config := BuildToolboxConfigDetail(&dbmodel.ToolboxConfig{
		Id:      10,
		ToolID:  1,
		Visible: true,
	})

	assert.Equal(t, int64(10), config.ConfigId)
	assert.True(t, config.Visible)
	assert.Nil(t, config.Name)
	assert.Nil(t, config.StudentId)
	assert.Nil(t, config.Version)
}
