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

func TestAcademicService_GetGPA(t *testing.T) {
	Convey("GetGPA", t, func() {

		Convey("should return error when user is not logged in", func() {
			// Given: 未登录的用户上下文
			ctx := context.Background()
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取GPA信息
			result, err := service.GetGPA()

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

			getGPAPatch := mockey.Mock((*jwch.Student).GetGPA).Return(
				nil, fmt.Errorf("network connection failed"),
			).Build()
			defer getGPAPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 尝试获取GPA信息
			result, err := service.GetGPA()

			// Then: 应该返回网络错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get gpa info fail")
		})

		Convey("should return GPA data when request is successful", func() {
			// Given: 已登录用户且系统正常
			testLoginData := &model.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			expectedGPA := &jwch.GPABean{
				Time: "2023-06-01",
				Data: []jwch.GPAData{
					{
						Type:  "Mathematics",
						Value: "4.0",
					},
					{
						Type:  "Physics",
						Value: "3.5",
					},
				},
			}

			getGPAPatch := mockey.Mock((*jwch.Student).GetGPA).Return(
				expectedGPA, nil,
			).Build()
			defer getGPAPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			service := &AcademicService{ctx: ctx}

			// When: 获取GPA信息
			result, err := service.GetGPA()

			// Then: 应该返回正确的GPA数据
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(result.Time, ShouldEqual, "2023-06-01")
			So(len(result.Data), ShouldEqual, 2)

			So(result.Data[0].Type, ShouldEqual, "Mathematics")
			So(result.Data[0].Value, ShouldEqual, "4.0")
			So(result.Data[1].Type, ShouldEqual, "Physics")
			So(result.Data[1].Value, ShouldEqual, "3.5")
		})
	})
}
