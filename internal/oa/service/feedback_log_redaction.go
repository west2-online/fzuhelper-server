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

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

const feedbackLogRedactedValue = "***"

var (
	feedbackKeyNormalizer = strings.NewReplacer("-", "", "_", "", " ", "")
	feedbackDigitPattern  = regexp.MustCompile(`[0-9]+`)
	feedbackBearerPattern = regexp.MustCompile(`(?i)\b(?:bearer|basic)\s+[^\s"',;]+`)
	feedbackJWTToken      = regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)
	feedbackURLPattern    = regexp.MustCompile(`(?i)(?:https?://[^\s<>"']+|/[^\s<>"'?]*\?[^\s<>"']+)`)
	feedbackHeaderPattern = regexp.MustCompile(`(?im)(^|\n)([ \t]*(?:authorization|proxy-authorization|cookie|cookies|set-cookie)[ \t]*:[ \t]*)[^\r\n]*`)
	// 先匹配完整键值，再按同一份字段表判断，避免把 token_count 等诊断字段误认为 token。
	feedbackAssignmentPattern = regexp.MustCompile(`([A-Za-z][A-Za-z0-9_-]*)(["']?\s*[:=]\s*)("(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|[^\s&,;"'\]}]+)`)
	feedbackSensitiveKeys     = map[string]struct{}{
		"authorization": {}, "proxyauthorization": {}, "accesstoken": {}, "refreshtoken": {},
		"cookie": {}, "cookies": {}, "setcookie": {}, "id": {}, "password": {}, "passwd": {},
		"secret": {}, "token": {}, "apikey": {}, "studentid": {}, "stuid": {}, "stunum": {},
		"phone": {}, "phonenumber": {}, "mobile": {}, "contactphone": {},
	}
)

// redactFeedbackLog 逐行校验、脱敏并写入 dst。失败时 dst 可能已有部分输出，调用方必须丢弃，不能上传。
func redactFeedbackLog(data []byte, dst io.Writer) error {
	if !utf8.Valid(data) {
		return fmt.Errorf("feedback log must be UTF-8 JSON Lines")
	}
	encoder := json.NewEncoder(dst)
	// 输出是日志 JSON，不嵌入 HTML；避免无谓转义 URL 中的 & 等字符。
	encoder.SetEscapeHTML(false)
	hasJSON := false
	for lineNumber := 1; len(data) > 0; lineNumber++ {
		var line []byte
		line, data, _ = bytes.Cut(data, []byte{'\n'})
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var value interface{}
		decoder := json.NewDecoder(bytes.NewReader(line))
		decoder.UseNumber()
		if err := decoder.Decode(&value); err != nil {
			return fmt.Errorf("invalid JSON at line %d", lineNumber)
		}
		// Decode 默认允许尾随另一个值；每行必须恰好包含一个 JSON 值。
		if err := decoder.Decode(new(interface{})); err != io.EOF {
			return fmt.Errorf("invalid JSON at line %d", lineNumber)
		}
		if err := encoder.Encode(redactFeedbackLogValue(value)); err != nil {
			return fmt.Errorf("marshal JSON line %d: %w", lineNumber, err)
		}
		hasJSON = true
	}
	if !hasJSON {
		return fmt.Errorf("feedback log contains no JSON values")
	}
	return nil
}

func redactFeedbackLogValue(value interface{}) interface{} {
	switch current := value.(type) {
	case map[string]interface{}:
		for key, item := range current {
			if isFeedbackLogSensitiveKey(key) {
				current[key] = feedbackLogRedactedValue
				continue
			}
			current[key] = redactFeedbackLogValue(item)
		}
	case []interface{}:
		for i, item := range current {
			current[i] = redactFeedbackLogValue(item)
		}
	case string:
		return redactFeedbackLogString(current)
	}
	// 无字段语义时保留数值及类型，不能按位数区分时间戳、耗时与学号。
	return value
}

func isFeedbackLogSensitiveKey(key string) bool {
	key = feedbackKeyNormalizer.Replace(strings.ToLower(key))
	_, ok := feedbackSensitiveKeys[key]
	return ok
}

func redactFeedbackLogURL(value string) string {
	parsed, err := url.Parse(value)
	if err != nil {
		// 解析失败时不保留可能含凭据的原始 URL。
		return feedbackLogRedactedValue
	}
	if parsed.User != nil {
		parsed.User = url.User(feedbackLogRedactedValue)
	}
	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		parsed.RawQuery = ""
	} else {
		for key := range query {
			if isFeedbackLogSensitiveKey(key) {
				query.Set(key, feedbackLogRedactedValue)
			}
		}
		parsed.RawQuery = query.Encode()
	}
	// fragment 不参与 HTTP 请求，可能携带 OAuth 凭据。
	parsed.Fragment = ""
	parsed.RawFragment = ""
	return parsed.String()
}

func redactFeedbackLogString(value string) string {
	value = feedbackHeaderPattern.ReplaceAllString(value, "${1}${2}"+feedbackLogRedactedValue)
	value = feedbackURLPattern.ReplaceAllStringFunc(value, redactFeedbackLogURL)
	value = feedbackBearerPattern.ReplaceAllString(value, feedbackLogRedactedValue)
	value = feedbackJWTToken.ReplaceAllString(value, feedbackLogRedactedValue)
	value = feedbackAssignmentPattern.ReplaceAllStringFunc(value, func(assignment string) string {
		parts := feedbackAssignmentPattern.FindStringSubmatch(assignment)
		if !isFeedbackLogSensitiveKey(parts[1]) {
			return assignment
		}
		replacement := feedbackLogRedactedValue
		if parts[3][0] == '"' || parts[3][0] == '\'' {
			replacement = parts[3][:1] + replacement + parts[3][:1]
		}
		return parts[1] + parts[2] + replacement
	})
	// 自由文本仅对手机号作格式兜底；学号和普通长字符串按字段名识别。
	return feedbackDigitPattern.ReplaceAllStringFunc(value, func(number string) string {
		if len(number) == 11 && number[0] == '1' && number[1] >= '3' && number[1] <= '9' {
			return feedbackLogRedactedValue
		}
		return number
	})
}
