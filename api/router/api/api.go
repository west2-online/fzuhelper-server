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

// Code generated by hertz generator. DO NOT EDIT.

package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	api "github.com/west2-online/fzuhelper-server/api/handler/api"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_api := root.Group("/api", _apiMw()...)
		{
			_login := _api.Group("/login", _loginMw()...)
			_login.POST("/validateCode", append(_validatecodeforandroidMw(), api.ValidateCodeForAndroid)...)
		}
		{
			_v1 := _api.Group("/v1", _v1Mw()...)
			_v1.GET("/downloadUrl", append(_getdownloadurlforandroidMw(), api.GetDownloadUrlForAndroid)...)
			_v1.GET("/list", append(_listdirfilesforandroidMw(), api.ListDirFilesForAndroid)...)
			{
				_common := _v1.Group("/common", _commonMw()...)
				_common.GET("/notice", append(_getnoticeMw(), api.GetNotice)...)
				{
					_classroom := _common.Group("/classroom", _classroomMw()...)
					_classroom.GET("/empty", append(_getemptyclassroomsMw(), api.GetEmptyClassrooms)...)
				}
			}
			{
				_internal := _v1.Group("/internal", _internalMw()...)
				{
					_user := _internal.Group("/user", _userMw()...)
					_user.GET("/login", append(_getlogindataMw(), api.GetLoginData)...)
				}
			}
			{
				_jwch := _v1.Group("/jwch", _jwchMw()...)
				_jwch.GET("/ping", append(_testauthMw(), api.TestAuth)...)
				{
					_academic := _jwch.Group("/academic", _academicMw()...)
					_academic.GET("/credit", append(_getcreditMw(), api.GetCredit)...)
					_academic.GET("/gpa", append(_getgpaMw(), api.GetGPA)...)
					_academic.GET("/plan", append(_getplanMw(), api.GetPlan)...)
					_academic.GET("/scores", append(_getscoresMw(), api.GetScores)...)
					_academic.GET("/unified-exam", append(_getunifiedexamMw(), api.GetUnifiedExam)...)
				}
				{
					_classroom0 := _jwch.Group("/classroom", _classroom0Mw()...)
					_classroom0.GET("/exam", append(_getexamroominfoMw(), api.GetExamRoomInfo)...)
				}
				{
					_course := _jwch.Group("/course", _courseMw()...)
					_course.GET("/list", append(_getcourselistMw(), api.GetCourseList)...)
				}
				{
					_term := _jwch.Group("/term", _termMw()...)
					_term.GET("/list", append(_gettermlistMw(), api.GetTermList)...)
				}
				{
					_user0 := _jwch.Group("/user", _user0Mw()...)
					_user0.GET("/info", append(_getuserinfoMw(), api.GetUserInfo)...)
				}
			}
			{
				_launch_screen := _v1.Group("/launch-screen", _launch_screenMw()...)
				_launch_screen.DELETE("/image", append(_deleteimageMw(), api.DeleteImage)...)
				_image := _launch_screen.Group("/image", _imageMw()...)
				_image.GET("/point-time", append(_addimagepointtimeMw(), api.AddImagePointTime)...)
				_launch_screen.GET("/image", append(_getimageMw(), api.GetImage)...)
				_image0 := _launch_screen.Group("/image", _image0Mw()...)
				_image0.PUT("/property", append(_changeimagepropertyMw(), api.ChangeImageProperty)...)
				_launch_screen.POST("/image", append(_createimageMw(), api.CreateImage)...)
				_launch_screen.PUT("/image", append(_changeimageMw(), api.ChangeImage)...)
				_launch_screen.GET("/screen", append(_mobilegetimageMw(), api.MobileGetImage)...)
			}
			{
				_login0 := _v1.Group("/login", _login0Mw()...)
				_login0.GET("/access-token", append(_gettokenMw(), api.GetToken)...)
				_login0.GET("/refresh-token", append(_refreshtokenMw(), api.RefreshToken)...)
			}
			{
				_paper := _v1.Group("/paper", _paperMw()...)
				_paper.GET("/download", append(_getdownloadurlMw(), api.GetDownloadUrl)...)
				_paper.GET("/list", append(_listdirfilesMw(), api.ListDirFiles)...)
			}
			{
				_terms := _v1.Group("/terms", _termsMw()...)
				_terms.GET("/info", append(_gettermMw(), api.GetTerm)...)
				_terms.GET("/list", append(_gettermslistMw(), api.GetTermsList)...)
			}
			{
				_user1 := _v1.Group("/user", _user1Mw()...)
				_user1.POST("/validate-code", append(_validatecodeMw(), api.ValidateCode)...)
			}
		}
		{
			_v2 := _api.Group("/v2", _v2Mw()...)
			{
				_common0 := _v2.Group("/common", _common0Mw()...)
				_common0.GET("/fzu-helper.css", append(_getcssMw(), api.GetCSS)...)
				_common0.GET("/fzu-helper.html", append(_gethtmlMw(), api.GetHtml)...)
				_common0.GET("/user-agreement.html", append(_getuseragreementMw(), api.GetUserAgreement)...)
			}
			{
				_url := _v2.Group("/url", _urlMw()...)
				_url.GET("/beta.apk", append(_downloadbetaapkMw(), api.DownloadBetaApk)...)
				_url.GET("/dump", append(_getdumpMw(), api.GetDump)...)
				_url.GET("/getcloud", append(_getcloudMw(), api.GetCloud)...)
				_url.POST("/login", append(_login1Mw(), api.Login)...)
				_url.GET("/release.apk", append(_downloadreleaseapkMw(), api.DownloadReleaseApk)...)
				_url.POST("/setcloud", append(_setcloudMw(), api.SetCloud)...)
				_url.GET("/settings.php", append(_getsettingMw(), api.GetSetting)...)
				_url.POST("/test", append(_gettestMw(), api.GetTest)...)
				_url.POST("/upload", append(_uploadversionMw(), api.UploadVersion)...)
				_url.POST("/upload-params", append(_uploadparamsMw(), api.UploadParams)...)
				_url.GET("/version.json", append(_getreleaseversionMw(), api.GetReleaseVersion)...)
				_url.GET("/versionbeta.json", append(_getbetaversionMw(), api.GetBetaVersion)...)
			}
		}
	}
}
