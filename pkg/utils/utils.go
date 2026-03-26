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

const (
	HoursInADay            = 24
	DaysInAWeek            = 7
	DefaultFilePermissions = 0o666 // 默认文件权限
)

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
	defer func() {
		if err := l.Close(); err != nil {
			logger.Errorf("utils.AddrCheck: failed to close listener: %v", err.Error())
		}
	}()
	return true
}

// ParseCookies 将cookie字符串解析为 http.Cookie
// 这里只能解析这样的数组： "Key=Value; Key=Value"
func ParseCookies(rawData string) []*http.Cookie {
	var cookies []*http.Cookie
	maxSplitNumber := 2

	// 按照分号分割每个 Cookie
	pairs := strings.Split(rawData, ";")
	for _, pair := range pairs {
		// 去除空格并检查是否为空
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// 按等号分割键和值
		parts := strings.SplitN(pair, "=", maxSplitNumber)
		if len(parts) != maxSplitNumber {
			continue // 如果格式不正确，跳过
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 创建 http.Cookie 并添加到切片
		cookie := &http.Cookie{
			Name:  key,
			Value: value,
		}
		cookies = append(cookies, cookie)
	}

	return cookies
}

// ParseCookiesToString 会尝试解析 cookies 到 string
// 只会返回 "Key=Value; Key=Value"样式的文本数组
func ParseCookiesToString(cookies []*http.Cookie) string {
	var cookieStrings []string
	for _, cookie := range cookies {
		cookieStrings = append(cookieStrings, cookie.Name+"="+cookie.Value)
	}
	return strings.Join(cookieStrings, "; ")
}

// GetAvailablePort 会尝试获取可用的监听地址
// NOTICE: 调用这个函数前需要确保 config 已经初始化
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

// FileToByteArray 用于将客户端发来的文件转换为[][]byte格式，用于流式传输
func FileToByteArray(file *multipart.FileHeader) (fileBuf [][]byte, err error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, errno.ParamError
	}
	defer func() {
		// 捕获并处理关闭文件时可能发生的错误
		if err := fileContent.Close(); err != nil {
			logger.Errorf("utils.FileToByteArray: failed to close file: %v", err.Error())
		}
	}()

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
	defer func() {
		// 捕获并处理关闭文件时可能发生的错误
		if err := file.Close(); err != nil {
			logger.Errorf("utils.CheckImageFileType: failed to close file: %v", err.Error())
		}
	}()

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

// GetWeekdayByDate 根据学期开始日期和目标日期计算出当前是第几周的第几天
func GetWeekdayByDate(termStartDate string, date string) (week int, day int, err error) {
	startDate, err := TimeParse(termStartDate)
	if err != nil {
		return 0, 0, fmt.Errorf("utils.GetWeekdayByDate: Parse term start date fail: %w", err)
	}

	targetDate, err := TimeParse(date)
	if err != nil {
		return 0, 0, fmt.Errorf("utils.GetWeekdayByDate: Parse date fail: %w", err)
	}

	// 如果不是周一，向前调整到周一
	offset := (int(startDate.Weekday()) - 1 + DaysInAWeek) % DaysInAWeek
	startDate = startDate.AddDate(0, 0, -offset)
	diff := targetDate.Sub(startDate)

	if diff < 0 {
		return 0, 0, fmt.Errorf("utils.GetWeekdayByDate: Target date is before term start date")
	}

	week = int(diff.Hours()/HoursInADay)/DaysInAWeek + 1
	day = int(diff.Hours()/HoursInADay)%DaysInAWeek + 1
	return week, day, nil
}
