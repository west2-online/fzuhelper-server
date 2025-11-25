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

package base

import (
	cli "github.com/cloudwego/hertz/pkg/app/client"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/oss"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// WithRedisClient will create redis object
func WithRedisClient(dbName int) Option {
	return func(clientSet *ClientSet) {
		redisClient, err := client.NewRedisClient(dbName)
		if err != nil {
			logger.Fatalf("init cache failed, err: %v", err)
		}
		clientSet.CacheClient = cache.NewCache(redisClient)
		clientSet.cleanups = append(clientSet.cleanups, func() {
			err = redisClient.Close()
			if err != nil {
				logger.Errorf("close cache failed, err: %v", err)
			}
		})

		logger.Infof("Cache Redis Connect Success")
	}
}

// WithDBClient will create database object
func WithDBClient() Option {
	return func(clientSet *ClientSet) {
		DB, err := client.InitMySQL()
		if err != nil {
			logger.Fatalf("init database failed, err: %v", err)
		}

		// TODO: currently our service only deploy on one or some servers, we do not need specific datacenterID
		sf, err := utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID)
		if err != nil {
			logger.Fatalf("init Snowflake object error: %v", err.Error())
		}
		clientSet.SFClient = sf
		clientSet.DBClient = db.NewDatabase(DB, sf)
		// gorm maintains a connection pool.
		// After initialization, all connections are managed by gorm and do not need to be closed manually.

		logger.Infof("Database MySQL Connect Success")
	}
}

func WithElasticSearch() Option {
	return func(clientSet *ClientSet) {
		es, err := client.NewEsClient()
		if err != nil {
			logger.Fatalf("init elastic search client error: %v", err.Error())
		}
		clientSet.ESClient = es
		logger.Infof("ElasticSearch Connect Success")
	}
}

func WithHzClient() Option {
	return func(clientSet *ClientSet) {
		hz, err := cli.NewClient()
		if err != nil {
			logger.Fatalf("init Hertz client error: %v", err)
		}
		clientSet.HzClient = hz
		logger.Infof("Hertz Client Create Success")
	}
}

func WithCommonRPCClient() Option {
	return func(clientSet *ClientSet) {
		client, err := client.InitCommonRPC()
		if err != nil {
			logger.Fatalf("init common rpc client error: %v", err)
		}
		clientSet.CommonClient = *client
		logger.Infof("Common RPC Client Create Success")
	}
}

func WithUserRPCClient() Option {
	return func(clientSet *ClientSet) {
		client, err := client.InitUserRPC()
		if err != nil {
			logger.Fatalf("init user rpc client error: %v", err)
		}
		clientSet.UserClient = *client
		logger.Infof("User RPC Client Create Success")
	}
}

func WithOssSet(provider string) Option {
	return func(clientSet *ClientSet) {
		ossSet := &oss.OSSSet{
			Provider: provider,
		}
		switch ossSet.Provider {
		case oss.UpYunProvider:
			ossSet.Upyun = oss.NewUpYunConfig()
		default:
			logger.Fatalf("unknown ossSet.Provider: %v", ossSet.Provider)
		}
		clientSet.OssSet = ossSet
		logger.Infof("OSS Client Create Success")
	}
}
