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

func TestAcademicService_GetCreditV2(t *testing.T) {
	Convey("GetCreditV2", t, func() {
		Convey("should return error when user is not logged in", func() {
			// Given: 未登录的用户上下文
			ctx := context.Background()
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取学分信息
			result, err := service.GetCreditV2()

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

			getCreditV2Patch := mockey.Mock((*jwch.Student).GetCreditV2).Return(
				nil, fmt.Errorf("network connection failed"),
			).Build()
			defer getCreditV2Patch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取学分信息
			result, err := service.GetCreditV2()

			// Then: 应该返回网络错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get credit fail")
		})

		Convey("should return credit statistics when request is successful with overflow cases", func() {
			// Given: 已登录用户且系统正常，包含有超出部分的学分
			testLoginData := &model.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			// 构建预期的返回数据，包括超出部分的学分
			expectedCreditStats := jwch.CreditResponse{
				{
					Type: "主修专业",
					Data: []*jwch.CreditDetail{
						{
							Key:   "公共基础必修课",
							Value: "35.0/32.0(已修满)", // 35.0 > 32.0，有超出部分
						},
						{
							Key:   "学科基础必修课",
							Value: "54.0/54.5(还需0.5分)",
						},
						{
							Key:   "专业必修课",
							Value: "50.0/44.0(已修满)", // 50.0 > 44.0，有超出部分
						},
						{
							Key:   "专业选修课",
							Value: "5.0/12.0(还需7.0分)",
						},
						// 用于验证"总计"的还需学分
						// 总共还差：0.5 + 7.0 = 7.5
						{
							Key:   "总计",
							Value: "144.0/142.5(还需7.5分)",
						},
					},
				},
			}

			getCreditV2Patch := mockey.Mock((*jwch.Student).GetCreditV2).Return(
				expectedCreditStats, nil,
			).Build()
			defer getCreditV2Patch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 获取学分信息
			result, err := service.GetCreditV2()

			// Then: 应该返回正确的学分统计数据
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(len(*result), ShouldEqual, 1)

			So((*result)[0].Type, ShouldEqual, "主修专业")
			// 验证有超出部分的学分
			So((*result)[0].Data[0].Key, ShouldEqual, "公共基础必修课")
			So((*result)[0].Data[0].Value, ShouldEqual, "35.0/32.0(已修满)")
			// 验证有超出部分的学分
			So((*result)[0].Data[2].Key, ShouldEqual, "专业必修课")
			So((*result)[0].Data[2].Value, ShouldEqual, "50.0/44.0(已修满)")
			// 验证有超出部分时的"总计"行
			So((*result)[0].Data[4].Key, ShouldEqual, "总计")
			So((*result)[0].Data[4].Value, ShouldEqual, "144.0/142.5(还需7.5分)")
		})
	})
}
