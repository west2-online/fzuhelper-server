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
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	baseContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/jwch"
)

func TestAcademicService_GetPlan(t *testing.T) {
	Convey("GetPlan", t, func() {

		Convey("should return error when user is not logged in", func() {
			// Given: 未登录的用户上下文
			ctx := context.Background()
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取培养计划
			result, err := service.GetPlan()

			// Then: 应该返回登录错误
			So(result, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})

		Convey("should return error when remote service is unavailable", func() {
			// Given: 已登录用户但远程服务不可用
			testLoginData := &model.LoginData{
				Id:      "test_student_id",
				Cookies: "test_session=abc123",
			}

			getPlanPatch := mockey.Mock((*jwch.Student).GetCultivatePlan).Return(
				"", fmt.Errorf("cultivate plan not found"),
			).Build()
			defer getPlanPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取培养计划
			result, err := service.GetPlan()

			// Then: 应该返回网络错误
			So(result, ShouldEqual, "")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "AcademicService.GetPlan")
		})

		Convey("should return error when URL format is invalid", func() {
			// Given: 已登录用户但返回的URL格式不正确
			testLoginData := &model.LoginData{
				Id:      "test_student_id",
				Cookies: "test_session=abc123",
			}

			getPlanPatch := mockey.Mock((*jwch.Student).GetCultivatePlan).Return(
				"https://jwch.fzu.edu.cn/plan/view", nil, // 没有 &id 参数
			).Build()
			defer getPlanPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取培养计划
			result, err := service.GetPlan()

			// Then: 应该返回格式错误
			So(result, ShouldEqual, "")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "AcademicService.GetPlan")
		})

		Convey("should return plan URL when request is successful", func() {
			// Given: 已登录用户且系统正常
			testLoginData := &model.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			fullURL := "https://jwch.fzu.edu.cn/plan/view?type=cultivate&id=123456789"
			expectedURL := "https://jwch.fzu.edu.cn/plan/view?type=cultivate"

			getPlanPatch := mockey.Mock((*jwch.Student).GetCultivatePlan).Return(
				fullURL, nil,
			).Build()
			defer getPlanPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 获取培养计划
			result, err := service.GetPlan()

			// Then: 应该返回正确的URL（去掉id参数）
			So(err, ShouldBeNil)
			So(result, ShouldEqual, expectedURL)
		})
	})
}
