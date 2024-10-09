package pack

import (
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"io"
	"mime/multipart"
	"path/filepath"
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
