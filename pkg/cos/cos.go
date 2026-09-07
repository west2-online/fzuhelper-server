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

// Package cos 腾讯云 COS 对象存储操作封装
package cos

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const defaultRequestTimeout = 30 * time.Second

var (
	client *tencentyun.Client
	once   sync.Once
)

// NewCos 使用当前服务对应的 coss 配置显式初始化 COS 客户端
func NewCos() *tencentyun.Client {
	once.Do(func() {
		client = newClient()
	})
	return client
}

// GetClient 返回 COS 客户端,未初始化时懒加载
func GetClient() *tencentyun.Client {
	return NewCos()
}

func newClient() *tencentyun.Client {
	if config.Cos == nil {
		logger.Fatalf("cos.newClient: coss config not loaded for this service")
	}
	// 统一通过 NewBucketURL 构造桶地址,避免手拼 URL
	bucketURL, err := tencentyun.NewBucketURL(config.Cos.Bucket, config.Cos.Region, true)
	if err != nil {
		logger.Fatalf("cos.newClient: failed to create bucket URL: %v", err)
	}
	return tencentyun.NewClient(&tencentyun.BaseURL{BucketURL: bucketURL}, &http.Client{
		Timeout: defaultRequestTimeout,
		Transport: &tencentyun.AuthorizationTransport{
			SecretID:  config.Cos.SecretID,
			SecretKey: config.Cos.SecretKey,
		},
	})
}

// objectKey 转换为 COS 对象 key(去掉前导 '/')
func objectKey(path string) string {
	return strings.TrimPrefix(path, "/")
}

// keyFromURL 从下载 URL 中解析出对象 key
func keyFromURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return objectKey(u.Path), nil
}
