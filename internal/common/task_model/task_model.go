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

package task_model

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/common/pack"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/github"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/umeng"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/jwch"
)

type NoticeSyncTask struct {
	db *db.Database
}

func NewNoticeSyncTask(db *db.Database) *NoticeSyncTask {
	return &NoticeSyncTask{
		db: db,
	}
}

func (t *NoticeSyncTask) Execute() error {
	// 默认爬取第一页的内容（教务处不太可能一次性更新出一页的数据），然后和数据库做 diff 操作
	content, _, err := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password).GetNoticeInfo(&jwch.NoticeInfoReq{PageNum: 1})
	if err != nil {
		logger.Errorf("notice sycn task: failed to get notice info: %v", err)
		return fmt.Errorf("failed to get notice info: %w", err)
	}

	for _, row := range content {
		// 判断是否已存在
		ctx := context.Background()
		ok, err := t.db.Notice.IsURLExists(ctx, row.URL)
		if err != nil {
			return fmt.Errorf("notice sycn task: failed to check url exists: %w", err)
		}

		// 数据库已存在，无需处理
		if ok {
			continue
		}

		info := &model.Notice{
			Title:       row.Title,
			URL:         row.URL,
			PublishedAt: row.Date,
		}

		if err = t.db.Notice.CreateNotice(ctx, info); err != nil {
			return fmt.Errorf("notice sycn task: failed to create notice: %w", err)
		}

		// 进行消息推送
		err = umeng.SendAndroidGroupcastWithUrl(config.Umeng.Android.AppKey, config.Umeng.Android.AppMasterSecret,
			"", "教务处通知", info.Title, constants.UmengJwchNoticeTag, info.URL, map[string]string{
				"channel_activity":          "com.west2online.umeng.MfrMessageActivity",
				"xiaomi_channel_id":         config.Vendors.Xiaomi.JwchNotice,
				"huawei_channel_importance": "NORMAL",
			})
		if err != nil {
			logger.Errorf("notice sycn task: failed to send notice to Android: %v", err)
		}

		// ios 无法跳转url
		err = umeng.SendIOSGroupcast(config.Umeng.IOS.AppKey, config.Umeng.IOS.AppMasterSecret,
			"教务处通知", "", info.Title, constants.UmengJwchNoticeTag)
		if err != nil {
			logger.Errorf("notice sycn task: failed to send notice to IOS: %v", err)
		}

		logger.Infof("notice sycn task: send notice successfully")
		break
	}
	return nil
}

func (t *NoticeSyncTask) GetScheduleTime() time.Duration {
	return constants.NoticeUpdateTime
}

type ContributorInfoSyncTask struct {
	cache *cache.Cache
}

func NewContributorInfoSyncTask(c *cache.Cache) *ContributorInfoSyncTask {
	return &ContributorInfoSyncTask{
		cache: c,
	}
}

func (t *ContributorInfoSyncTask) Execute() error {
	urls := []string{
		constants.ContributorFzuhelperApp,
		constants.ContributorFzuhelperServer,
		constants.ContributorJwch,
		constants.ContributorYJSY,
	}
	contributorKeys := []string{
		constants.ContributorFzuhelperAppKey,
		constants.ContributorFzuhelperServerKey,
		constants.ContributorJwchKey,
		constants.ContributorYJSYKey,
	}
	for i, url := range urls {
		rawContributors, err := github.FetchContributorsFromURL(url)
		if err != nil {
			return fmt.Errorf("contributor info sync: failed to fetch contributors from %s: %w", url, err)
		}
		contributors := pack.BuildContributors(rawContributors)
		for i, contributor := range contributors {
			newAvatarUrl, err := uploadAvatar(contributor.AvatarUrl, contributor.Name)
			if err != nil {
				return fmt.Errorf("contributor info sync: failed to upload avatar for contributor %s: %w", contributor.Name, err)
			}
			// 替换头像 url
			contributors[i].AvatarUrl = newAvatarUrl
		}

		if err := t.cache.Common.SetContributorInfo(context.Background(), contributorKeys[i], contributors); err != nil {
			return fmt.Errorf("contributor info sync: failed to write updated contributors to Redis: %w", err)
		}
	}

	return nil
}

const (
	baseUrl    = "https://avatars.githubusercontent.com/u/"
	uploadBase = "http://v0.api.upyun.com/fzuhelper-filedown"
	readBase   = "https://download.w2fzu.com"
)

func uploadAvatar(avatarUrl string, name string) (string, error) {
	if strings.HasPrefix(avatarUrl, baseUrl) {
		// 1.将原始 URL 替换成反代 URL
		parsedUrl, err := url.Parse(avatarUrl)
		if err != nil {
			return "", err
		}
		// parsedUrl.Path[3:]会去掉 `/u/`
		newAvatarUrl := fmt.Sprintf(constants.AvatarProxy, parsedUrl.Path[3:])

		// 2.下载图片并上传又拍云
		resp, err := http.Get(newAvatarUrl)
		if err != nil {
			return "", fmt.Errorf("failed to download avatar from %s: %w", avatarUrl, err)
		}
		imgData, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read avatar image: %w", err)
		}
		// 生成上传用Url
		newAvatarUrl = upyun.GenerateContributorAvatarUrl(name)
		err = upyun.URlUploadFile(imgData, newAvatarUrl)
		if err != nil {
			return "", fmt.Errorf("failed to upload avatar to image host: %w", err)
		}
		_ = resp.Body.Close()

		// 3.最终换成加速域名
		return strings.Replace(newAvatarUrl, uploadBase, readBase, 1), nil
	}

	return "", nil
}

func (t *ContributorInfoSyncTask) GetScheduleTime() time.Duration {
	return constants.ContributorInfoUpdateTime
}
