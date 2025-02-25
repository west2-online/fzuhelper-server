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
	"strings"
	"time"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/umeng"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

// SetScoresCacheTask 定义
type SetScoresCacheTask struct {
	key     string
	scores  []*jwch.Mark
	cache   *cache.Cache
	context context.Context
}

func NewSetScoresCacheTask(key string, scores []*jwch.Mark, cache *cache.Cache, context context.Context) *SetScoresCacheTask {
	return &SetScoresCacheTask{
		key:     key,
		scores:  scores,
		cache:   cache,
		context: context,
	}
}

func (t *SetScoresCacheTask) Execute() error {
	return t.cache.Academic.SetScoresCache(t.context, t.key, t.scores)
}

type PutScoresToDatabaseTask struct {
	ctx    context.Context
	db     *db.Database
	id     string
	scores []*jwch.Mark
}

func NewPutScoresToDatabaseTask(ctx context.Context, db *db.Database, id string, scores []*jwch.Mark,
) *PutScoresToDatabaseTask {
	return &PutScoresToDatabaseTask{
		ctx:    ctx,
		db:     db,
		id:     id,
		scores: scores,
	}
}

func (t *PutScoresToDatabaseTask) Execute() error {
	stuID := utils.ParseJwchStuId(t.id)
	// 获取旧成绩 hash
	oldSha256, err := t.db.Academic.GetScoreSha256ByStuId(t.ctx, stuID)
	if err != nil {
		return err
	}
	// 生成新成绩json 和 hash
	json, err := utils.JSONEncode(t.scores)
	if err != nil {
		return err
	}

	newSha256 := utils.SHA256(json)

	// 成绩信息不存在，直接存数据库后返回
	if oldSha256 == "" {
		_, err = t.db.Academic.CreateUserScore(t.ctx, &model.Score{
			StuID:            stuID,
			ScoresInfo:       json,
			ScoresInfoSHA256: newSha256,
		})
		if err != nil {
			return err
		}
		return nil
	} else if oldSha256 != newSha256 {
		// 处理推送逻辑
		err = t.handleScoreChange(stuID)
		if err != nil {
			return err
		}
		// 更新成绩信息
		err = t.db.Academic.UpdateUserScores(t.ctx, &model.Score{
			StuID:            stuID,
			ScoresInfo:       json,
			ScoresInfoSHA256: newSha256,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *PutScoresToDatabaseTask) handleScoreChange(stuID string) (err error) {
	// 成绩信息存在并且和db中的不同，比较成绩是否更新
	var old *model.Score
	old, err = t.db.Academic.GetScoreByStuId(t.ctx, stuID)
	if err != nil {
		return err
	}
	var oldScores []*jwch.Mark
	err = sonic.Unmarshal([]byte(old.ScoresInfo), &oldScores)
	if err != nil {
		return err
	}
	// 反转 oldScores 和 t.scores，方便判断是新课程还是成绩更新
	reverseScores := func(scores []*jwch.Mark) {
		for i := 0; i < len(scores)/2; i++ {
			scores[i], scores[len(scores)-1-i] = scores[len(scores)-1-i], scores[i]
		}
	}
	// 反转两个切片
	reverseScores(t.scores)
	reverseScores(oldScores)
	// 循环比较新旧成绩数据
	for i := range t.scores {
		// 如果下标更大，表示已经遍历到新增课程，结束循环
		if i >= len(oldScores) {
			break
		}

		if !(t.scores[i].Score == oldScores[i].Score) {
			// 尝试获取课程信息
			courseHash := utils.GenerateCourseHash(t.scores[i].Name, t.scores[i].Semester, t.scores[i].Teacher,
				t.scores[i].ElectiveType)
			existingCourse, err := t.db.Academic.GetCourseByHash(t.ctx, courseHash)
			if err != nil {
				return err
			}
			// 课程信息不存在，说明还未发过通知
			if existingCourse == nil {
				// md5 作为tag
				tag := utils.MD5(strings.Join([]string{
					t.scores[i].Name, t.scores[i].Semester, t.scores[i].Teacher,
					t.scores[i].ElectiveType,
				}, "|"))
				err = t.sendNotifications(t.scores[i].Name, tag)
				if err != nil {
					return err
				}
				// 写入课程信息，代表发送过通知
				_, err = t.db.Academic.CreateCourseOffering(t.ctx, &model.CourseOffering{
					Name:         t.scores[i].Name,
					Term:         t.scores[i].Semester,
					Teacher:      t.scores[i].Teacher,
					ElectiveType: t.scores[i].ElectiveType,
					CourseHash:   courseHash,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *PutScoresToDatabaseTask) sendNotifications(courseName, tag string) (err error) {
	err = umeng.SendAndroidGroupcast(config.Umeng.Android.AppKey, config.Umeng.Android.AppMasterSecret,
		"", fmt.Sprintf("%v成绩更新啦", courseName), "",
		tag)
	if err != nil {
		logger.Errorf("task queue: failed to send notice to Android: %v", err)
	}
	err = umeng.SendIOSGroupcast(config.Umeng.IOS.AppKey, config.Umeng.IOS.AppMasterSecret,
		fmt.Sprintf("%v成绩更新啦", courseName), "", "",
		tag)
	if err != nil {
		logger.Errorf("task queue: failed to send notice to IOS: %v", err)
	}

	logger.Infof("task queue: send notice to app, tag:%v", tag)
	// 停止5秒防止 umeng 限流
	time.Sleep(constants.UmengRateLimitDelay)
	return nil
}
