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
	"net"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func Init() {
	topic := "Info"

	conn := GetConn()
	defer closeFn(conn)
	if err := conn.CreateTopics( // 创建topic是幂等的，如果topic已经存在会返回nil
		kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	); err != nil {
		panic(err)
	}

	w := NewWriter()
	defer closeFn(w)
	if err := w.WriteMessages(
		context.TODO(),
		kafka.Message{
			Topic: topic,
			Key:   []byte("Info Key"),
			Value: []byte("Hello world"),
			Time:  time.Now(),
		}); err != nil {
		panic(err)
	}

	r := NewReader(topic)
	defer closeFn(r)
	msg, err := r.ReadMessage(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message`s data, key:%s, value: %s, time:%v\n", msg.Key, msg.Value, msg.Time)
}

// GetConn conn不能保证并发安全,仅可作为单线程的长连接使用。
func GetConn() *kafka.Conn {
	conn, err := kafka.Dial("tcp", "127.0.0.1:9093")
	if err != nil {
		panic(err)
	}
	return conn
}

// NewReader 创建一个reader示例，reader是并发安全的
func NewReader(topic string) *kafka.Reader {
	cfg := kafka.ReaderConfig{
		Brokers:     []string{"127.0.0.1:9093"},
		Topic:       topic,
		MinBytes:    1,
		MaxBytes:    1 * 1024 * 1024,
		MaxAttempts: 3,
	}
	return kafka.NewReader(cfg)
}

// NewWriter 创建一个writer示例，writer是并发安全的。 errLogger可以传入带有es hook的logger
func NewWriter() *kafka.Writer {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9093")
	if err != nil {
		panic(err)
	}

	return &kafka.Writer{
		Addr:                   addr,
		Balancer:               &kafka.RoundRobin{},
		MaxAttempts:            3,
		RequiredAcks:           kafka.RequireOne,
		Async:                  true,
		AllowAutoTopicCreation: false,
	}
}

func closeFn(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(err)
	}
}
