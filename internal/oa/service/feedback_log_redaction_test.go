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
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// 仅供测试与旧缓冲方案基准使用；生产上传直接将脱敏结果写入 gzip。
func redactedFeedbackLogBytes(data []byte) ([]byte, error) {
	var dst bytes.Buffer
	if err := redactFeedbackLog(data, &dst); err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}

type feedbackFailWriter struct {
	err error
}

func (w feedbackFailWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

func TestFeedbackLogRedactionWriterError(t *testing.T) {
	want := errors.New("writer failed")
	err := redactFeedbackLog([]byte(`{"ts":1788912000}`), feedbackFailWriter{err: want})
	require.ErrorIs(t, err, want)
}

func TestFeedbackLogRedactionRegression(t *testing.T) {
	tests := []struct{ name, input, want string }{
		{"credentials in message", `{"message":"password=hunter2&token=short-token"}`, `{"message":"password=***&token=***"}`},
		{"quoted password", `{"message":"password='two words', token=short"}`, `{"message":"password='***', token=***"}`},
		{
			"headers",
			`{"headers":{"Set-Cookie":["sid=one; Secure","other=two"],"Proxy-Authorization":"Basic value","ACCESS_TOKEN":"short"}}`,
			`{"headers":{"Set-Cookie":"***","Proxy-Authorization":"***","ACCESS_TOKEN":"***"}}`,
		},
		{
			"raw headers",
			`{"message":"Cookie: sid=one; other=two\nSet-Cookie: sid=three; Secure\nstatus: 200"}`,
			`{"message":"Cookie: ***\nSet-Cookie: ***\nstatus: 200"}`,
		},
		{
			"numeric diagnostic fields",
			`{"timestamp":1788912000,"duration":123456789,"count":13800138000,"ratio":1.25,"ts":1788912000123}`,
			`{"timestamp":1788912000,"duration":123456789,"count":13800138000,"ratio":1.25,"ts":1788912000123}`,
		},
		{
			"string diagnostics",
			`{"timestamp":"1788912000","trace_id":"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",` +
				`"token_count":12,"session_id":"session"}`,
			`{"timestamp":"1788912000","trace_id":"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",` +
				`"token_count":12,"session_id":"session"}`,
		},
		{
			"identity by key",
			`{"nested":[{"stu_id":102301000,"student-id":"052106112","phone":13800138000}]}`,
			`{"nested":[{"stu_id":"***","student-id":"***","phone":"***"}]}`,
		},
		{"phone fallback", `{"message":"call 13800138000, ts 1788912000"}`, `{"message":"call ***, ts 1788912000"}`},
		{
			"URL encoded keys and repeated values",
			`{"url":"https://example.com/path?%74oken=one&token=two&page=2#access_token=three"}`,
			`{"url":"https://example.com/path?page=2&token=***"}`,
		},
		{"URL userinfo", `{"url":"https://user:password@example.com/path"}`, `{"url":"https://%2A%2A%2A@example.com/path"}`},
		{"malformed query", `{"url":"https://example.com/path?token=secret%ZZ"}`, `{"url":"https://example.com/path"}`},
		{"relative encoded query", `{"url":"/login?%74oken=short&page=2"}`, `{"url":"/login?page=2&token=***"}`},
		{"bare query", `{"message":"/login?password=short&student_id=052106112&page=2"}`, `{"message":"/login?page=2&password=***&student_id=***"}`},
		{"auth schemes", `{"message":"Bearer short-token, Basic dXNlcjpwYXNz"}`, `{"message":"***, ***"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := redactedFeedbackLogBytes([]byte(tt.input))
			require.NoError(t, err)
			require.JSONEq(t, tt.want, string(got))
			again, err := redactedFeedbackLogBytes(got)
			require.NoError(t, err)
			require.Equal(t, got, again, "redaction must be idempotent")
		})
	}
}

func TestFeedbackLogRedactionJSONLines(t *testing.T) {
	for _, input := range [][]byte{
		nil, []byte(" \r\n"),
		{0xff},
		[]byte(`{"token":"secret"} {}`),
		[]byte("{}\ninvalid"), []byte(`{"password":"do-not-echo",}`),
	} {
		out, err := redactedFeedbackLogBytes(input)
		require.Error(t, err)
		require.Nil(t, out, "never return partially redacted content")
		require.NotContains(t, err.Error(), "do-not-echo")
	}
	input := "\r\n{\"phone\":13800138000}\r\n\n{\"ts\":1788912000}"
	out, err := redactedFeedbackLogBytes([]byte(input))
	require.NoError(t, err)
	require.Equal(t, "{\"phone\":\"***\"}\n{\"ts\":1788912000}\n", string(out))
	// 不引入 bufio.Scanner 默认 64 KB 单行限制。
	out, err = redactedFeedbackLogBytes([]byte(`{"message":"` + strings.Repeat("x", 70*1024) + `"}`))
	require.NoError(t, err)
	require.True(t, json.Valid(bytes.TrimSpace(out)))
}

func BenchmarkFeedbackLogRedaction(b *testing.B) {
	data := bytes.Repeat([]byte(`{"ts":1788912000,"message":"password=short","url":"https://example.com/path?token=short&page=2"}`+"\n"), 1000)
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	for b.Loop() {
		if _, err := redactedFeedbackLogBytes(data); err != nil {
			b.Fatal(err)
		}
	}
}
