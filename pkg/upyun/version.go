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

package upyun

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// gmtDate returns the current date and time in GMT format.
func gmtDate() string {
	return time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

// SignStr generates the signature string for authentication.
func SignStr(policy string) string {
	// Generate MD5 hash of the password
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(config.VersionUploadService.Pass))
	key := fmt.Sprintf("%x", md5Hasher.Sum(nil))

	gmtdate := gmtDate()
	var msg string
	if policy == "" {
		msg = "POST" + "&/" + config.VersionUploadService.Bucket + "&" + gmtdate
	} else {
		msg = "POST" + "&/" + config.VersionUploadService.Bucket + "&" + gmtdate + "&" + policy
	}

	// Generate HMAC-SHA1 hash
	hmacHasher := hmac.New(sha1.New, []byte(key))
	hmacHasher.Write([]byte(msg))
	signature := base64.StdEncoding.EncodeToString(hmacHasher.Sum(nil))

	return "UPYUN " + config.VersionUploadService.Operator + ":" + signature
}

// GetPolicy generates the policy string for requests.
func GetPolicy() string {
	gmtdate := gmtDate()
	expiration := time.Now().Unix() + config.VersionUploadService.TokenTimeout
	// expiration := timeout
	policy := map[string]interface{}{
		"bucket":     config.VersionUploadService.Bucket,
		"save-key":   config.VersionUploadService.Path,
		"expiration": expiration,
		"date":       gmtdate,
	}

	policyJSON, _ := json.Marshal(policy)
	return base64.StdEncoding.EncodeToString(policyJSON)
}

// URlUploadFile 又拍云上传文件
func URlUploadFile(file []byte, url string) error {
	body := bytes.NewReader(file)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.UpYun.Operator, config.UpYun.Password)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warnf("URlUploadFile : failed to close response body: %v", err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return errno.UpcloudError
	}
	return nil
}

// URlGetFile 又拍云下载文件
func URlGetFile(url string) (*[]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.UpYun.Operator, config.UpYun.Password)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warnf("URlGetFile : failed to close response body: %v", err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errno.UpcloudError
	}

	file, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errno.UpcloudError
	}
	return &file, nil
}

// JoinFileName 生成文件名字
func JoinFileName(fileName string) string {
	return strings.Join([]string{
		config.UpYun.UssDomain, config.UpYun.Path,
		fileName,
	}, "")
}
