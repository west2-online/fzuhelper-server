package utils

import (
	"errors"
	"io"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/west2-online/fzuhelper-server/config"
)

func UploadImg(body io.Reader) (string, error) {
	name := uuid.NewV4().String()
	url := "http://v0.api.upyun.com/fzuhelper-assets" + config.USS.Path + name

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.USS.User, config.USS.Pass)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New(res.Status)
	}
	return config.USS.DomainName + config.USS.Path + name, nil
}
