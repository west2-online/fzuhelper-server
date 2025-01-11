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
	"fmt"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/jwch"
)

func (s *CommonService) GetNotice(pageNum int) (list *[]model.Notice, total int, err error) {
	list, err = s.db.Notice.GetNoticeByPage(s.ctx, pageNum)
	if err != nil {
		return nil, 0, fmt.Errorf("CommonService.GetNotice get notice from database:%w", err)
	}
	// 爬取总页数
	_, total, err = jwch.NewStudent().GetNoticeInfo(&jwch.NoticeInfoReq{PageNum: 1})
	if err != nil {
		return nil, 0, errno.Errorf(errno.BizJwchCookieExceptionCode, "dal.GetNoticeByPage error: %s", err)
	}
	return list, total, nil
}
