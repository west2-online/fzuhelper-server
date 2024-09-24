package utils

import (
	"errors"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	config "github.com/west2-online/fzuhelper-server/config"
)

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

func IsVideoFile(header *multipart.FileHeader) bool {
	contentType := header.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "video/") {
		return true
	}

	filename := header.Filename
	extensions := []string{".mp4", ".avi", ".mkv", ".mov"} // Add more video extensions if needed
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}

	return false
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
