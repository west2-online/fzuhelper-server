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
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// gmtDate returns the current date and time in GMT format.
func gmtDate() string {
	return time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

// SignStr generates the signature string for authentication.
func SignStr(opename, opepass, bucket, policy string) string {
	// Generate MD5 hash of the password
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(opepass))
	key := fmt.Sprintf("%x", md5Hasher.Sum(nil))

	gmtdate := gmtDate()
	var msg string
	if policy == "" {
		msg = "POST" + "&/" + bucket + "&" + gmtdate
	} else {
		msg = "POST" + "&/" + bucket + "&" + gmtdate + "&" + policy
	}

	// Generate HMAC-SHA1 hash
	hmacHasher := hmac.New(sha1.New, []byte(key))
	hmacHasher.Write([]byte(msg))
	signature := base64.StdEncoding.EncodeToString(hmacHasher.Sum(nil))

	return "UPYUN " + opename + ":" + signature
}

// GetPolicy generates the policy string for requests.
func GetPolicy(bucket, savepath string, timeout int) string {
	gmtdate := gmtDate()
	expiration := time.Now().Unix() + int64(timeout)
	// expiration := timeout
	policy := map[string]interface{}{
		"bucket":     bucket,
		"save-key":   savepath,
		"expiration": expiration,
		"date":       gmtdate,
	}

	policyJSON, _ := json.Marshal(policy)
	return base64.StdEncoding.EncodeToString(policyJSON)
}
