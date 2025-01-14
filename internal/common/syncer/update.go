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
	"fmt"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

// update 定期爬取教务处教学通知，并进行 diff 操作
func (ns *NoticeSyncer) update() (err error) {
	// 默认爬取第一页的内容（教务处不太可能一次性更新出一页的数据），然后和数据库做 diff 操作
	content, _, err := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password).GetNoticeInfo(&jwch.NoticeInfoReq{PageNum: 1})
	if err != nil {
		logger.Errorf("syncer update: failed to get notice info: %v", err)
	}
	for _, row := range content {
		// 判断是否已存在
		ctx := context.Background()
		ok, err := ns.db.Notice.IsURLExists(ctx, row.URL)
		if err != nil {
			return fmt.Errorf("syncer update: failed to check url exists: %w", err)
		}
		// 数据库已存在，无需处理
		if ok {
			continue
		}
		// TODO: SDK 进行消息推送等操作
		info := &model.Notice{
			Title:       row.Title,
			URL:         row.URL,
			PublishedAt: row.Date,
		}
		if err = ns.db.Notice.CreateNotice(ctx, info); err != nil {
			return fmt.Errorf("syncer update: failed to create notice: %w", err)
		}
	}
	return nil
}
