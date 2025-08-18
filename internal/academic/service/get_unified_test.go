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

func TestAcademicService_GetUnifiedExam(t *testing.T) {
	Convey("GetUnifiedExam", t, func() {

		Convey("should return error when user is not logged in", func() {
			// Given: 未登录的用户上下文
			ctx := context.Background()
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取统一考试信息
			result, err := service.GetUnifiedExam()

			// Then: 应该返回登录错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get login data fail")
		})

		Convey("should return error when CET service is unavailable", func() {
			// Given: 已登录用户但CET服务不可用
			testLoginData := &model.LoginData{
				Id:      "test_student_id",
				Cookies: "test_session=abc123",
			}

			getCETPatch := mockey.Mock((*jwch.Student).GetCET).Return(
				nil, fmt.Errorf("CET service unavailable"),
			).Build()
			defer getCETPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取统一考试信息
			result, err := service.GetUnifiedExam()

			// Then: 应该返回CET错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get cet info fail")
		})

		Convey("should return error when JS service is unavailable", func() {
			// Given: 已登录用户，CET服务正常但JS服务不可用
			testLoginData := &model.LoginData{
				Id:      "test_student_id",
				Cookies: "test_session=abc123",
			}

			cetExams := []*jwch.UnifiedExam{
				{
					Name:  "CET-4",
					Score: "520",
					Term:  "2023年06月",
				},
			}

			getCETPatch := mockey.Mock((*jwch.Student).GetCET).Return(
				cetExams, nil,
			).Build()
			defer getCETPatch.UnPatch()

			getJSPatch := mockey.Mock((*jwch.Student).GetJS).Return(
				nil, fmt.Errorf("JS service unavailable"),
			).Build()
			defer getJSPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取统一考试信息
			result, err := service.GetUnifiedExam()

			// Then: 应该返回JS错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get js info fail")
		})

		Convey("should return unified exam data when request is successful", func() {
			// Given: 已登录用户且所有服务正常
			testLoginData := &model.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			cetExams := []*jwch.UnifiedExam{
				{
					Name:  "CET-4",
					Score: "520",
					Term:  "2023年06月",
				},
			}

			jsExams := []*jwch.UnifiedExam{
				{
					Name:  "全国计算机等级考试",
					Score: "85",
					Term:  "2023年03月",
				},
			}

			getCETPatch := mockey.Mock((*jwch.Student).GetCET).Return(
				cetExams, nil,
			).Build()
			defer getCETPatch.UnPatch()

			getJSPatch := mockey.Mock((*jwch.Student).GetJS).Return(
				jsExams, nil,
			).Build()
			defer getJSPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 获取统一考试信息
			result, err := service.GetUnifiedExam()

			// Then: 应该返回合并的考试数据
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(len(result), ShouldEqual, 2)

			// 验证CET数据
			So(result[0].Name, ShouldEqual, "CET-4")
			So(result[0].Score, ShouldEqual, "520")
			So(result[0].Term, ShouldEqual, "2023年06月")

			// 验证JS数据
			So(result[1].Name, ShouldEqual, "全国计算机等级考试")
			So(result[1].Score, ShouldEqual, "85")
			So(result[1].Term, ShouldEqual, "2023年03月")
		})
	})
}
