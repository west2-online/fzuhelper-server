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
	"fmt"
	"io/ioutil"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
	"github.com/west2-online/fzuhelper-server/cmd/paper/dal/db"
	"github.com/west2-online/fzuhelper-server/cmd/paper/dal/upyun"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *PaperService) UploadFile(req *paper.UploadFileRequest) error {

	fileUid := uuid.Must(uuid.NewV4(), nil).String()

	savePath := path.Join(constants.CacheDst, fileUid)

	err := os.MkdirAll(constants.CacheDst, os.ModePerm)
	if err != nil {
		return fmt.Errorf("service.UploadFile failed to create cache dir: %w", err)
	}

	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("service.UploadFile failed to create file: %w", err)
	}
	defer out.Close()

	_, err = out.Write(req.Content)
	if err != nil {
		return fmt.Errorf("service.UploadFile failed to write file: %w", err)
	}

	db.AddCheck(req.Filename, req.UssPath, req.User, fileUid)

	go func() {
		files, err := ioutil.ReadDir(constants.CacheDst)
		if err != nil {
			logger.Errorf("service.UploadFile failed to read cache dir: %v", err)
			return
		}

		for _, file := range files {
			filename := file.Name()
			filepath := path.Join(constants.CacheDst, filename)
			ussPath := path.Join(config.UpYun.UnCheckedDir, filename)

			err = upyun.UploadFile(filepath, ussPath)
			if err != nil {
				logger.Errorf("service.UploadFile failed to upload file: %v", err)
				continue
			}

			err = os.Remove(filepath)
			if err != nil {
				logger.Errorf("service.UploadFile failed to remove file: %v", err)

			}
		}
	}()

	return nil
}
