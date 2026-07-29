package constants

import "time"

const (
	SignedLocationTimeout = 15 * time.Second
	// LocationServiceSuccessCode 参考: https://github.com/west2-online/location/blob/master/pkg/base/code.go
	LocationServiceSuccessCode = 10000
)
