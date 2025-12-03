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

package oss

type LaunchScreenOSSRepo interface {
	// UploadImg 又拍云上传文件
	UploadImg(file []byte, url string) error
	// DeleteImg 又拍云删除文件
	DeleteImg(url string) error
	// GenerateImgName 生成图片名字
	GenerateImgName(suffix string) (string, string, error)
	// GetRemotePathFromUrl 获得远程path
	GetRemotePathFromUrl(url string) string
}
