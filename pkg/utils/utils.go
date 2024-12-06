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

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const DefaultFilePermissions = 0o666 // 默认文件权限

// TimeParse 会将文本日期解析为标准时间对象
func TimeParse(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

// LoadCNLocation 载入cn时间
func LoadCNLocation() *time.Location {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	return Loc
}

// GetMysqlDSN 会拼接 mysql 的 DSN
func GetMysqlDSN() (string, error) {
	if config.Mysql == nil {
		return "", errors.New("config not found")
	}

	dsn := strings.Join([]string{
		config.Mysql.Username, ":", config.Mysql.Password,
		"@tcp(", config.Mysql.Addr, ")/",
		config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true",
	}, "")

	return dsn, nil
}

// GetEsHost 会获取 ElasticSearch 的客户端
func GetEsHost() (string, error) {
	if config.Elasticsearch == nil {
		return "", errors.New("elasticsearch not found")
	}

	return config.Elasticsearch.Host, nil
}

// AddrCheck 会检查当前的监听地址是否已被占用
func AddrCheck(addr string) bool {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	l.Close()
	return true
}

// ParseCookies 将cookie字符串解析为http.Cookie
func ParseCookies(rawData []string) []*http.Cookie {
	var cookies []*http.Cookie
	for _, raw := range rawData {
		cookie := &http.Cookie{}
		parts := strings.Split(raw, ";")
		for i, part := range parts {
			if i == 0 {
				cookieParts := strings.Split(part, "=")
				cookie.Name = cookieParts[0]
				cookie.Value = cookieParts[1]
			} else {
				cookieParts := strings.Split(part, "=")
				switch cookieParts[0] {
				case "Domain":
					cookie.Domain = cookieParts[1]
				case "Path":
					cookie.Path = cookieParts[1]
				case "Expires":
					if expires, err := time.Parse(time.RFC1123, cookieParts[1]); err == nil {
						cookie.Expires = expires
					}
				case "Max-Age":
					if maxAge, err := strconv.Atoi(cookieParts[1]); err == nil {
						cookie.MaxAge = maxAge
					}
				case "Secure":
					cookie.Secure = true
				case "HttpOnly":
					cookie.HttpOnly = true
				}
			}
		}
		cookies = append(cookies, cookie)
	}
	return cookies
}

// ParseCookiesToString 会尝试解析 cookies 到 string
func ParseCookiesToString(cookies []*http.Cookie) []string {
	var cookieStrings []string
	for _, cookie := range cookies {
		var parts []string
		parts = append(parts, cookie.Name+"="+cookie.Value)
		if cookie.Domain != "" {
			parts = append(parts, "Domain="+cookie.Domain)
		}
		if cookie.Path != "" {
			parts = append(parts, "Path="+cookie.Path)
		}
		if !cookie.Expires.IsZero() {
			parts = append(parts, "Expires="+cookie.Expires.Format(time.RFC1123))
		}
		if cookie.MaxAge > 0 {
			parts = append(parts, "Max-Age="+strconv.Itoa(cookie.MaxAge))
		}
		if cookie.Secure {
			parts = append(parts, "Secure")
		}
		if cookie.HttpOnly {
			parts = append(parts, "HttpOnly")
		}
		cookieStrings = append(cookieStrings, strings.Join(parts, "; "))
	}
	return cookieStrings
}

// GetAvailablePort 会尝试获取可用的监听地址
func GetAvailablePort() (string, error) {
	if config.Service.AddrList == nil {
		return "", errors.New("utils.GetAvailablePort: config.Service.AddrList is nil")
	}
	logger.Debugf("Available AddrList: %v", config.Service.AddrList)
	for _, addr := range config.Service.AddrList {
		if ok := AddrCheck(addr); ok {
			logger.Debugf("Finally Choose to listen: %v", addr)
			return addr, nil
		}
	}
	return "", errors.New("utils.GetAvailablePort: not available port from config")
}

// FileToByteArray 用于将客户端发来的文件转换为[][]byte格式，用于流式传输
func FileToByteArray(file *multipart.FileHeader) (fileBuf [][]byte, err error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, errno.ParamError
	}
	defer fileContent.Close()
	for {
		buf := make([]byte, constants.StreamBufferSize)
		_, err = fileContent.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errno.InternalServiceError
		}
		fileBuf = append(fileBuf, buf)
	}
	return fileBuf, nil
}

// CheckImageFileType 检查文件格式是否合规
func CheckImageFileType(header *multipart.FileHeader) (string, bool) {
	file, err := header.Open()
	if err != nil {
		return "", false
	}
	defer file.Close()

	buffer := make([]byte, constants.CheckFileTypeBufferSize)
	_, err = file.Read(buffer)
	if err != nil {
		return "", false
	}

	kind, _ := filetype.Match(buffer)

	// 检查是否为jpg、png
	switch kind {
	case types.Get("jpg"):
		return "jpg", true
	case types.Get("png"):
		return "png", true
	default:
		return "", false
	}
}

// GetImageFileType 获得图片格式
func GetImageFileType(fileBytes *[]byte) (string, error) {
	buffer := (*fileBytes)[:constants.CheckFileTypeBufferSize]

	kind, _ := filetype.Match(buffer)

	// 检查是否为jpg、png
	switch kind {
	case types.Get("jpg"):
		return "jpg", nil
	case types.Get("png"):
		return "png", nil
	default:
		return "", errno.InternalServiceError
	}
}

// GenerateRedisKeyByStuId 开屏页通过学号与sType与device生成缓存对应Key
func GenerateRedisKeyByStuId(stuId string, sType int64, device string) string {
	return strings.Join([]string{stuId, device, strconv.FormatInt(sType, 10)}, ":")
}

// SaveImageFromBytes 仅用于测试流式传输结果是否正确
func SaveImageFromBytes(imgBytes []byte, format string) error {
	// 使用 bytes.NewReader 将 []byte 转换为 io.Reader
	imgReader := bytes.NewReader(imgBytes)

	// 解码图片，自动检测图片格式（jpeg, png 等）
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return fmt.Errorf("can't decode img: %w", err)
	}

	// 创建保存图片的文件
	outFile, err := os.OpenFile("testImg.jpg", os.O_CREATE|os.O_WRONLY, DefaultFilePermissions)
	if err != nil {
		return fmt.Errorf("can't create img file: %w", err)
	}
	defer outFile.Close()

	// 根据格式保存图片
	switch format {
	case "jpeg", "jpg":
		// 将图片保存为 JPEG 格式
		err = jpeg.Encode(outFile, img, nil)
	case "png":
		// 将图片保存为 PNG 格式
		err = png.Encode(outFile, img)
	default:
		return fmt.Errorf("unsupport img type: %v", format)
	}

	if err != nil {
		return fmt.Errorf("save img failed: %w", err)
	}

	return nil
}

// SaveJSON 保存 JSON 数据到文件
func SaveJSON(fileName string, saveJson []byte) error {
	// 写入操作，可以覆盖文件
	if err := ioutil.WriteFile(fileName, saveJson, DefaultFilePermissions); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// GetJSON 从 JSON 文件读取数据
func GetJSON(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", fileName, err)
	}
	defer file.Close()

	// 读取文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	return data, nil
}
