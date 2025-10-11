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
				nil, nil, fmt.Errorf("network connection failed"),
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

		Convey("should return major credits when available", func() {
			// Given: 已登录用户且系统正常，包含有超出部分的学分
			testLoginData := &model.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			majorCredits := []*jwch.CreditStatistics{
				{
					Type:  "公共基础必修课",
					Gain:  "35.0",
					Total: "32.0",
				},
				{
					Type:  "学科基础必修课",
					Gain:  "54.0",
					Total: "54.5",
				},
				{
					Type:  "专业必修课",
					Gain:  "50.0",
					Total: "44.0",
				},
				{
					Type:  "专业选修课",
					Gain:  "5.0",
					Total: "12.0",
				},
				{
					Type:  "总计",
					Gain:  "144.0",
					Total: "142.5",
				},
			}

			// Mock jwch's GetCreditV2 to return the data
			getCreditV2Patch := mockey.Mock((*jwch.Student).GetCreditV2).Return(
				majorCredits, nil, nil, // majorCredits, minorCredits, error
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
			So((*result)[0].Data[0].Value, ShouldEqual, "35 / 32")
			// 验证有超出部分时的"总计"行
			So((*result)[0].Data[4].Key, ShouldEqual, "总计")
			So((*result)[0].Data[4].Value, ShouldEqual, "144 / 142.5 (还需 7.5 分)")
		})

		Convey("should return both major and minor credits when available", func() {
			// Given: 已登录用户且系统正常，有主修专业和辅修专业
			testLoginData := &model.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			// 主修专业数据
			majorCredits := []*jwch.CreditStatistics{
				{
					Type:  "公共基础必修课",
					Gain:  "20.0",
					Total: "20.0",
				},
				{
					Type:  "学科基础必修课",
					Gain:  "25.0",
					Total: "24.0",
				},
				{
					Type:  "专业选修课",
					Gain:  "15.0",
					Total: "12.0",
				},
				{
					Type:  "总计",
					Gain:  "60.0",
					Total: "56.0",
				},
			}

			// 辅修专业数据
			minorCredits := []*jwch.CreditStatistics{
				{
					Type:  "辅修必修课",
					Gain:  "15.0",
					Total: "18.0",
				},
				{
					Type:  "辅修选修课",
					Gain:  "12.0",
					Total: "10.0",
				},
				{
					Type:  "总计",
					Gain:  "27.0",
					Total: "28.0",
				},
			}

			// Mock jwch's GetCreditV2 to return the data
			getCreditV2Patch := mockey.Mock((*jwch.Student).GetCreditV2).Return(
				majorCredits, minorCredits, nil, // majorCredits, minorCredits, error
			).Build()
			defer getCreditV2Patch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 获取学分信息
			result, err := service.GetCreditV2()

			// Then: 应该返回主修专业和辅修专业数据
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(len(*result), ShouldEqual, 2)

			// Check main major credits
			So((*result)[0].Type, ShouldEqual, "主修专业")
			// 验证刚好修满的学分
			So((*result)[0].Data[0].Key, ShouldEqual, "公共基础必修课")
			So((*result)[0].Data[0].Value, ShouldEqual, "20 / 20")
			// 验证有超出部分的学分
			So((*result)[0].Data[1].Key, ShouldEqual, "学科基础必修课")
			So((*result)[0].Data[1].Value, ShouldEqual, "25 / 24")
			// 验证有超出部分的学分
			So((*result)[0].Data[2].Key, ShouldEqual, "专业选修课")
			So((*result)[0].Data[2].Value, ShouldEqual, "15 / 12")
			// 验证有超出部分时的"总计"行
			So((*result)[0].Data[3].Key, ShouldEqual, "总计")
			So((*result)[0].Data[3].Value, ShouldEqual, "60 / 56")

			// Check minor credits
			So((*result)[1].Type, ShouldEqual, "辅修专业")
			// 验证未修满的学分
			So((*result)[1].Data[0].Key, ShouldEqual, "辅修必修课")
			So((*result)[1].Data[0].Value, ShouldEqual, "15 / 18 (还需 3 分)")
			// 验证有超出部分的学分
			So((*result)[1].Data[1].Key, ShouldEqual, "辅修选修课")
			So((*result)[1].Data[1].Value, ShouldEqual, "12 / 10")
			// 验证有超出部分时的"总计"行
			So((*result)[1].Data[2].Key, ShouldEqual, "总计")
			So((*result)[1].Data[2].Value, ShouldEqual, "27 / 28 (还需 3 分)")
		})
	})
}
