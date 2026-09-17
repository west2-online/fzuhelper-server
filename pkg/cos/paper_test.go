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
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestGetDir(t *testing.T) {
	setupCosConfig(t, "paper")

	const dirPath = "/__Unchecked_adjowe"
	// GetDir 会去掉前导 '/' 作为列举前缀
	const prefix = "__Unchecked_adjowe/"

	page1 := &tencentyun.BucketGetResult{
		CommonPrefixes: []string{
			prefix + "C语言/",
			prefix + "王鸿/",
			prefix + "__Unchecked_tmp/", // 以 __ 开头的文件夹应被过滤
		},
		Contents: []tencentyun.Object{
			{Key: prefix + "a.pdf"},
			{Key: prefix + "placeholder/"}, // 目录占位对象应被跳过
		},
		IsTruncated: true,
		NextMarker:  "marker-1",
	}
	page2 := &tencentyun.BucketGetResult{
		Contents:    []tencentyun.Object{{Key: prefix + "b.zip"}},
		IsTruncated: false,
	}

	defer mockey.UnPatchAll()
	mockey.PatchConvey("ListWithPaginationAndFilter", t, func() {
		var markers []string
		mockey.Mock((*tencentyun.BucketService).Get).To(
			func(_ context.Context, opt *tencentyun.BucketGetOptions) (*tencentyun.BucketGetResult, *tencentyun.Response, error) {
				assert.Equal(t, prefix, opt.Prefix)
				assert.Equal(t, "/", opt.Delimiter)
				markers = append(markers, opt.Marker)
				if len(markers) == 1 {
					return page1, nil, nil
				}
				return page2, nil, nil
			}).Build()

		fileDir, err := GetDir(dirPath)
		assert.Nil(t, err)
		assert.Equal(t, dirPath, *fileDir.BasePath)
		assert.Equal(t, []string{"C语言", "王鸿"}, fileDir.Folders)
		assert.Equal(t, []string{"a.pdf", "b.zip"}, fileDir.Files)
		// 第二次请求应携带上一页的 NextMarker
		assert.Equal(t, []string{"", "marker-1"}, markers)
	})

	mockey.PatchConvey("ListError", t, func() {
		mockey.Mock((*tencentyun.BucketService).Get).Return(nil, nil, assert.AnError).Build()

		fileDir, err := GetDir(dirPath)
		assert.NotNil(t, err)
		assert.Nil(t, fileDir)
	})
}

func TestGetDownloadUrl(t *testing.T) {
	setupCosConfig(t, "paper")
	const uri = "/__Unchecked_adjowe/C语言/10份练习.zip"
	const tokenTimeout = int64(60)
	const timeTolerance = int64(2)

	t.Run("UnsignedWhenTokenSecretEmpty", func(t *testing.T) {
		config.Cos.TokenSecret = ""
		config.Cos.TokenTimeout = tokenTimeout

		result, err := GetDownloadUrl(uri)
		assert.Nil(t, err)
		assert.Equal(t, config.Cos.DownloadDomain+utils.UriEncode(uri), result)
		assert.NotContains(t, result, "sign=")
	})

	t.Run("SignedWithEOTypeA", func(t *testing.T) {
		config.Cos.TokenSecret = "test-secret"
		config.Cos.TokenTimeout = tokenTimeout
		defer func() { config.Cos.TokenSecret = "" }()

		now := time.Now().Unix()
		result, err := GetDownloadUrl(uri)
		assert.Nil(t, err)

		// 格式: <domain><encoded-uri>?sign=<md5>-<etime>-<rand>-<uid>
		base := config.Cos.DownloadDomain + utils.UriEncode(uri)
		assert.True(t, strings.HasPrefix(result, base+"?sign="))

		signParam := strings.TrimPrefix(result, base+"?sign=")
		parts := strings.Split(signParam, "-")
		assert.Len(t, parts, 4)

		// 按算法独立复算签名,锁定 EO TypeA 计算方式
		etime, err := strconv.ParseInt(parts[1], 10, 64)
		assert.Nil(t, err)
		assert.InDelta(t, float64(now+tokenTimeout), float64(etime), float64(timeTolerance))
		expectedSign := utils.MD5(strings.Join([]string{uri, parts[1], parts[2], parts[3], config.Cos.TokenSecret}, "-"))
		assert.Equal(t, expectedSign, parts[0])
	})
}
