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
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// listMaxKeys COS 单次列举条数上限
const listMaxKeys = 1000

// GetDir 获取目录下的文件和文件夹
func GetDir(path string) (*model.UpYunFileDir, error) {
	prefix := objectKey(path)
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	fileDir := &model.UpYunFileDir{
		BasePath: &path,
		Folders:  []string{},
		Files:    []string{},
	}

	client := GetClient()
	marker := ""
	for {
		res, _, err := client.Bucket.Get(context.Background(), &tencentyun.BucketGetOptions{
			Prefix:    prefix,
			Delimiter: "/",
			Marker:    marker,
			MaxKeys:   listMaxKeys,
		})
		if err != nil {
			return nil, err
		}
		for _, commonPrefix := range res.CommonPrefixes {
			name := strings.TrimSuffix(strings.TrimPrefix(commonPrefix, prefix), "/")
			// 过滤以 __ 开头的特殊文件夹
			if name != "" && !strings.HasPrefix(name, "__") {
				fileDir.Folders = append(fileDir.Folders, name)
			}
		}
		for _, obj := range res.Contents {
			name := strings.TrimPrefix(obj.Key, prefix)
			// 跳过目录占位对象(key 以 '/' 结尾)
			if name == "" || strings.HasSuffix(obj.Key, "/") {
				continue
			}
			fileDir.Files = append(fileDir.Files, name)
		}
		if !res.IsTruncated {
			break
		}
		marker = res.NextMarker
	}

	return fileDir, nil
}

// GetDownloadUrl 基于路径获取对应的下载链接
// 配置了 token-secret 时按 EO TypeA Token 鉴权拼接防盗链参数,留空则不签名
func GetDownloadUrl(uri string) (string, error) {
	base := config.Cos.DownloadDomain + utils.UriEncode(uri)
	if config.Cos.TokenSecret == "" {
		return base, nil
	}

	etime := time.Now().Unix() + config.Cos.TokenTimeout
	randStr, uidStr := "0", "0"
	sign := utils.MD5(strings.Join([]string{
		uri, // 参与签名的 uri 为未编码的原始路径
		strconv.FormatInt(etime, 10),
		randStr,
		uidStr,
		config.Cos.TokenSecret,
	}, "-"))
	return fmt.Sprintf("%s?sign=%s-%d-%s-%s", base, sign, etime, randStr, uidStr), nil
}

// UploadFile 上传本地文件到指定路径
func UploadFile(filepath, ussDir string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = GetClient().Object.Put(context.Background(), objectKey(ussDir), f, nil)
	return err
}
