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
	"strings"
	"time"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/config"
	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/umeng"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (s *AcademicService) GetScores(loginData *loginmodel.LoginData) ([]*jwch.Mark, error) {
	key := fmt.Sprintf("scores:%s", context.ExtractIDFromLoginData(loginData))
	loginData, err := context.GetLoginData(s.ctx)
	stuId := context.ExtractIDFromLoginData(loginData)
	if err != nil {
		return nil, errno.NewErrNo(errno.AuthErrorCode, fmt.Sprintf("service.GetScores: Get login data fail %v", err))
	}

	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		scores, err := s.cache.Academic.GetScoresCache(s.ctx, key)
		if err != nil {
			return nil, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("service.GetScores: Get scores info from redis error %v", err))
		}
		return scores, nil
	} else {
		stu := jwch.NewStudent().WithLoginData(loginData.Id, utils.ParseCookies(loginData.Cookies))
		scores, err := stu.GetMarks()
		if err = base.HandleJwchError(err); err != nil {
			return nil, fmt.Errorf("service.GetScores: Get scores info fail %w", err)
		}
		s.taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
			return cache.SetSliceCache(s.cache, s.ctx, key, scores, constants.AcademicScoresExpire, "Academic.SetScores")
		}})
		s.taskQueue.Add(stuId, taskqueue.QueueTask{Execute: func() error {
			return s.checkScoreChange(stuId, scores)
		}})
		return scores, nil
	}
}

func (s *AcademicService) GetScoresYjsy(loginData *loginmodel.LoginData) ([]*yjsy.Mark, error) {
	key := fmt.Sprintf("scores:%s", loginData.Id[len(loginData.Id)-constants.StudentIDLength:])
	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		scores, err := s.cache.Academic.GetScoresCacheYjsy(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetScoresYjsy: Get scores info from redis error %w", err)
		}
		return scores, nil
	} else {
		stu := yjsy.NewStudent().WithLoginData(utils.ParseCookies(loginData.Cookies))
		scores, err := stu.GetMarks()
		if err = base.HandleYjsyError(err); err != nil {
			return nil, fmt.Errorf("service.GetScoresYjsy: Get scores info fail %w", err)
		}
		s.taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
			return cache.SetSliceCache(s.cache, s.ctx, key, scores, constants.AcademicScoresExpire, "Academic.SetScoresYjsy")
		}})
		// 研究生暂时不做成绩推送，也就不做持久化
		return scores, nil
	}
}

func (s *AcademicService) checkScoreChange(stuId string, scores []*jwch.Mark) error {
	// 获取旧成绩 hash
	oldSha256, err := s.db.Academic.GetScoreSha256ByStuId(s.ctx, stuId)
	if err != nil {
		return err
	}
	// 生成新成绩json 和 hash
	json, err := utils.JSONEncode(scores)
	if err != nil {
		return err
	}

	newSha256 := utils.SHA256(json)

	// 成绩信息不存在，直接存数据库后返回
	if oldSha256 == "" {
		_, err = s.db.Academic.CreateUserScore(s.ctx, &model.Score{
			StuID:            stuId,
			ScoresInfo:       json,
			ScoresInfoSHA256: newSha256,
		})
		if err != nil {
			return err
		}
		return nil
	} else if oldSha256 != newSha256 {
		// 处理推送逻辑
		err = s.handleScoreChange(stuId, scores)
		if err != nil {
			return err
		}
		// 更新成绩信息
		err = s.db.Academic.UpdateUserScores(s.ctx, &model.Score{
			StuID:            stuId,
			ScoresInfo:       json,
			ScoresInfoSHA256: newSha256,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AcademicService) handleScoreChange(stuID string, scores []*jwch.Mark) (err error) {
	// 成绩信息存在并且和db中的不同，比较成绩是否更新
	var old *model.Score
	old, err = s.db.Academic.GetScoreByStuId(s.ctx, stuID)
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
	reverseScores(scores)
	reverseScores(oldScores)
	// 循环比较新旧成绩数据
	for i := range scores {
		// 如果下标更大，表示已经遍历到新增课程，结束循环
		if i >= len(oldScores) {
			break
		}

		if !(scores[i].Score == oldScores[i].Score) {
			// 尝试获取课程信息
			courseHash := utils.GenerateCourseHash(scores[i].Name, scores[i].Semester, scores[i].Teacher,
				scores[i].ElectiveType)
			existingCourse, err := s.db.Academic.GetCourseByHash(s.ctx, courseHash)
			if err != nil {
				return err
			}
			// 课程信息不存在，说明还未发过通知
			if existingCourse == nil {
				// md5 作为tag
				tag := utils.MD5(strings.Join([]string{
					scores[i].Name, scores[i].Semester, scores[i].Teacher,
					scores[i].ElectiveType,
				}, "|"))
				err = s.sendNotifications(scores[i].Name, tag)
				if err != nil {
					return err
				}
				// 写入课程信息，代表发送过通知
				_, err = s.db.Academic.CreateCourseOffering(s.ctx, &model.CourseOffering{
					Name:         scores[i].Name,
					Term:         scores[i].Semester,
					Teacher:      scores[i].Teacher,
					ElectiveType: scores[i].ElectiveType,
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

func (s *AcademicService) sendNotifications(courseName, tag string) (err error) {
	err = umeng.SendAndroidGroupcastWithGoApp(config.Umeng.Android.AppKey, config.Umeng.Android.AppMasterSecret,
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
