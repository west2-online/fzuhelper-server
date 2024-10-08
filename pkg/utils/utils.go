package utils

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

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
