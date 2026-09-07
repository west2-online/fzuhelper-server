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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestJoinFileName(t *testing.T) {
	if err := config.InitForTest("version"); err != nil {
		t.Fatalf("TestJoinFileName: init config failed: %v", err)
	}

	result := JoinFileName("version.json")
	assert.Equal(t, config.Cos.DownloadDomain+config.Cos.Path+"version.json", result)
}

func TestURlGetFile(t *testing.T) {
	setupCosConfig(t, "version")
	const fileURL = "https://download.w2fzu.com/statistic/version.json"

	defer mockey.UnPatchAll()
	mockey.PatchConvey("GetSuccess", t, func() {
		var gotKey string
		mockey.Mock((*tencentyun.ObjectService).Get).To(
			func(_ context.Context, name string, _ *tencentyun.ObjectGetOptions, _ ...string) (*tencentyun.Response, error) {
				gotKey = name
				return &tencentyun.Response{Response: &http.Response{
					Body: io.NopCloser(strings.NewReader("file-content")),
				}}, nil
			}).Build()

		file, err := URlGetFile(fileURL)
		assert.Nil(t, err)
		// key 应从 URL 解析且不带前导 '/'
		assert.Equal(t, "statistic/version.json", gotKey)
		assert.Equal(t, "file-content", string(*file))
	})

	mockey.PatchConvey("GetError", t, func() {
		mockey.Mock((*tencentyun.ObjectService).Get).Return(nil, assert.AnError).Build()

		file, err := URlGetFile(fileURL)
		assert.Equal(t, errno.UpcloudError, err)
		assert.Nil(t, file)
	})
}

func TestURlUploadFile(t *testing.T) {
	setupCosConfig(t, "version")
	const fileURL = "https://download.w2fzu.com/statistic/version.json"

	defer mockey.UnPatchAll()
	mockey.PatchConvey("PutSuccess", t, func() {
		var (
			gotKey  string
			gotBody string
		)
		mockey.Mock((*tencentyun.ObjectService).Put).To(
			func(_ context.Context, name string, r io.Reader, _ *tencentyun.ObjectPutOptions) (*tencentyun.Response, error) {
				gotKey = name
				body, _ := io.ReadAll(r)
				gotBody = string(body)
				return nil, nil
			}).Build()

		err := URlUploadFile([]byte("upload-content"), fileURL)
		assert.Nil(t, err)
		assert.Equal(t, "statistic/version.json", gotKey)
		assert.Equal(t, "upload-content", gotBody)
	})

	mockey.PatchConvey("PutError", t, func() {
		mockey.Mock((*tencentyun.ObjectService).Put).Return(nil, assert.AnError).Build()

		err := URlUploadFile([]byte("upload-content"), fileURL)
		assert.Equal(t, errno.UpcloudError, err)
	})
}

// policyConditions 将 policy 的 conditions 解析为 map 列表
func policyConditions(t *testing.T, policy string) []map[string]string {
	t.Helper()
	policyJSON, err := base64.StdEncoding.DecodeString(policy)
	assert.Nil(t, err)

	var p struct {
		Expiration string            `json:"expiration"`
		Conditions []json.RawMessage `json:"conditions"`
	}
	assert.Nil(t, json.Unmarshal(policyJSON, &p))

	conds := make([]map[string]string, 0, len(p.Conditions))
	for _, raw := range p.Conditions {
		var cond map[string]string
		if err := json.Unmarshal(raw, &cond); err == nil {
			conds = append(conds, cond)
		}
	}
	return conds
}

func findCondition(conds []map[string]string, key string) string {
	for _, cond := range conds {
		if v, ok := cond[key]; ok {
			return v
		}
	}
	return ""
}

func TestGetPolicy(t *testing.T) {
	if err := config.InitForTest("version"); err != nil {
		t.Fatalf("TestGetPolicy: init config failed: %v", err)
	}

	policy := GetPolicy()
	conds := policyConditions(t, policy)

	assert.Equal(t, config.CosUpload.Bucket, findCondition(conds, "bucket"))
	assert.Equal(t, "sha1", findCondition(conds, "q-sign-algorithm"))
	assert.Equal(t, config.CosUpload.SecretID, findCondition(conds, "q-ak"))

	// q-sign-time 为 "start;end" 且窗口等于 TokenTimeout
	keyTime := findCondition(conds, "q-sign-time")
	segments := strings.Split(keyTime, ";")
	assert.Len(t, segments, 2)
	start, err := strconv.ParseInt(segments[0], 10, 64)
	assert.Nil(t, err)
	end, err := strconv.ParseInt(segments[1], 10, 64)
	assert.Nil(t, err)
	assert.Equal(t, config.CosUpload.TokenTimeout, end-start)
}

func TestSignStr(t *testing.T) {
	if err := config.InitForTest("version"); err != nil {
		t.Fatalf("TestSignStr: init config failed: %v", err)
	}

	const testSecretID = "AKIDtest"
	const testSecretKey = "test-key"
	config.CosUpload.SecretID = testSecretID
	config.CosUpload.SecretKey = testSecretKey
	defer func() {
		config.CosUpload.SecretID = ""
		config.CosUpload.SecretKey = ""
	}()

	policy := GetPolicy()
	auth := SignStr(policy)

	// 解析 authorization 串
	fields := make(map[string]string)
	for _, kv := range strings.Split(auth, "&") {
		parts := strings.SplitN(kv, "=", 2)
		assert.Len(t, parts, 2)
		fields[parts[0]] = parts[1]
	}
	assert.Equal(t, "sha1", fields["q-sign-algorithm"])
	assert.Equal(t, testSecretID, fields["q-ak"])

	// 签名的 key-time 必须与 policy 中的 q-sign-time 一致,否则 COS 会拒绝
	keyTime := fields["q-sign-time"]
	assert.Equal(t, keyTime, findCondition(policyConditions(t, policy), "q-sign-time"))
	assert.Equal(t, keyTime, fields["q-key-time"])

	// 按 COS PostObject 算法独立复算签名:
	// signKey = hex(hmac_sha1(SecretKey, keyTime))
	// signature = hex(hmac_sha1(signKey, "sha1\n<keyTime>\n<sha1hex(policy)>\n"))
	h := hmac.New(sha1.New, []byte(testSecretKey))
	h.Write([]byte(keyTime))
	signKey := hex.EncodeToString(h.Sum(nil))

	s := sha1.New()
	s.Write([]byte(policy))
	policySha1 := hex.EncodeToString(s.Sum(nil))
	stringToSign := strings.Join([]string{"sha1", keyTime, policySha1, ""}, "\n")

	h2 := hmac.New(sha1.New, []byte(signKey))
	h2.Write([]byte(stringToSign))
	expected := hex.EncodeToString(h2.Sum(nil))

	assert.Equal(t, expected, fields["q-signature"])
}
