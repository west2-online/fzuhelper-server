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
			_v1 := _api.Group("/v1", _v1Mw()...)
			{
				_common := _v1.Group("/common", _commonMw()...)
				{
					_classroom := _common.Group("/classroom", _classroomMw()...)
					_classroom.GET("/empty", append(_getemptyclassroomsMw(), api.GetEmptyClassrooms)...)
				}
			}
			{
				_jwch := _v1.Group("/jwch", _jwchMw()...)
				{
					_academic := _jwch.Group("/academic", _academicMw()...)
					_academic.GET("/credit", append(_getcreditMw(), api.GetCredit)...)
					_academic.GET("/gpa", append(_getgpaMw(), api.GetGPA)...)
					_academic.GET("/scores", append(_getscoresMw(), api.GetScores)...)
					_academic.GET("/unifiedExam", append(_getunifiedexamMw(), api.GetUnifiedExam)...)
				}
				{
					_course := _jwch.Group("/course", _courseMw()...)
					_course.GET("/list", append(_getcourselistMw(), api.GetCourseList)...)
				}
				{
					_user := _jwch.Group("/user", _userMw()...)
					_user.GET("/login", append(_getlogindataMw(), api.GetLoginData)...)
					_user.POST("/validateCode", append(_validatecodeMw(), api.ValidateCode)...)
				}
			}
			{
				_paper := _v1.Group("/paper", _paperMw()...)
				_paper.GET("/download", append(_getdownloadurlMw(), api.GetDownloadUrl)...)
				_paper.GET("/list", append(_listdirfilesMw(), api.ListDirFiles)...)
			}
			{
				_url := _v1.Group("/url", _urlMw()...)
				_url.GET("/beta.apk", append(_getdownloadbetaMw(), api.GetDownloadBeta)...)
				_url.GET("/dump", append(_dumpvisitMw(), api.DumpVisit)...)
				_url.POST("/login", append(_apiloginMw(), api.APILogin)...)
				_url.GET("/release.apk", append(_getdownloadreleaseMw(), api.GetDownloadRelease)...)
				_url.GET("/settings.php", append(_getcloudsettingMw(), api.GetCloudSetting)...)
				_url.GET("/version.json", append(_getreleaseversionMw(), api.GetReleaseVersion)...)
				_url.GET("/versionbeta.json", append(_getbetaversionMw(), api.GetBetaVersion)...)
				{
					_api0 := _url.Group("/api", _api0Mw()...)
					_api0.GET("/getcloud", append(_getallcloudsettingMw(), api.GetAllCloudSetting)...)
					_api0.POST("/setcloud", append(_setallcloudsettingMw(), api.SetAllCloudSetting)...)
					_api0.POST("/test", append(_testsettingMw(), api.TestSetting)...)
					_api0.POST("/upload", append(_uploadversioninfoMw(), api.UploadVersionInfo)...)
					_api0.POST("/uploadparams", append(_getuploadparamsMw(), api.GetUploadParams)...)
				}
				{
					_onekey := _url.Group("/onekey", _onekeyMw()...)
					_onekey.GET("/FZUHelper.css", append(_fzuhelpercssMw(), api.FZUHelperCSS)...)
					_onekey.GET("/FZUHelper.html", append(_fzuhelperhtmlMw(), api.FZUHelperHTML)...)
					_onekey.GET("/UserAgreement.html", append(_useragreementhtmlMw(), api.UserAgreementHTML)...)
				}
			}
		}
	}
	{
		_launch_screen := root.Group("/launch_screen", _launch_screenMw()...)
		{
			_api1 := _launch_screen.Group("/api", _api1Mw()...)
			_api1.DELETE("/image", append(_deleteimageMw(), api.DeleteImage)...)
			_image := _api1.Group("/image", _imageMw()...)
			_image.GET("/point", append(_addimagepointtimeMw(), api.AddImagePointTime)...)
			_api1.GET("/image", append(_getimageMw(), api.GetImage)...)
			_image0 := _api1.Group("/image", _image0Mw()...)
			_image0.PUT("/img", append(_changeimageMw(), api.ChangeImage)...)
			_api1.POST("/image", append(_createimageMw(), api.CreateImage)...)
			_api1.PUT("/image", append(_changeimagepropertyMw(), api.ChangeImageProperty)...)
			_api1.GET("/screen", append(_mobilegetimageMw(), api.MobileGetImage)...)
		}
	}
}
