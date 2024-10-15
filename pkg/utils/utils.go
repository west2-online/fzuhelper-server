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
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/errno"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/jwch"

	config "github.com/west2-online/fzuhelper-server/config"
)

func TimeParse(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func GetMysqlDSN() (string, error) {
	if config.Mysql == nil {
		return "", errors.New("config not found")
	}

	dsn := strings.Join([]string{config.Mysql.Username, ":", config.Mysql.Password, "@tcp(", config.Mysql.Addr, ")/", config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true"}, "")

	return dsn, nil
}

func GetMQUrl() (string, error) {
	if config.RabbitMQ == nil {
		return "", errors.New("config not found")
	}

	url := strings.Join([]string{"amqp://", config.RabbitMQ.Username, ":", config.RabbitMQ.Password, "@", config.RabbitMQ.Addr, "/"}, "")

	return url, nil
}

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

func GetAvailablePort() (string, error) {
	if config.Service.AddrList == nil {
		return "", errors.New("utils.GetAvailablePort: config.Service.AddrList is nil")
	}
	for _, addr := range config.Service.AddrList {
		if ok := AddrCheck(addr); ok {
			return addr, nil
		}
	}
	return "", errors.New("utils.GetAvailablePort: not available port from config")
}

func RetryLogin(stu *jwch.Student) error {
	var err error
	delay := constants.InitialDelay

	for attempt := 1; attempt <= constants.MaxRetries; attempt++ {
		err = stu.Login()
		if err == nil {
			return nil // 登录成功
		}

		if attempt < constants.MaxRetries {
			time.Sleep(delay) // 等待一段时间后再重试
			delay *= 2        // 指数退避，逐渐增加等待时间
		}
	}

	return fmt.Errorf("failed to login after %d attempts: %w", constants.MaxRetries, err)
}

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

func LoadCNLocation() *time.Location {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	return Loc
}

func GenerateRedisKeyByStuId(stuId int64, sType int64) string {
	return strings.Join([]string{strconv.FormatInt(stuId, 10), strconv.FormatInt(sType, 10)}, ":")
}
