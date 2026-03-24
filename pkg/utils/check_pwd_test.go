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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
)

func TestCheckPwd(t *testing.T) {
	tests := []struct {
		name   string
		pwd    string
		secret string
		want   bool
	}{
		{
			name:   "Success",
			pwd:    "114514",
			secret: "114514",
			want:   true,
		},
		{
			name:   "Invalid",
			pwd:    "114514",
			secret: "1919810",
			want:   false,
		},
	}
	_ = config.InitForTest("api")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Admin.Secret = tt.secret
			result := CheckPwd(tt.pwd)
			assert.Equal(t, result, tt.want)
		})
	}
}
