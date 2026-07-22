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
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/umeng"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (s *CourseService) GetCourseList(req *course.CourseListRequest, loginData *kitexModel.LoginData) ([]*kitexModel.Course, error) {
	var err error
	stuId := context.ExtractIDFromLoginData(loginData)
	termKey := fmt.Sprintf("terms:%s", stuId)
	courseKey := fmt.Sprintf("course:%s:%s", stuId, req.Term)
	terms := new(jwch.Term)
	// 学期缓存存在
	isRefresh := false
	if req.IsRefresh != nil {
		isRefresh = *req.IsRefresh
	}
	// 不刷新且cache存在
	if !isRefresh && s.cache.IsKeyExist(s.ctx, termKey) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, termKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetCourseList: Get term fail: %w", err)
		}
		terms.Terms = termsList
		// 只有最新的两个学期的课程才会被放入缓存
		if slices.Contains(pack.GetTop2Terms(terms).Terms, req.Term) &&
			s.cache.IsKeyExist(s.ctx, courseKey) {
			courses, err := s.cache.Course.GetCoursesCache(s.ctx, courseKey)
			if err != nil {
				return nil, fmt.Errorf("service.GetCourseList: Get courses fail: %w", err)
			}
			return s.removeDuplicateCourses(pack.BuildCourse(courses)), nil
		}
	}

	stu := jwch.NewStudent().WithLoginData(loginData.GetId(), utils.ParseCookies(loginData.GetCookies()))

	terms, err = stu.GetTerms()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get terms failed: %w", err)
	}

	// validate term
	if !slices.Contains(terms.Terms, req.Term) {
		return nil, errors.New("service.GetCourseList: Invalid term")
	}

	courses, err := stu.GetSemesterCourses(req.Term, terms.ViewState, terms.EventValidation)
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get semester courses failed: %w", err)
	}

	// async put course list to db
	// 数据库存储原始的课表信息（不包含调课信息）
	originalCourses := pack.BuildCourse(courses)
	s.taskQueue.Add(fmt.Sprintf("putCourse:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		if err := s.putCourseToDatabase(stuId, req.Term, originalCourses); err != nil {
			return err
		}
		return s.putExamToDatabase(stuId, req.Term, courses)
	}})

	adjustCourses, err := s.GetAutoAdjustCourseList(req.Term)
	if err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get adjust course failed: %w", err)
	}

	for _, c := range courses {
		adjustRules := getAdjustRules(c.ScheduleRules, adjustCourses)
		c.ScheduleRules = jwch.ApplyAdjustRules(
			jwch.ApplyAdjustRules(c.ScheduleRules, c.AdjustRules),
			adjustRules,
		)
	}

	if slices.Contains(pack.GetTop2Terms(terms).Terms, req.Term) {
		// async put course list to cache
		// 缓存存储调课后的课表信息
		s.taskQueue.Add(courseKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetSliceCache(s.cache, s.ctx, courseKey, courses,
				constants.CourseTermsKeyExpire, "Course.SetCourseCache")
		}})
		s.taskQueue.Add(termKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetValueSliceCache(s.cache, s.ctx, termKey, terms.Terms, constants.CourseTermsKeyExpire, "Course.SetTermsCache")
		}})
	}

	// 学期列表异步存库
	s.taskQueue.Add(fmt.Sprintf("putTerms:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		return s.putTermToDatabase(stuId, pack.BuildTermOnDB(terms.Terms))
	}})

	return s.removeDuplicateCourses(pack.BuildCourse(courses)), nil
}

func (s *CourseService) putCourseToDatabase(stuId string, term string, courses []*kitexModel.Course) error {
	old, err := s.db.Course.GetUserTermCourseSha256ByStuIdAndTerm(s.ctx, stuId, term)
	if err != nil {
		return err
	}

	json, err := utils.JSONEncode(courses)
	if err != nil {
		return err
	}

	newSha256 := utils.SHA256(json)

	if old == nil {
		dbId, err := s.sf.NextVal()
		if err != nil {
			return err
		}

		_, err = s.db.Course.CreateUserTermCourse(s.ctx, &model.UserCourse{
			Id:                dbId,
			StuId:             stuId,
			Term:              term,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return err
		}
	} else if old.TermCoursesSha256 != newSha256 {
		_, err = s.db.Course.UpdateUserTermCourse(s.ctx, &model.UserCourse{
			Id:                old.Id,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CourseService) putExamToDatabase(stuId string, term string, rawCourses []*jwch.Course) error {
	// 考试通知采用“先抢占全局去重记录，再异步发送，成功后提交快照”的顺序。
	// 发送或快照提交失败时释放去重记录，确保后续刷新仍然可以重试。
	exams := buildCourseExamInfo(rawCourses)
	examInfo, err := utils.JSONEncode(exams)
	if err != nil {
		return errno.Errorf(errno.InternalJSONErrorCode,
			"service.putExamToDatabase: encode exam info failed: %v", err)
	}
	examInfoSHA256, err := courseExamInfoHash(exams)
	if err != nil {
		return errno.Errorf(errno.InternalJSONErrorCode,
			"service.putExamToDatabase: hash exam info failed: %v", err)
	}

	old, err := s.db.Course.GetUserTermCourseByStuIdAndTerm(s.ctx, stuId, term)
	if err != nil {
		return err
	}
	if old == nil {
		return nil
	}
	if old.ExamInfoSHA256 == examInfoSHA256 {
		return nil
	}

	var oldExams []CourseExamInfo
	if old.ExamInfo != "" {
		if err = sonic.Unmarshal([]byte(old.ExamInfo), &oldExams); err != nil {
			return errno.Errorf(errno.InternalJSONErrorCode,
				"service.putExamToDatabase: decode exam info failed: %v", err)
		}
	}
	if old.ExamInfoSHA256 == "" {
		// 历史数据没有考试快照时只建立基线，不把已有考试信息当作新增变化通知。
		return s.updateExamSnapshot(old.Id, examInfo, examInfoSHA256)
	}

	changes := buildCourseExamChanges(term, oldExams, exams)
	if len(changes) == 0 {
		// 内容没有实际变化时只更新快照，避免重复进入全局去重和推送流程。
		return s.updateExamSnapshot(old.Id, examInfo, examInfoSHA256)
	}

	claimed := make([]courseExamChange, 0, len(changes))
	for _, change := range changes {
		// CreateExamOffering 依赖 exam_hash 唯一索引原子抢占发送资格。
		// 返回 nil 表示其他用户已经处理过相同变化，本次不再重复推送。
		offering, createErr := s.db.Course.CreateExamOffering(s.ctx, &model.ExamOffering{
			ExamHash: change.ExamHash,
			Tag:      change.Tag,
		})
		if createErr != nil {
			// 本批次后续无法继续抢占时，释放已经抢到的去重记录。
			// 否则这些考试变化虽然没有完成通知流程，却会被全局去重表永久拦截。
			return errors.Join(createErr, s.releaseExamOfferings(claimed))
		}
		if offering != nil {
			claimed = append(claimed, change)
		}
	}

	if len(claimed) == 0 {
		// 所有变化都已被其他任务去重，本次只同步本地快照，不重复发送通知。
		return s.updateExamSnapshot(old.Id, examInfo, examInfoSHA256)
	}
	// 去重记录已经抢占成功，但通知和快照更新放入同一个异步任务。
	// 只有两步都成功，exam_offerings 记录才会继续保留。
	if !umeng.EnqueueAsync(func() error {
		if err := s.sendExamNotifications(claimed); err != nil {
			// 推送失败时释放去重记录；考试快照也不会更新，后续刷新仍能识别到这次变化。
			// 这样可以允许后续课程刷新重新抢占并重试通知。
			return errors.Join(err, s.releaseExamOfferings(claimed))
		}
		if err := s.updateExamSnapshot(old.Id, examInfo, examInfoSHA256); err != nil {
			// 快照更新失败时释放去重记录，避免通知状态未完整落库却被永久去重。
			// 下次刷新会再次发现变化，并重新执行通知流程。
			return errors.Join(err, s.releaseExamOfferings(claimed))
		}
		return nil
	}) {
		// 队列已满时任务没有进入异步流程，当前请求不会再发送这些通知。
		// 因此需要释放已抢到的去重记录，避免这次未发送的变化无法重试。
		return errors.Join(
			errno.NewErrNo(errno.InternalQueueErrorCode, "service.putExamToDatabase: exam notification queue is full"),
			s.releaseExamOfferings(claimed),
		)
	}
	return nil
}

func (s *CourseService) updateExamSnapshot(id int64, examInfo, examInfoSHA256 string) error {
	// 快照更新是本次考试变化处理的提交步骤；成功后下一次刷新不会再次识别同一变化。
	_, err := s.db.Course.UpdateUserTermCourse(s.ctx, &model.UserCourse{
		Id:             id,
		ExamInfo:       examInfo,
		ExamInfoSHA256: examInfoSHA256,
	})
	return err
}

func (s *CourseService) releaseExamOfferings(changes []courseExamChange) error {
	// exam_offerings 既是全局去重记录，也是本次通知流程的占位状态。
	// 仅在发送或快照提交失败时调用，释放后下一次刷新可以重新抢占并重试。
	var releaseErr error
	for _, change := range changes {
		// exam_offerings 记录表示本次变化已经被某个刷新任务抢占。
		// 只有通知发送和考试快照更新都成功后才应保留；失败时删除该记录，
		// 让后续刷新可以重新抢占，避免考试变化被错误地永久去重。
		if err := s.db.Course.DeleteExamOfferingByHash(s.ctx, change.ExamHash); err != nil {
			releaseErr = errors.Join(releaseErr, err)
		}
	}
	return releaseErr
}

func (s *CourseService) sendExamNotifications(changes []courseExamChange) error {
	// 同一批变化同时发送 Android 和 iOS；记录首个错误，但继续处理剩余通知。
	var firstErr error
	for _, change := range changes {
		title := fmt.Sprintf("%v考试信息更新啦", change.Exam.Name)
		description := fmt.Sprintf("考试信息更新%v", change.Tag[:12])
		if err := umeng.SendAndroidGroupcastWithGoApp(
			title, "", "", change.Tag, description, constants.UmengGradeDeeplink,
		); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
		if err := umeng.SendIOSGroupcast(
			title, "", "", change.Tag, description, constants.UmengGradeDeeplink,
		); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (s *CourseService) GetCourseListYjsy(req *course.CourseListRequest, loginData *kitexModel.LoginData) ([]*kitexModel.Course, error) {
	var err error

	stuId := context.ExtractIDFromLoginData(loginData)
	termKey := fmt.Sprintf("terms:%s", stuId)
	courseKey := fmt.Sprintf("course:%s:%s", stuId, req.Term)
	terms := new(yjsy.Term)
	// 学期缓存存在
	isRefresh := false
	if req.IsRefresh != nil {
		isRefresh = *req.IsRefresh
	}
	if !isRefresh && s.cache.IsKeyExist(s.ctx, termKey) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, termKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetCourseListYjsy: Get terms fail: %w", err)
		}
		terms.Terms = termsList

		// 检查是否有该学期的课程缓存
		if slices.Contains(pack.GetTop2TermsYjsy(terms).Terms, req.Term) && s.cache.IsKeyExist(s.ctx, courseKey) {
			courses, err := s.cache.Course.GetCoursesCacheYjsy(s.ctx, courseKey)
			if err != nil {
				return nil, fmt.Errorf("service.GetCourseListYjsy: Get courses fail: %w", err)
			}
			return pack.BuildCourseYjsy(courses), nil
		}
	}

	// 获取学期信息
	stu := yjsy.NewStudent().WithLoginData(utils.ParseCookies(loginData.Cookies))
	terms, err = stu.GetTerms()
	if err = base.HandleYjsyError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseListYjsy: Get terms failed: %w", err)
	}

	// 验证学期是否有效
	if !slices.Contains(terms.Terms, req.Term) {
		return nil, errors.New("service.GetCourseListYjsy: Invalid term")
	}

	// 获取该学期的课程
	courses, err := stu.GetSemesterCourses(req.Term)
	if err = base.HandleYjsyError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseListYjsy: Get semester courses failed: %w", err)
	}

	// 如果是前两个学期，异步缓存课程列表
	if slices.Contains(pack.GetTop2TermsYjsy(terms).Terms, req.Term) {
		s.taskQueue.Add(courseKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetSliceCache(s.cache, s.ctx, courseKey, courses,
				constants.CourseTermsKeyExpire, "Course.SetCourseCache")
		}})
		s.taskQueue.Add(termKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetValueSliceCache(s.cache, s.ctx, termKey, terms.Terms, constants.CourseTermsKeyExpire, "Course.SetTermsCache")
		}})
	}

	// 异步将课程列表存入数据库
	s.taskQueue.Add(fmt.Sprintf("putCourse:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		return s.putCourseToDatabase(stuId, req.Term, pack.BuildCourseYjsy(courses))
	}})
	// 学期列表异步存库
	s.taskQueue.Add(fmt.Sprintf("putTerms:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		return s.putTermToDatabase(stuId, pack.BuildTermOnDB(terms.Terms))
	}})

	return pack.BuildCourseYjsy(courses), nil
}

// removeDuplicateCourses 移除重复课程，只保留第一个出现的。
func (s *CourseService) removeDuplicateCourses(courses []*kitexModel.Course) []*kitexModel.Course {
	seen := make(map[string]struct{})
	var result []*kitexModel.Course

	for _, c := range courses {
		srIDs := make([]string, 0, len(c.ScheduleRules))
		for _, rule := range c.ScheduleRules {
			part := fmt.Sprintf("%d-%d-%d-%d",
				rule.StartClass, rule.EndClass,
				rule.StartWeek, rule.EndWeek)
			srIDs = append(srIDs, part)
		}
		sort.Strings(srIDs)

		// 把“课程名 + 教师 + 排课信息”拼成一个全局唯一的 key
		identifier := fmt.Sprintf("%s-%s-%s", c.Name, c.Teacher, strings.Join(srIDs, "|"))

		// 如果 map 里还没出现过这个标识，那就是新课程
		if _, exists := seen[identifier]; !exists {
			seen[identifier] = struct{}{}
			result = append(result, c)
		}
	}

	return result
}

func (s *CourseService) getSemesterCourses(stuID string, term string, isGraduate bool) (course []*kitexModel.Course, err error) {
	courseKey := fmt.Sprintf("course:%s:%s", stuID, term)
	if s.cache.IsKeyExist(s.ctx, courseKey) {
		courses, err := s.cache.Course.GetCoursesCache(s.ctx, courseKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetSemesterCourses: Get courses fail: %w", err)
		}
		return s.removeDuplicateCourses(pack.BuildCourse(courses)), nil
	}
	// 从数据中获取课程表
	var courses *model.UserCourse
	courses, err = s.db.Course.GetUserTermCourseByStuIdAndTerm(s.ctx, stuID, term)
	if err != nil {
		return nil, fmt.Errorf("service.GetSemesterCourses: Get courses fail: %w", err)
	}
	if courses == nil {
		return nil, errno.NewErrNo(errno.InternalServiceErrorCode, "service.GetSemesterCourses: there is no course in database, please login app and retry")
	}
	// 将数据库中的课程表进行解析转化
	list := make([]*kitexModel.Course, 0)

	if courses.TermCourses != "" {
		if err = sonic.Unmarshal([]byte(courses.TermCourses), &list); err != nil {
			return nil, fmt.Errorf("service.GetSemesterCourses: Unmarshal fail: %w", err)
		}
	}

	// 只处理本科生的调课信息
	if !isGraduate {
		adjustCourses, err := s.GetAutoAdjustCourseList(term)
		if err != nil {
			return nil, fmt.Errorf("service.getSemesterCourses: Get adjust course failed: %w", err)
		}

		for _, c := range list {
			jwchRules := pack.ToJwchScheduleRules(c.ScheduleRules)
			adjustRules := getAdjustRules(jwchRules, adjustCourses)
			c.ScheduleRules = pack.FromJwchScheduleRules(jwch.ApplyAdjustRules(jwchRules, adjustRules))
		}
	}

	// 写入 cache
	s.taskQueue.Add(courseKey, taskqueue.QueueTask{Execute: func() error {
		return cache.SetSliceCache(s.cache, s.ctx, courseKey, list,
			constants.CourseTermsKeyExpire, "Course.SetCourseCache")
	}})
	return list, nil
}

func getAdjustRules(scheduleRules []jwch.CourseScheduleRule, adjustCourses []*model.AutoAdjustCourse) (adjustRules []jwch.CourseAdjustRule) {
	for _, c := range adjustCourses {
		if !c.Enabled {
			continue
		}

		fromWeek := int(c.FromWeek)
		fromWeekday := int(c.FromWeekday)

		canceled := c.ToDate == nil

		for _, r := range scheduleRules {
			if r.StartWeek <= fromWeek && r.EndWeek >= fromWeek && r.Weekday == fromWeekday {
				if canceled {
					adjustRules = append(adjustRules, jwch.CourseAdjustRule{
						OldWeek:       fromWeek,
						OldWeekday:    r.Weekday,
						OldStartClass: r.StartClass,
						OldEndClass:   r.EndClass,
						Canceled:      true,
					})
					continue
				}

				toWeek := int(*c.ToWeek)
				toWeekday := int(*c.ToWeekday)

				adjustRules = append(adjustRules, jwch.CourseAdjustRule{
					OldWeek:       fromWeek,
					OldWeekday:    r.Weekday,
					OldStartClass: r.StartClass,
					OldEndClass:   r.EndClass,
					Canceled:      false,
					NewWeek:       toWeek,
					NewWeekday:    toWeekday,
					NewStartClass: r.StartClass,
					NewEndClass:   r.EndClass,
					NewLocation:   r.Location,
				})
			}
		}
	}

	return adjustRules
}
