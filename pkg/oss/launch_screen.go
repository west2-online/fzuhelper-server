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

package oss

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// LaunchScreenOSSCli 根据需求定制Cli
type LaunchScreenOSSCli struct {
	client         *tencentyun.Client
	downloadDomain string
	path           string
	sf             *utils.Snowflake
}

func NewLaunchScreenOSSCli(cfg *CosConfig, sf *utils.Snowflake) LaunchScreenOSSRepo {
	return &LaunchScreenOSSCli{
		client:         cfg.client,
		downloadDomain: cfg.DownloadDomain,
		path:           cfg.Path,
		sf:             sf,
	}
}

// UploadImg 上传文件到指定path
func (c *LaunchScreenOSSCli) UploadImg(file []byte, remotePath string) error {
	// COS 对象 key 不带前导 '/'
	_, err := c.client.Object.Put(context.Background(), strings.TrimPrefix(remotePath, "/"), bytes.NewReader(file), nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteImg 删除指定path的文件
func (c *LaunchScreenOSSCli) DeleteImg(remotePath string) error {
	_, err := c.client.Object.Delete(context.Background(), strings.TrimPrefix(remotePath, "/"), nil)
	if err != nil {
		return err
	}
	return nil
}

// GenerateImgName 生成唯一图片名字
func (c *LaunchScreenOSSCli) GenerateImgName(suffix string) (string, string, error) {
	// 唯一id
	sfid, err := c.sf.NextVal()
	if err != nil {
		return "", "", errno.Errorf(errno.InternalSFErrorCode, "failed to generate next val:%v", err)
	}
	newFileName := fmt.Sprintf(
		"%d.%s",
		sfid,
		suffix,
	)
	remotePath := strings.Join([]string{
		c.path,
		newFileName,
	}, "")

	return strings.Join([]string{
		c.downloadDomain,
		c.path,
		newFileName,
	}, ""), remotePath, nil
}

// GetRemotePathFromUrl 获得远程path
func (c *LaunchScreenOSSCli) GetRemotePathFromUrl(url string) string {
	return strings.TrimPrefix(url, c.downloadDomain)
}
