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

package cos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
)

// setupCosConfig 初始化指定服务的测试配置,并覆盖为合法的桶名/地域,
// 避免 NewBucketURL 校验失败触发 Fatalf
func setupCosConfig(t *testing.T, service string) {
	t.Helper()
	if err := config.InitForTest(service); err != nil {
		t.Fatalf("setupCosConfig: init config failed: %v", err)
	}
	config.Cos.Bucket = "fzuhelper-test-bucket-125000000"
	config.Cos.Region = "ap-guangzhou"
}

func TestObjectKey(t *testing.T) {
	testCases := []struct {
		name string
		path string
		want string
	}{
		{"WithLeadingSlash", "/statistic/version.json", "statistic/version.json"},
		{"WithoutLeadingSlash", "statistic/version.json", "statistic/version.json"},
		{"RootPath", "/", ""},
		{"EmptyPath", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, objectKey(tc.path))
		})
	}
}

func TestKeyFromURL(t *testing.T) {
	key, err := keyFromURL("https://download.w2fzu.com/statistic/version.json")
	assert.Nil(t, err)
	assert.Equal(t, "statistic/version.json", key)

	// 无路径的 URL 返回空 key
	key, err = keyFromURL("https://download.w2fzu.com")
	assert.Nil(t, err)
	assert.Equal(t, "", key)

	// 非法 URL
	_, err = keyFromURL("://invalid")
	assert.NotNil(t, err)
}
