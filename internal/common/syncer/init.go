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

package syncer

import (
	"context"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

// Init 当包被导入时，默认执行一次全局爬取 需要显式调用
// TODO: 失败后的重试机制
func (ns *NoticeSyncer) initNoticeSyncer() {
	stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
	_, totalPage, err := stu.GetNoticeInfo(&jwch.NoticeInfoReq{PageNum: 1})
	if err != nil {
		logger.Errorf("syncer init: failed to get notice info: %v", err)
	}
	// 初始化数据库
	for i := 1; i <= totalPage; i++ {
		content, _, err := stu.GetNoticeInfo(&jwch.NoticeInfoReq{PageNum: i})
		if err != nil {
			logger.Errorf("syncer init: failed to get notice info in page %d: %v", i, err)
		}
		for _, row := range content {
			ctx := context.Background()
			info := &model.Notice{
				Title:       row.Title,
				PublishedAt: row.Date,
				URL:         row.URL,
			}
			err = ns.db.Notice.CreateNotice(ctx, info)
			if err != nil {
				logger.Errorf("syncer init: failed to create notice in page %d: %v", i, err)
			}
		}
	}
	logger.Infof("syncer init: notice syncer init success")
}
