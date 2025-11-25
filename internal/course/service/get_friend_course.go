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

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *CourseService) GetFriendCourse(req *course.GetFriendCourseRequest, loginData *kitexModel.LoginData) ([]*kitexModel.Course, error) {
	var err error
	stuId := context.ExtractIDFromLoginData(loginData)
	// 验证好友
	resp, err := s.userClient.VerifyFriend(s.ctx, &user.VerifyFriendRequest{Id: stuId, FriendId: req.Id})
	if err != nil {
		return nil, fmt.Errorf("CourseService.gerFriendCourse: verify friend failed: %w", err)
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}
	ok := resp.FriendExist
	if !ok {
		return nil, fmt.Errorf("service.GetFriendCourse: only friend course available")
	}
	termKey := fmt.Sprintf("terms:%s", req.Id)
	courseKey := fmt.Sprintf("course:%s:%s", req.Id, req.Term)
	/* 这里如果terms Cache没命中 无法验证term的合法性 也就会拒绝返回好友课表
	   而近term也是会在学生刷新课表时缓存 并且term似乎目前并不在db内存储
	   此外因为jwch与yjsy的区别 term也有两个结构 这边就直接用string来处理了
	*/
	var terms []string
	if s.cache.IsKeyExist(s.ctx, termKey) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, termKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetFriendCourse: Get term fail: %w", err)
		}
		terms = termsList
	}
	if !slices.Contains(terms, req.Term) && !slices.Contains(pack.GetTop2TermLists(terms), req.Term) {
		return nil, errors.New("service.GetFriendCourse: Invalid term")
	}
	/* cache 返回的两个course结构有区别 而目前判别研究生身份的方法需要loginData.Id
	在cache命中的情况下 先后两次尝试获取并返回课表
	*/
	if s.cache.IsKeyExist(s.ctx, courseKey) {
		courses, err := s.cache.Course.GetCoursesCache(s.ctx, courseKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetFriendCourse: Get courses fail: %w", err)
		}
		if courses != nil {
			return s.removeDuplicateCourses(pack.BuildCourse(courses)), nil
		}
		// cache 命中却没有course数据 做出查找研究生课表的尝试
		yjsyCourses, err := s.cache.Course.GetCoursesCacheYjsy(s.ctx, courseKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetYjsyFriendCourse: Get courses fail: %w", err)
		}
		return pack.BuildCourseYjsy(yjsyCourses), nil
	} else {
		var courses *model.UserCourse
		courses, err = s.db.Course.GetUserTermCourseByStuIdAndTerm(s.ctx, req.Id, req.Term)
		if err != nil {
			return nil, fmt.Errorf("service.GetSemesterCourses: Get courses fail: %w", err)
		}
		if courses == nil {
			return nil, errno.NewErrNo(errno.InternalServiceErrorCode, "service.GetSemesterCourses: there is no course in database, please login app and retry")
		}
		list := make([]*kitexModel.Course, 0)
		if err = sonic.Unmarshal([]byte(courses.TermCourses), &list); err != nil {
			return nil, fmt.Errorf("service.GetSemesterCourses: Unmarshal fail: %w", err)
		}
		return list, nil
	}
}
