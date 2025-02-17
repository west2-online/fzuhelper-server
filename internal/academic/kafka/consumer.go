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

package kafka

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/kafka"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

type AcademicConsumer struct {
	kafkaInstance *kafka.Kafka
	cache         *cache.Cache
}

func InitAcademicConsumer(cache *cache.Cache, instance *kafka.Kafka) *AcademicConsumer {
	return &AcademicConsumer{
		kafkaInstance: instance,
		cache:         cache,
	}
}

func (c *AcademicConsumer) Close() {
	c.kafkaInstance.Close()
}

func (c *AcademicConsumer) ConsumeMessage(topic string, consumerNum int, groupID string) {
	consumerCh := c.kafkaInstance.Consume(context.Background(), topic, consumerNum, groupID)
	for i := 0; i < consumerNum; i++ {
		go c.handleKafkaMessages(consumerCh)
	}
}

func (c *AcademicConsumer) handleKafkaMessages(consumerCh <-chan *kafka.Message) {
	for kafkaMsg := range consumerCh {
		key := utils.BytesToInt64(kafkaMsg.K)

		switch key {
		case constants.AcademicSetScoresCacheEventKey:
			var payload ScoresCacheMessage
			err := sonic.Unmarshal(kafkaMsg.V, &payload)
			if err != nil {
				logger.Errorf("Kafka consumer: Failed to unmarshal ScoresCache payload: %v", err)
				continue
			}
			c.cache.Academic.SetScoresCache(context.Background(), payload.Key, payload.Scores)

		default:
			logger.Warnf("Kafka consumer: Unknown message key: %v", key)
		}
	}
}
