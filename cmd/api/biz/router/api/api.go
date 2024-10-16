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

	api "github.com/west2-online/fzuhelper-server/cmd/api/biz/handler/api"
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
					_user := _jwch.Group("/user", _userMw()...)
					_user.GET("/login", append(_getlogindataMw(), api.GetLoginData)...)
					_user.GET("/validateCode", append(_getvalidatecodeMw(), api.GetValidateCode)...)
				}
			}
		}
	}
	{
		_launch_screen := root.Group("/launch_screen", _launch_screenMw()...)
		{
			_api0 := _launch_screen.Group("/api", _api0Mw()...)
			_api0.DELETE("/image", append(_deleteimageMw(), api.DeleteImage)...)
			_image := _api0.Group("/image", _imageMw()...)
			_image.GET("/point", append(_addimagepointtimeMw(), api.AddImagePointTime)...)
			_api0.GET("/image", append(_getimageMw(), api.GetImage)...)
			_image0 := _api0.Group("/image", _image0Mw()...)
			_image0.PUT("/img", append(_changeimageMw(), api.ChangeImage)...)
			_api0.POST("/image", append(_createimageMw(), api.CreateImage)...)
			_api0.PUT("/image", append(_changeimagepropertyMw(), api.ChangeImageProperty)...)
			_api0.GET("/images", append(_getimagesbyuseridMw(), api.GetImagesByUserId)...)
			_api0.GET("/screen", append(_mobilegetimageMw(), api.MobileGetImage)...)
		}
	}
}
