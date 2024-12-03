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

// Code generated by hertz generator.

package api

import (
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/api/middleware"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _apiMw() []app.HandlerFunc {
	return nil
}

func _v1Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commonMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _classroomMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getemptyclassroomsMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _jwchMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.JwtMiddleware.MiddlewareFunc(),
	}
}

func _courseMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getcourselistMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.GetHeaderParams(),
	}
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getlogindataMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _paperMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getdownloadurlMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _listdirfilesMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _launch_screenMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _api0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _imageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _deleteimageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _addimagepointtimeMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _image0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getimageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _changeimageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _createimageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _changeimagepropertyMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _mobilegetimageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _academicMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.GetHeaderParams(),
	}
}

func _getcreditMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getgpaMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getscoresMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getunifiedexamMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _classroom0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getexamroominfoMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _urlMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getdownloadbetaMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _dumpvisitMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _apiloginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getdownloadreleaseMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getcloudsettingMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getreleaseversionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getbetaversionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getallcloudsettingMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _setallcloudsettingMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _testsettingMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _uploadversioninfoMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getuploadparamsMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _onekeyMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _fzuhelpercssMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _fzuhelperhtmlMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _useragreementhtmlMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _api1Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _validatecodeMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _validatecodeforandroidMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getdownloadurlforandroidMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _listdirfilesforandroidMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _gettokenMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _refreshtokenMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _login0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _testauthMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _internalMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _user0Mw() []app.HandlerFunc {
	// your code...
	return nil
}
