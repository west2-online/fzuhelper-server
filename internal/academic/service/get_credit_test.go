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

func TestAcademicService_GetCredit(t *testing.T) {
	Convey("GetCredit", t, func() {

		Convey("should return error when user is not logged in", func() {
			// Given: 未登录的用户上下文
			ctx := context.Background()
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取学分信息
			result, err := service.GetCredit()

			// Then: 应该返回登录错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get login data fail")
		})

		Convey("should return error when remote service is unavailable", func() {
			// Given: 已登录用户但远程服务不可用
			testLoginData := &model.LoginData{
				Id:      "test_student_id",
				Cookies: "test_session=abc123",
			}

			getCreditPatch := mockey.Mock((*jwch.Student).GetCredit).Return(
				nil, fmt.Errorf("network connection failed"),
			).Build()
			defer getCreditPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取学分信息
			result, err := service.GetCredit()

			// Then: 应该返回网络错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get credit info fail")
		})

		Convey("should return credit statistics when request is successful", func() {
			// Given: 已登录用户且系统正常
			testLoginData := &model.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			expectedCreditStats := []*jwch.CreditStatistics{
				{
					Type:  "公共基础必修课",
					Gain:  "29.5",
					Total: "32",
				},
				{
					Type:  "学科基础必修课",
					Gain:  "54",
					Total: "54",
				},
			}

			getCreditPatch := mockey.Mock((*jwch.Student).GetCredit).Return(
				expectedCreditStats, nil,
			).Build()
			defer getCreditPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 获取学分信息
			result, err := service.GetCredit()

			// Then: 应该返回正确的学分统计数据
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(len(result), ShouldEqual, 2)

			So(result[0].Type, ShouldEqual, "公共基础必修课")
			So(result[0].Gain, ShouldEqual, "29.5")
			So(result[0].Total, ShouldEqual, "32")

			So(result[1].Type, ShouldEqual, "学科基础必修课")
			So(result[1].Gain, ShouldEqual, "54")
			So(result[1].Total, ShouldEqual, "54")
		})
	})
}
