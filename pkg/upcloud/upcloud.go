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

package upcloud

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// UploadImg 又拍云上传文件
func UploadImg(file []byte, name string) error {
	body := bytes.NewReader(file)
	url := config.Upcloud.DomainName + config.Upcloud.Path + name

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.Upcloud.User, config.Upcloud.Pass)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errno.UpcloudError
	}
	return nil
}

// DeleteImg 又拍云删除文件
func DeleteImg(name string) error {
	// body := bytes.NewReader(file)
	url := config.Upcloud.DomainName + config.Upcloud.Path + name

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.Upcloud.User, config.Upcloud.Pass)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errno.UpcloudError
	}
	return nil
}

// GenerateImgName 生成图片名字
func GenerateImgName(uid int64) string {
	currentTime := time.Now()
	// 获取年月日和小时分钟
	year, month, day := currentTime.Date()
	hour, minute := currentTime.Hour(), currentTime.Minute()
	second := currentTime.Second()
	return fmt.Sprintf("%v_%d%02d%02d_%02d%02d%02d.jpg", uid, year, month, day, hour, minute, second)
}
