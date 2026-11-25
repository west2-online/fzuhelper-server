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
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// 合成接近 2 MB 的客户端日志，包含不同时间、追踪标识、网络摘要、事件及敏感字段，不使用真实用户数据。
func feedbackLogBenchmarkData() []byte {
	var data bytes.Buffer
	for i := 0; ; i++ {
		trace := sha256.Sum256(fmt.Append(nil, i))
		line := fmt.Sprintf(
			"{\"ts\":%d,\"level\":\"info\",\"event\":\"request\",\"trace_id\":\"%x\",\"duration\":%d,"+
				"\"status\":%d,\"url\":\"https://example.com/api/course?page=%d&token=secret\","+
				"\"headers\":{\"Authorization\":\"Bearer secret\"},\"message\":\"request completed\"}\n",
			int64(1788912000000)+int64(i), trace, i%1500, 200+i%4, i%20,
		)
		if data.Len()+len(line) > constants.FeedbackLogMaxSize {
			return data.Bytes()
		}
		data.WriteString(line)
	}
}

func benchmarkBufferedFeedbackLog(data []byte) ([]byte, error) {
	redacted, err := redactedFeedbackLogBytes(data)
	if err != nil {
		return nil, err
	}
	var dst bytes.Buffer
	writer := gzip.NewWriter(&dst)
	if _, err = writer.Write(redacted); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}

func BenchmarkFeedbackLogCompression(b *testing.B) {
	data := feedbackLogBenchmarkData()
	for _, tc := range []struct {
		name     string
		buffered bool
		level    int
	}{
		{name: "buffered-default", buffered: true},
		{name: "stream-default", level: gzip.DefaultCompression},
		{name: "stream-speed", level: gzip.BestSpeed},
	} {
		compress := func() ([]byte, error) {
			if tc.buffered {
				return benchmarkBufferedFeedbackLog(data)
			}
			var dst bytes.Buffer
			err := compressFeedbackLog(data, &dst, tc.level)
			return dst.Bytes(), err
		}
		for _, parallel := range []bool{false, true} {
			mode := "serial"
			if parallel {
				mode = "parallel"
			}
			b.Run(tc.name+"/"+mode, func(b *testing.B) {
				sample, err := compress()
				if err != nil {
					b.Fatal(err)
				}
				b.ReportAllocs()
				b.SetBytes(int64(len(data)))
				b.ResetTimer()
				if parallel {
					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							if _, err := compress(); err != nil {
								b.Error(err)
								return
							}
						}
					})
				} else {
					for b.Loop() {
						if _, err := compress(); err != nil {
							b.Fatal(err)
						}
					}
				}
				b.ReportMetric(float64(len(sample)), "compressed-B/op")
				b.ReportMetric(float64(len(sample))/float64(len(data)), "compressed/input")
			})
		}
	}
}
