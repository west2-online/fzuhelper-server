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

package eszap

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	elastic "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ElasticHook is a zap hook for ElasticSearch
type ElasticHook struct {
	client *elastic.Client // es 的客户端
	host   string          // es 的 host
	index  string          // 获取索引的名字
	ctx    context.Context
	cancel func()

	enc zapcore.Encoder
	lvl zap.AtomicLevel
	ws  zapcore.WriteSyncer
}

// IndexNameFunc get index name
type IndexNameFunc func() string

func NewElasticHook(client *elastic.Client, host string, index string, level zapcore.Level) *ElasticHook {
	hook := defaultHookConfig()

	hook.client = client
	hook.host = host
	hook.index = index
	hook.lvl = zap.NewAtomicLevelAt(level)

	return hook
}

// 发送到 es 的信息结构
type message struct {
	Host      string
	Timestamp string `json:"@timestamp"`
	Message   string
	Level     string
}

func createMessage(entry *zapcore.Entry, hook *ElasticHook) *message {
	return &message{
		Host:      hook.host,
		Timestamp: entry.Time.UTC().Format(time.RFC3339),
		Message:   entry.Message,
		Level:     strings.ToUpper(entry.Level.String()),
	}
}

func (hook *ElasticHook) Fire(entry zapcore.Entry) error {
	data, err := json.Marshal(createMessage(&entry, hook))
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:   hook.index,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(hook.ctx, hook.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	var r map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	} else {
		// Print the response status and indexed document version.
		log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	}
	return err
}

func (hook *ElasticHook) Enc() zapcore.Encoder {
	return hook.enc
}

func (hook *ElasticHook) Lvl() zap.AtomicLevel {
	return hook.lvl
}

func (hook *ElasticHook) Ws() zapcore.WriteSyncer {
	return hook.ws
}

func defaultHookConfig() *ElasticHook {
	ctx, cancel := context.WithCancel(context.Background())
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	lvl := zap.NewAtomicLevelAt(zap.DebugLevel)
	ws := zapcore.AddSync(os.Stdout)

	return &ElasticHook{
		ctx:    ctx,
		cancel: cancel,
		enc:    enc,
		lvl:    lvl,
		ws:     ws,
	}
}
