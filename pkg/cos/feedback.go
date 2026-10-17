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
	"bytes"
	"context"
	"fmt"
	"path"
	"strings"

	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const (
	FeedbackImageCategory    = "img"
	FeedbackLogCategory      = "log"
	FeedbackLogFileExtension = "json.gz"
)

type FeedbackCOSRepo interface {
	Upload(file []byte, remotePath string) error
	GenerateFileName(category, suffix string) (url, remotePath string, err error)
}

type FeedbackCOSCli struct {
	client         *tencentyun.Client
	path           string
	downloadDomain string
	sf             *utils.Snowflake
}

func NewFeedbackCOSCli(client *tencentyun.Client, rootPath, downloadDomain string, sf *utils.Snowflake) FeedbackCOSRepo {
	return &FeedbackCOSCli{
		client:         client,
		path:           rootPath,
		downloadDomain: downloadDomain,
		sf:             sf,
	}
}

func (c *FeedbackCOSCli) Upload(file []byte, remotePath string) error {
	_, err := c.client.Object.Put(context.Background(), objectKey(remotePath), bytes.NewReader(file), nil)
	if err != nil {
		return errno.UpcloudError
	}
	return nil
}

func (c *FeedbackCOSCli) GenerateFileName(category, suffix string) (string, string, error) {
	id, err := c.sf.NextVal()
	if err != nil {
		return "", "", errno.Errorf(errno.InternalSFErrorCode, "failed to generate feedback file name: %v", err)
	}

	fileName := fmt.Sprintf("%d.%s", id, suffix)
	remotePath := "/" + strings.TrimLeft(path.Join(c.path, category, fileName), "/")
	url := strings.TrimRight(c.downloadDomain, "/") + remotePath
	return url, remotePath, nil
}
