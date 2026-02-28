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

	"github.com/west2-online/fzuhelper-server/config"
	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	baseContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	academicCache "github.com/west2-online/fzuhelper-server/pkg/cache/academic"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	academicDB "github.com/west2-online/fzuhelper-server/pkg/db/academic"
	dbModel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/umeng"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

// 测试初始化函数
func init() {
	// 初始化测试配置，这会读取config.example.yaml
	if err := config.InitForTest("academic"); err != nil {
		panic(fmt.Sprintf("Failed to initialize test config: %v", err))
	}
}

func TestAcademicService_GetScores(t *testing.T) {
	Convey("GetScores", t, func() {
		Convey("should return error when user is not logged in", func() {
			// Given: 未登录的用户上下文
			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 尝试获取成绩信息
			result, err := service.GetScores(&loginmodel.LoginData{
				Id:      "test_student_id",
				Cookies: "test_session=abc123",
			})

			// Then: 应该返回登录错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get login data fail")
		})

		Convey("should return scores from cache when cache exists", func() {
			// Given: 已登录用户且缓存中有成绩数据
			testLoginData := &loginmodel.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			expectedScores := []*jwch.Mark{
				{
					Name:    "数据结构",
					Score:   "90",
					Credits: "4.0",
					GPA:     "3.9",
				},
				{
					Name:    "计算机网络",
					Score:   "85",
					Credits: "3.0",
					GPA:     "3.6",
				},
			}

			// Mock 缓存存在
			cacheExistsPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(true).Build()
			defer cacheExistsPatch.UnPatch()

			// Mock 从缓存获取成绩
			getCachePatch := mockey.Mock((*academicCache.CacheAcademic).GetScoresCache).Return(
				expectedScores, nil,
			).Build()
			defer getCachePatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 获取成绩信息
			result, err := service.GetScores(testLoginData)

			// Then: 应该返回缓存中的成绩数据
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(len(result), ShouldEqual, 2)
			So(result[0].Name, ShouldEqual, "数据结构")
			So(result[0].Score, ShouldEqual, "90")
			So(result[1].Name, ShouldEqual, "计算机网络")
			So(result[1].Score, ShouldEqual, "85")
		})

		Convey("should return error when cache exists but cache retrieval fails", func() {
			// Given: 已登录用户但缓存读取失败
			testLoginData := &loginmodel.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			// Mock 缓存存在
			cacheExistsPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(true).Build()
			defer cacheExistsPatch.UnPatch()

			// Mock 缓存读取失败
			getCachePatch := mockey.Mock((*academicCache.CacheAcademic).GetScoresCache).Return(
				nil, fmt.Errorf("redis connection failed"),
			).Build()
			defer getCachePatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 尝试获取成绩信息
			result, err := service.GetScores(testLoginData)

			// Then: 应该返回缓存错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get scores info from redis error")
		})

		Convey("should return scores from jwch when cache does not exist", func() {
			// Given: 已登录用户且缓存为空，需要从jwch获取
			testLoginData := &loginmodel.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			expectedScores := []*jwch.Mark{
				{
					Name:    "操作系统",
					Score:   "88",
					Credits: "3.0",
					GPA:     "3.7",
				},
			}

			// Mock 缓存不存在
			cacheExistsPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(false).Build()
			defer cacheExistsPatch.UnPatch()

			// Mock jwch 获取成绩成功
			getMarksPatch := mockey.Mock((*jwch.Student).GetMarks).Return(
				expectedScores, nil,
			).Build()
			defer getMarksPatch.UnPatch()

			// Mock 任务队列（防止真实执行异步任务）
			taskQueuePatch := mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()
			defer taskQueuePatch.UnPatch()

			// Mock umeng 推送（防止真实发送推送，因为任务队列可能会触发）
			umengAndroidPatch := mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(nil).Build()
			defer umengAndroidPatch.UnPatch()

			umengIOSPatch := mockey.Mock(umeng.SendIOSGroupcast).Return(nil).Build()
			defer umengIOSPatch.UnPatch()

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 获取成绩信息
			result, err := service.GetScores(testLoginData)

			// Then: 应该返回从jwch获取的成绩数据
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(len(result), ShouldEqual, 1)
			So(result[0].Name, ShouldEqual, "操作系统")
			So(result[0].Score, ShouldEqual, "88")
		})

		Convey("should return error when cache does not exist and jwch service fails", func() {
			// Given: 已登录用户，缓存为空，jwch服务不可用
			testLoginData := &loginmodel.LoginData{
				Id:      "222200311",
				Cookies: "ASP.NET_SessionId=lzs1t42mpkml4ag2jrxvib4z",
			}

			// Mock 缓存不存在
			cacheExistsPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(false).Build()
			defer cacheExistsPatch.UnPatch()

			// Mock jwch 获取成绩失败
			getMarksPatch := mockey.Mock((*jwch.Student).GetMarks).Return(
				nil, fmt.Errorf("network connection failed"),
			).Build()
			defer getMarksPatch.UnPatch()

			// 注意：jwch失败时不会走到任务队列，所以不需要mock umeng

			ctx := baseContext.WithLoginData(context.Background(), testLoginData)
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 尝试获取成绩信息
			result, err := service.GetScores(testLoginData)

			// Then: 应该返回jwch错误
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Get scores info fail")
		})
	})
}

func TestAcademicService_checkScoreChange(t *testing.T) {
	Convey("checkScoreChange", t, func() {
		Convey("should create new score record when student has no score history", func() {
			// Given: 学生没有成绩历史记录
			testScores := []*jwch.Mark{
				{
					Name:    "数据结构",
					Score:   "90",
					Credits: "4.0",
					GPA:     "3.9",
				},
			}

			// Mock 数据库查询返回空的SHA256（表示没有历史记录）
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return("", nil).Build()
			defer getSha256Patch.UnPatch()

			// Mock 创建用户成绩记录
			createScorePatch := mockey.Mock((*academicDB.DBAcademic).CreateUserScore).Return(&dbModel.Score{
				StuID:            "222200311",
				ScoresInfo:       `[{"name":"数据结构","score":"90","credits":"4.0","gpa":"3.9","semester":"2024-1","teacher":"张老师","electiveType":"必修"}]`,
				ScoresInfoSHA256: "new_sha256",
			}, nil).Build()
			defer createScorePatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err := service.checkScoreChange("222200311", testScores)

			// Then: 应该成功创建新记录
			So(err, ShouldBeNil)
		})

		Convey("should update score record when scores have changed", func() {
			// Given: 学生有成绩历史记录，包含多门课程，部分成绩已更新
			testScores := []*jwch.Mark{
				{
					Name:         "数据结构",
					Score:        "95", // 分数更新了：90 -> 95
					Credits:      "4.0",
					GPA:          "4.0",
					Semester:     "2024-1",
					Teacher:      "张老师",
					ElectiveType: "必修",
				},
				{
					Name:         "计算机网络",
					Score:        "88", // 分数更新了：85 -> 88
					Credits:      "3.0",
					GPA:          "3.7",
					Semester:     "2024-1",
					Teacher:      "李老师",
					ElectiveType: "必修",
				},
				{
					Name:         "软件工程",
					Score:        "92", // 分数没有变化，仍然是92
					Credits:      "3.5",
					GPA:          "4.0",
					Semester:     "2024-1",
					Teacher:      "王老师",
					ElectiveType: "选修",
				},
			}

			// Mock 返回旧的SHA256（表示有历史记录）
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return("old_sha256", nil).Build()
			defer getSha256Patch.UnPatch()

			// Mock 获取旧成绩数据
			getScorePatch := mockey.Mock((*academicDB.DBAcademic).GetScoreByStuId).Return(&dbModel.Score{
				StuID: "222200311",
				ScoresInfo: `[{"name":"数据结构","score":"90","credits":"4.0","gpa":"3.9","semester":"2024-1","teacher":"张老师","electiveType":"必修"},
				{"name":"计算机网络","score":"85","credits":"3.0","gpa":"3.6","semester":"2024-1","teacher":"李老师","electiveType":"必修"},
				{"name":"软件工程","score":"92","credits":"3.5","gpa":"4.0","semester":"2024-1","teacher":"王老师","electiveType":"选修"}]`,
				ScoresInfoSHA256: "old_sha256",
			}, nil).Build()
			defer getScorePatch.UnPatch()

			// Mock 课程不存在（需要发送推送）
			getCourseByHashPatch := mockey.Mock((*academicDB.DBAcademic).GetCourseByHash).Return(nil, nil).Build()
			defer getCourseByHashPatch.UnPatch()

			// Mock 创建课程记录
			createCoursePatch := mockey.Mock((*academicDB.DBAcademic).CreateCourseOffering).Return(&dbModel.CourseOffering{}, nil).Build()
			defer createCoursePatch.UnPatch()

			// Mock 更新成绩记录
			updateScorePatch := mockey.Mock((*academicDB.DBAcademic).UpdateUserScores).Return(nil).Build()
			defer updateScorePatch.UnPatch()

			// Mock umeng 推送
			umengAndroidPatch := mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(nil).Build()
			defer umengAndroidPatch.UnPatch()

			umengIOSPatch := mockey.Mock(umeng.SendIOSGroupcast).Return(nil).Build()
			defer umengIOSPatch.UnPatch()
			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err := service.checkScoreChange("222200311", testScores)

			// Then: 应该成功更新记录并发送推送
			So(err, ShouldBeNil)
		})

		Convey("should not send notification when course notification already sent", func() {
			// Given: 学生成绩更新但已经发送过通知
			testScores := []*jwch.Mark{
				{
					Name:         "数据结构",
					Score:        "95",
					Credits:      "4.0",
					GPA:          "4.0",
					Semester:     "2024-1",
					Teacher:      "张老师",
					ElectiveType: "必修",
				},
			}

			// Mock 返回旧的SHA256
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return("old_sha256", nil).Build()
			defer getSha256Patch.UnPatch()

			// Mock 获取旧成绩数据
			getScorePatch := mockey.Mock((*academicDB.DBAcademic).GetScoreByStuId).Return(&dbModel.Score{
				StuID:            "222200311",
				ScoresInfo:       `[{"name":"数据结构","score":"90","credits":"4.0","gpa":"3.9","semester":"2024-1","teacher":"张老师","electiveType":"必修"}]`,
				ScoresInfoSHA256: "old_sha256",
			}, nil).Build()
			defer getScorePatch.UnPatch()

			// Mock 课程存在（已发送过通知）
			getCourseByHashPatch := mockey.Mock((*academicDB.DBAcademic).GetCourseByHash).Return(&dbModel.CourseOffering{
				Name:         "数据结构",
				Term:         "2024-1",
				Teacher:      "张老师",
				ElectiveType: "必修",
				CourseHash:   "test_hash",
			}, nil).Build()
			defer getCourseByHashPatch.UnPatch()

			// Mock 更新成绩记录
			updateScorePatch := mockey.Mock((*academicDB.DBAcademic).UpdateUserScores).Return(nil).Build()
			defer updateScorePatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err := service.checkScoreChange("222200311", testScores)

			// Then: 应该成功更新记录但不发送推送
			So(err, ShouldBeNil)
		})

		Convey("should do nothing when scores have not changed", func() {
			// Given: 学生成绩没有变化
			testScores := []*jwch.Mark{
				{
					Name:    "数据结构",
					Score:   "90",
					Credits: "4.0",
					GPA:     "3.9",
				},
			}
			json, err := utils.JSONEncode(testScores)
			if err != nil {
				t.Fatal(err)
			}
			sha256 := utils.SHA256(json)
			// Mock 返回相同的SHA256（表示成绩没有变化）
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return(sha256, nil).Build()
			defer getSha256Patch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err = service.checkScoreChange("222200311", testScores)

			// Then: 应该不做任何操作
			So(err, ShouldBeNil)
		})

		Convey("should return error when GetScoreSha256ByStuId fails", func() {
			// Given: 获取SHA256时发生数据库错误
			testScores := []*jwch.Mark{
				{
					Name:    "数据结构",
					Score:   "90",
					Credits: "4.0",
					GPA:     "3.9",
				},
			}

			// Mock 数据库查询返回错误
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return("", fmt.Errorf("database connection failed")).Build()
			defer getSha256Patch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err := service.checkScoreChange("222200311", testScores)

			// Then: 应该返回数据库错误
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "database connection failed")
		})

		Convey("should return error when CreateUserScore fails", func() {
			// Given: 学生没有成绩历史记录，但创建新记录失败
			testScores := []*jwch.Mark{
				{
					Name:    "数据结构",
					Score:   "90",
					Credits: "4.0",
					GPA:     "3.9",
				},
			}

			// Mock 返回空的SHA256（表示没有历史记录）
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return("", nil).Build()
			defer getSha256Patch.UnPatch()

			// Mock 创建用户成绩记录失败
			createScorePatch := mockey.Mock((*academicDB.DBAcademic).CreateUserScore).Return(nil, fmt.Errorf("insert failed")).Build()
			defer createScorePatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err := service.checkScoreChange("222200311", testScores)

			// Then: 应该返回创建错误
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "insert failed")
		})

		Convey("should return error when handleScoreChange fails", func() {
			// Given: 成绩已变化，但处理过程中出错
			testScores := []*jwch.Mark{
				{
					Name:         "数据结构",
					Score:        "95",
					Credits:      "4.0",
					GPA:          "4.0",
					Semester:     "2024-1",
					Teacher:      "张老师",
					ElectiveType: "必修",
				},
			}

			// Mock 返回旧的SHA256（表示有历史记录）
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return("old_sha256", nil).Build()
			defer getSha256Patch.UnPatch()

			// Mock 获取旧成绩数据失败（这会导致handleScoreChange失败）
			getScorePatch := mockey.Mock((*academicDB.DBAcademic).GetScoreByStuId).Return(nil, fmt.Errorf("query failed")).Build()
			defer getScorePatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err := service.checkScoreChange("222200311", testScores)

			// Then: 应该返回处理错误
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "query failed")
		})

		Convey("should return error when UpdateUserScores fails", func() {
			// Given: 成绩已变化，但更新数据库失败
			testScores := []*jwch.Mark{
				{
					Name:         "数据结构",
					Score:        "95",
					Credits:      "4.0",
					GPA:          "4.0",
					Semester:     "2024-1",
					Teacher:      "张老师",
					ElectiveType: "必修",
				},
			}

			// Mock 返回旧的SHA256（表示有历史记录）
			getSha256Patch := mockey.Mock((*academicDB.DBAcademic).GetScoreSha256ByStuId).Return("old_sha256", nil).Build()
			defer getSha256Patch.UnPatch()

			// Mock 获取旧成绩数据
			getScorePatch := mockey.Mock((*academicDB.DBAcademic).GetScoreByStuId).Return(&dbModel.Score{
				StuID:            "222200311",
				ScoresInfo:       `[{"name":"数据结构","score":"90","credits":"4.0","gpa":"3.9","semester":"2024-1","teacher":"张老师","electiveType":"必修"}]`,
				ScoresInfoSHA256: "old_sha256",
			}, nil).Build()
			defer getScorePatch.UnPatch()

			// Mock 课程不存在（需要发送推送）
			getCourseByHashPatch := mockey.Mock((*academicDB.DBAcademic).GetCourseByHash).Return(nil, nil).Build()
			defer getCourseByHashPatch.UnPatch()

			// Mock 创建课程记录
			createCoursePatch := mockey.Mock((*academicDB.DBAcademic).CreateCourseOffering).Return(&dbModel.CourseOffering{
				Name:         "数据结构",
				Term:         "2024-1",
				Teacher:      "张老师",
				ElectiveType: "必修",
				CourseHash:   "test_hash",
			}, nil).Build()
			defer createCoursePatch.UnPatch()

			// Mock 更新成绩记录失败
			updateScorePatch := mockey.Mock((*academicDB.DBAcademic).UpdateUserScores).Return(fmt.Errorf("update failed")).Build()
			defer updateScorePatch.UnPatch()

			// Mock umeng 推送
			umengAndroidPatch := mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(nil).Build()
			defer umengAndroidPatch.UnPatch()

			umengIOSPatch := mockey.Mock(umeng.SendIOSGroupcast).Return(nil).Build()
			defer umengIOSPatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				DBClient: &db.Database{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 检查成绩变化
			err := service.checkScoreChange("222200311", testScores)

			// Then: 应该返回更新错误
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "update failed")
		})
	})
}

func TestAcademicService_sendNotifications(t *testing.T) {
	Convey("sendNotifications", t, func() {
		Convey("should send notifications to both Android and iOS", func() {
			// Given: 准备发送推送的课程信息
			courseName := "数据结构"
			tag := "abcdefghijklmnopqrstuvwxyz123456"

			// Mock umeng 推送成功
			umengAndroidPatch := mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(nil).Build()
			defer umengAndroidPatch.UnPatch()

			umengIOSPatch := mockey.Mock(umeng.SendIOSGroupcast).Return(nil).Build()
			defer umengIOSPatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 发送通知
			err := service.sendNotifications(courseName, tag)

			// Then: 应该成功发送推送
			So(err, ShouldBeNil)
		})

		Convey("should handle notification errors gracefully", func() {
			// Given: 准备发送推送但可能出错
			courseName := "数据结构"
			tag := "abcdefghijklmnopqrstuvwxyz123456"

			// Mock umeng 推送失败
			umengAndroidPatch := mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(fmt.Errorf("android push failed")).Build()
			defer umengAndroidPatch.UnPatch()

			umengIOSPatch := mockey.Mock(umeng.SendIOSGroupcast).Return(fmt.Errorf("ios push failed")).Build()
			defer umengIOSPatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 发送通知
			err := service.sendNotifications(courseName, tag)

			// Then: 应该成功处理错误（函数内部处理了错误，不会返回错误）
			So(err, ShouldBeNil)
		})
	})
}

func TestAcademicService_GetScoresYjsy(t *testing.T) {
	Convey("GetScoresYjsy", t, func() {
		Convey("should return scores from cache when cache exists", func() {
			// Given: 缓存中有研究生成绩数据
			testScores := []*yjsy.Mark{
				{
					Name:         "高等数学",
					Score:        "95",
					Credits:      "4.0",
					GPA:          "4.0",
					ExamType:     "正常考试",
					ElectiveType: "必修",
				},
			}

			loginData := &loginmodel.LoginData{
				Id:      "202212345678",
				Cookies: "test_cookie",
			}

			// Mock 缓存存在
			isKeyExistPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(true).Build()
			defer isKeyExistPatch.UnPatch()

			// Mock 获取缓存成功
			getCachePatch := mockey.Mock((*academicCache.CacheAcademic).GetScoresCacheYjsy).Return(testScores, nil).Build()
			defer getCachePatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 调用GetScoresYjsy
			scores, err := service.GetScoresYjsy(loginData)

			// Then: 应该返回缓存中的成绩
			So(err, ShouldBeNil)
			So(scores, ShouldNotBeNil)
			So(len(scores), ShouldEqual, 1)
			So(scores[0].Name, ShouldEqual, "高等数学")
		})

		Convey("should return error when cache exists but cache retrieval fails", func() {
			// Given: 缓存存在但获取失败
			loginData := &loginmodel.LoginData{
				Id:      "202212345678",
				Cookies: "test_cookie",
			}

			// Mock 缓存存在
			isKeyExistPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(true).Build()
			defer isKeyExistPatch.UnPatch()

			// Mock 获取缓存失败
			getCachePatch := mockey.Mock((*academicCache.CacheAcademic).GetScoresCacheYjsy).Return(nil, fmt.Errorf("cache error")).Build()
			defer getCachePatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 调用GetScoresYjsy
			scores, err := service.GetScoresYjsy(loginData)

			// Then: 应该返回错误
			So(err, ShouldNotBeNil)
			So(scores, ShouldBeNil)
			So(err.Error(), ShouldContainSubstring, "cache error")
		})

		Convey("should return scores from yjsy when cache does not exist", func() {
			// Given: 缓存不存在，需要从研究生院系统获取
			testScores := []*yjsy.Mark{
				{
					Name:         "高等数学",
					Score:        "95",
					Credits:      "4.0",
					GPA:          "4.0",
					ExamType:     "正常考试",
					ElectiveType: "必修",
				},
			}

			loginData := &loginmodel.LoginData{
				Id:      "202212345678",
				Cookies: "test_cookie",
			}

			// Mock 缓存不存在
			isKeyExistPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(false).Build()
			defer isKeyExistPatch.UnPatch()

			withLoginDataPatch := mockey.Mock((*yjsy.Student).WithLoginData).Return(yjsy.NewStudent()).Build()
			defer withLoginDataPatch.UnPatch()

			getMarksPatch := mockey.Mock((*yjsy.Student).GetMarks).Return(testScores, nil).Build()
			defer getMarksPatch.UnPatch()

			// Mock 任务队列（防止真实执行异步任务）
			taskQueuePatch := mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()
			defer taskQueuePatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 调用GetScoresYjsy
			scores, err := service.GetScoresYjsy(loginData)

			// Then: 应该从yjsy获取成绩并返回
			So(err, ShouldBeNil)
			So(scores, ShouldNotBeNil)
			So(len(scores), ShouldEqual, 1)
			So(scores[0].Name, ShouldEqual, "高等数学")
		})

		Convey("should return error when cache does not exist and yjsy service fails", func() {
			// Given: 缓存不存在且yjsy服务调用失败
			loginData := &loginmodel.LoginData{
				Id:      "202212345678",
				Cookies: "test_cookie",
			}

			// Mock 缓存不存在
			isKeyExistPatch := mockey.Mock((*cache.Cache).IsKeyExist).Return(false).Build()
			defer isKeyExistPatch.UnPatch()

			withLoginDataPatch := mockey.Mock((*yjsy.Student).WithLoginData).Return(yjsy.NewStudent()).Build()
			defer withLoginDataPatch.UnPatch()

			// Mock yjsy 获取成绩失败
			getMarksPatch := mockey.Mock((*yjsy.Student).GetMarks).Return(nil, fmt.Errorf("yjsy service error")).Build()
			defer getMarksPatch.UnPatch()

			ctx := context.Background()
			mockClientSet := &base.ClientSet{
				CacheClient: &cache.Cache{},
			}
			service := NewAcademicService(ctx, mockClientSet, &taskqueue.BaseTaskQueue{})

			// When: 调用GetScoresYjsy
			scores, err := service.GetScoresYjsy(loginData)

			// Then: 应该返回错误
			So(err, ShouldNotBeNil)
			So(scores, ShouldBeNil)
			So(err.Error(), ShouldContainSubstring, "Get scores info fail")
		})
	})
}
