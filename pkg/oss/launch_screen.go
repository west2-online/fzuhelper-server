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
	"fmt"
	"strings"

	"github.com/upyun/go-sdk/v3/upyun"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// LaunchScreenOSSCli 根据需求定制Cli
type LaunchScreenOSSCli struct {
	upYun          *upyun.UpYun
	ussDomain      string
	path           string
	downloadDomain string
	sf             *utils.Snowflake
}

func NewLaunchScreenOSSCli(cfg *UpYunConfig, sf *utils.Snowflake) LaunchScreenOSSRepo {
	return &LaunchScreenOSSCli{
		upYun:          cfg.upyun,
		ussDomain:      cfg.UssDomain,
		path:           cfg.Path,
		downloadDomain: cfg.DownloadDomain,
		sf:             sf,
	}
}

// UploadImg 又拍云上传文件到指定path
func (c *LaunchScreenOSSCli) UploadImg(file []byte, remotePath string) error {
	err := c.upYun.Put(&upyun.PutObjectConfig{
		Path:   remotePath,
		Reader: bytes.NewReader(file),
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteImg 又拍云删除指定path的文件
func (c *LaunchScreenOSSCli) DeleteImg(remotePath string) error {
	err := c.upYun.Delete(&upyun.DeleteObjectConfig{
		Path: remotePath,
	})
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
	return strings.TrimPrefix(strings.TrimPrefix(url, c.downloadDomain), c.ussDomain)
}
