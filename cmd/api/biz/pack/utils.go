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

package pack

import (
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func FileToByte(file *multipart.FileHeader) ([]byte, error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, errno.ParamError
	}
	return io.ReadAll(fileContent)
}

func IsAllowImageExt(fileName string) bool {
	imageExt := filepath.Ext(fileName)
	allowExtImage := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
	}
	if _, ok := allowExtImage[imageExt]; !ok {
		return false
	}
	return true
}
