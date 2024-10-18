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

package mq

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/mq"

	kafukago "github.com/segmentio/kafka-go"
)

func Init() {
	topic := "Info"

	// 获取链接
	conn, err := mq.GetConn()
	if err != nil {
		panic(err)
	}
	defer closeFn(conn)

	if err = conn.CreateTopics( // 创建topic是幂等的，如果topic已经存在会返回nil
		kafukago.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	); err != nil {
		panic(err)
	}

	// 获取Writer
	w, err := mq.GetNewWriter()
	if err != nil {
		panic(err)
	}
	defer closeFn(w)

	// 写入信息
	if err = w.WriteMessages(
		context.TODO(),
		kafukago.Message{
			Topic: topic,
			Key:   []byte("Info Key"),
			Value: []byte("Hello world"),
			Time:  time.Now(),
		}); err != nil {
		panic(err)
	}

	// 获取Reader
	r := mq.GetNewReader(topic)
	defer closeFn(r)

	// 读取信息
	msg, err := r.ReadMessage(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message`s data, key:%s, value: %s, time:%v\n", msg.Key, msg.Value, msg.Time)
}

func closeFn(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(err)
	}
}
