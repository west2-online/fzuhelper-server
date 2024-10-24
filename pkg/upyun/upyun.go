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

package upyun

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/upyun/go-sdk/v3/upyun"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var UpYun *upyun.UpYun

func NewUpYun() {
	UpYun = upyun.NewUpYun(
		&upyun.UpYunConfig{
			Bucket:   config.UpYun.Bucket,
			Operator: config.UpYun.Operator,
			Password: config.UpYun.Password,
		},
	)
}

/**
* GetDir() 获取目录下的文件和文件夹
* @param path string 目录路径
* @return *model.UpYunFileDir 目录路径以及该路径下的文件和文件夹
* @return error 错误信息
 */
func GetDir(path string) (*model.UpYunFileDir, error) {
	var err error
	fileDir := &model.UpYunFileDir{
		BasePath: path,
		Folders:  []string{},
		Files:    []string{},
	}
	objsChan := make(chan *upyun.FileInfo, 10)
	go func() {
		err = UpYun.List(&upyun.GetObjectsConfig{
			Path:        path,
			ObjectsChan: objsChan,
		})
		logger.Error(err)
	}()
	for obj := range objsChan {
		if obj.IsDir {
			if !obj.IsEmptyDir && !strings.HasPrefix(obj.Name, "__") { // 过滤空和临时文件夹
				fileDir.Folders = append(fileDir.Folders, obj.Name)
			}
		} else {
			fileDir.Files = append(fileDir.Files, obj.Name)
		}
	}
	return fileDir, err
}

// 拼接url的时候需要转义一下，防止特殊字符出现在url中导致莫名其妙的bug
func GetDownloadUrl(uri string) (string, error) {
	etime := strconv.FormatInt(time.Now().Unix()+config.UpYun.TokenTimeout, 10)
	sign := utils.MD5(strings.Join([]string{config.UpYun.TokenSecret, etime, uri}, "&"))
	url := fmt.Sprintf("%s%s?_upt=%s%s", config.UpYun.UssDomain, utils.UriEncode(uri), sign[12:20], etime)
	return url, nil
}

func UploadFile(filepath, ussDir string) error {
	return UpYun.Put(&upyun.PutObjectConfig{
		Path:      ussDir,
		LocalPath: filepath,
		UseMD5:    true,
	})
}
