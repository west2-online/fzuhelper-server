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

package course

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/internal/course/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	dbModel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/singleflight"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// CourseServiceImpl implements the last service interface defined in the IDL.
type CourseServiceImpl struct {
	ClientSet *base.ClientSet
	taskQueue taskqueue.TaskQueue
}

func NewCourseService(clientSet *base.ClientSet, taskQueue taskqueue.TaskQueue) *CourseServiceImpl {
	return &CourseServiceImpl{
		ClientSet: clientSet,
		taskQueue: taskQueue,
	}
}

// GetCourseList implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) GetCourseList(ctx context.Context, req *course.CourseListRequest) (resp *course.CourseListResponse, err error) {
	resp = course.NewCourseListResponse()
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Course.GetCourseList: Get login data fail %w", err)
	}
	stuId := loginData.Id
	isGraduate := utils.IsGraduate(stuId)
	isRefresh := req.IsRefresh != nil && *req.IsRefresh
	key := singleflight.Key(constants.SingleflightCourseListPrefix, stuId, req.Term, isGraduate, isRefresh)

	res, err := singleflight.Do(key, func() ([]*model.Course, error) {
		svc := service.NewCourseService(ctx, s.ClientSet, s.taskQueue)
		if isGraduate {
			return svc.GetCourseListYjsy(req, loginData)
		} else {
			// 检查学期是否合法的逻辑在 service 里面实现了，这里不需要再检查
			// 原因：GetSemesterCourses() 要用到 jwch 里面的 GetTerms() 函数返回的 ViewState 和 EventValidation 参数，顺便检查可以减少请求次数
			return svc.GetCourseList(req, loginData)
		}
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = res

	// 获取自定义课程（降级处理：获取失败不影响主流程）
	customCourses, err := s.getCustomCourses(ctx, stuId, req.Term)
	if err != nil {
		logger.WithCtx(ctx).Errorf("get custom courses failed (fallback to empty): %v", err)
		resp.CustomCourses = nil
	} else {
		resp.CustomCourses = customCourses
	}

	return resp, nil
}

// getCustomCourses 获取自定义课程列表
func (s *CourseServiceImpl) getCustomCourses(ctx context.Context, stuId, term string) ([]*course.CustomCourseItem, error) {
	dbClient := s.ClientSet.DBClient
	customCourses, err := dbClient.Course.GetCustomCourses(ctx, stuId, term)
	if err != nil {
		return nil, err
	}

	// 转换为 Thrift 类型
	return pack.BuildCustomCourseItems(customCourses), nil
}

// UpsertCustomCourse 新增或更新自定义课程
func (s *CourseServiceImpl) UpsertCustomCourse(ctx context.Context, req *course.UpsertCustomCourseRequest) (resp *course.UpsertCustomCourseResponse, err error) {
	resp = course.NewUpsertCustomCourseResponse()
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Course.UpsertCustomCourse: Get login data fail %w", err)
	}
	stuId := loginData.Id
	dbClient := s.ClientSet.DBClient

	courseItem := req.Course
	courseId := courseItem.Id

	if courseId == nil || *courseId == "" {
		// 新增：先检查是否已存在相同课程（多端去重）
		isDuplicate, err := dbClient.Course.CheckDuplicateCustomCourse(ctx, stuId, req.Term,
			courseItem.Name, courseItem.Location,
			int(courseItem.StartClass), int(courseItem.EndClass),
			int(courseItem.StartWeek), int(courseItem.EndWeek),
			int(courseItem.Weekday))
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return resp, nil
		}

		if isDuplicate {
			// 已存在相同课程，跳过（去重）
			resp.Base = base.BuildSuccessResp()
			resp.CourseId = nil
			return resp, nil
		}

		// 服务端生成 courseId 并保存
		newCourseId := uuid.New().String()
		courseId = &newCourseId

		customCourse := &dbModel.UserCustomCourse{
			StuId:      stuId,
			Term:       req.Term,
			CourseId:   newCourseId,
			Name:       courseItem.Name,
			Teacher:    getStringValue(courseItem.Teacher),
			Location:   courseItem.Location,
			StartClass: int(courseItem.StartClass),
			EndClass:   int(courseItem.EndClass),
			StartWeek:  int(courseItem.StartWeek),
			EndWeek:    int(courseItem.EndWeek),
			Weekday:    int(courseItem.Weekday),
			IsSingle:   getBoolValue(courseItem.Single),
			IsDouble:   getBoolValue(courseItem.Double_),
			Color:      getStringValueWithDefault(courseItem.Color, "#FF5733"),
			Remark:     getStringValue(courseItem.Remark),
		}

		if err := dbClient.Course.CreateCustomCourse(ctx, customCourse); err != nil {
			resp.Base = base.BuildBaseResp(err)
			return resp, nil
		}
	} else {
		// 更新：根据 course_id 更新
		updates := map[string]interface{}{
			"name":        courseItem.Name,
			"teacher":     getStringValue(courseItem.Teacher),
			"location":    courseItem.Location,
			"start_class": int(courseItem.StartClass),
			"end_class":   int(courseItem.EndClass),
			"start_week":  int(courseItem.StartWeek),
			"end_week":    int(courseItem.EndWeek),
			"weekday":     int(courseItem.Weekday),
			"is_single":   getBoolValue(courseItem.Single),
			"is_double":   getBoolValue(courseItem.Double_),
			"color":       getStringValueWithDefault(courseItem.Color, "#FF5733"),
			"remark":      getStringValue(courseItem.Remark),
		}

		if err := dbClient.Course.UpdateCustomCourse(ctx, stuId, req.Term, *courseId, updates); err != nil {
			resp.Base = base.BuildBaseResp(err)
			return resp, nil
		}
	}

	resp.Base = base.BuildSuccessResp()
	resp.CourseId = courseId
	return resp, nil
}

// DeleteCustomCourse 删除自定义课程
func (s *CourseServiceImpl) DeleteCustomCourse(ctx context.Context, req *course.DeleteCustomCourseRequest) (resp *course.DeleteCustomCourseResponse, err error) {
	resp = course.NewDeleteCustomCourseResponse()
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Course.DeleteCustomCourse: Get login data fail %w", err)
	}
	stuId := loginData.Id
	dbClient := s.ClientSet.DBClient

	if err := dbClient.Course.DeleteCustomCourse(ctx, stuId, req.Term, req.CourseId); err != nil {
		resp.Base = base.BuildBaseResp(errno.InternalServiceError.WithError(err))
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	return resp, nil
}

// 辅助函数：获取字符串值
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// 辅助函数：获取字符串值，带默认值
func getStringValueWithDefault(s *string, defaultVal string) string {
	if s == nil || *s == "" {
		return defaultVal
	}
	return *s
}

// 辅助函数：获取布尔值
func getBoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func (s *CourseServiceImpl) GetTermList(ctx context.Context, req *course.TermListRequest) (resp *course.TermListResponse, err error) {
	resp = course.NewTermListResponse()
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Course.GetTermList: Get login data fail %w", err)
	}
	stuId := loginData.Id
	isGraduate := utils.IsGraduate(stuId)
	key := singleflight.Key(constants.SingleflightCourseTermsPrefix, stuId, isGraduate)

	res, err := singleflight.Do(key, func() ([]string, error) {
		svc := service.NewCourseService(ctx, s.ClientSet, s.taskQueue)
		if isGraduate {
			return svc.GetTermsListYjsy(loginData)
		} else {
			return svc.GetTermsList(loginData)
		}
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = res
	return resp, nil
}

func (s *CourseServiceImpl) GetCalendar(ctx context.Context, req *course.GetCalendarRequest) (resp *course.GetCalendarResponse, err error) {
	resp = course.NewGetCalendarResponse()

	resp.Ics, err = service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetCalendar(req.StuId)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()

	return resp, nil
}

func (s *CourseServiceImpl) GetLocateDate(ctx context.Context, _ *course.GetLocateDateRequest) (resp *course.GetLocateDateResponse, err error) {
	resp = course.NewGetLocateDateResponse()

	res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetLocateDate()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.LocateDate = res
	return resp, nil
}

func (s *CourseServiceImpl) GetFriendCourse(ctx context.Context, req *course.GetFriendCourseRequest) (
	resp *course.GetFriendCourseResponse, err error,
) {
	resp = new(course.GetFriendCourseResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("Course.GetFriendCourse: Get login data fail %w", err)
	}
	res, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetFriendCourse(req, loginData)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = res
	return resp, nil
}

func (s *CourseServiceImpl) GetAutoAdjustCourseList(ctx context.Context, req *course.GetAutoAdjustCourseListRequest) (
	resp *course.GetAutoAdjustCourseListResponse, err error,
) {
	resp = new(course.GetAutoAdjustCourseListResponse)

	list, err := service.NewCourseService(ctx, s.ClientSet, s.taskQueue).GetAutoAdjustCourseList(req.Term)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = pack.BuildAdjustCourseList(list)
	return resp, nil
}

func (s *CourseServiceImpl) UpdateAdjustCourse(ctx context.Context, req *course.UpdateAdjustCourseRequest) (
	resp *course.UpdateAdjustCourseResponse, err error,
) {
	resp = new(course.UpdateAdjustCourseResponse)
	err = service.NewCourseService(ctx, s.ClientSet, s.taskQueue).UpdateAutoAdjustCourse(req)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	return resp, nil
}
