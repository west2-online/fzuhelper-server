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
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common/commonservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

type CourseService struct {
	ctx          context.Context
	db           *db.Database
	sf           *utils.Snowflake
	cache        *cache.Cache
	taskQueue    taskqueue.TaskQueue
	commonClient commonservice.Client
	userClient   userservice.Client
}

func NewCourseService(ctx context.Context, clientset *base.ClientSet, taskQueue taskqueue.TaskQueue) *CourseService {
	return &CourseService{
		ctx:          ctx,
		db:           clientset.DBClient,
		sf:           clientset.SFClient,
		cache:        clientset.CacheClient,
		taskQueue:    taskQueue,
		commonClient: clientset.CommonClient,
		userClient:   clientset.UserClient,
	}
}
