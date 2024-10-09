package pack

import (
	"fmt"
	"time"
)

func UploadImg(file []byte, name string) error {
	//body := bytes.NewReader(file)
	//url := config.Upcloud.DomainName + config.Upcloud.Path + name
	//
	//req, err := http.NewRequest("PUT", url, body)
	//if err != nil {
	//	return err
	//}
	//req.SetBasicAuth(config.Upcloud.User, config.Upcloud.Pass)
	//req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	//res, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	return err
	//}
	//defer res.Body.Close()
	//if res.StatusCode != 200 {
	//	return errno.UpcloudError
	//}
	return nil
}

func GenerateImgName(uid int64) string {
	currentTime := time.Now()
	// 获取年月日和小时分钟
	year, month, day := currentTime.Date()
	hour, minute := currentTime.Hour(), currentTime.Minute()
	second := currentTime.Second()
	return fmt.Sprintf("%v_%d%02d%02d_%02d%02d%02d_cover.jpg", uid, year, month, day, hour, minute, second)
}
