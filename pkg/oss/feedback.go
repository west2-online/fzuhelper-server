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
	"path"
	"strings"

	"github.com/upyun/go-sdk/v3/upyun"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

type FeedbackOSSCli struct {
	upYun          *upyun.UpYun
	path           string
	downloadDomain string
	sf             *utils.Snowflake
}

func NewFeedbackOSSCli(cfg *UpYunConfig, sf *utils.Snowflake) FeedbackOSSRepo {
	return &FeedbackOSSCli{
		upYun:          cfg.upyun,
		path:           cfg.Path,
		downloadDomain: cfg.DownloadDomain,
		sf:             sf,
	}
}

func (c *FeedbackOSSCli) Upload(file []byte, remotePath string) error {
	if err := c.upYun.Put(&upyun.PutObjectConfig{
		Path:   remotePath,
		Reader: bytes.NewReader(file),
	}); err != nil {
		return errno.UpcloudError
	}
	return nil
}

func (c *FeedbackOSSCli) GenerateFileName(category, suffix string) (string, string, error) {
	id, err := c.sf.NextVal()
	if err != nil {
		return "", "", errno.Errorf(errno.InternalSFErrorCode, "failed to generate feedback file name: %v", err)
	}

	fileName := fmt.Sprintf("%d.%s", id, suffix)
	remotePath := "/" + strings.TrimLeft(path.Join(c.path, category, fileName), "/")
	url := strings.TrimRight(c.downloadDomain, "/") + remotePath
	return url, remotePath, nil
}
