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

package singleflight

import (
	"errors"

	xsingleflight "golang.org/x/sync/singleflight"
)

var ErrInvalidType = errors.New("singleflight: invalid type assertion")

type Loader[T any] func() (T, error)

var group xsingleflight.Group

func Do[T any](key string, load Loader[T]) (T, error) {
	v, _, err := DoShared(key, load)
	return v, err
}

func DoShared[T any](key string, load Loader[T]) (T, bool, error) {
	// 把shared暴露出来，日志可能会需要
	v, err, shared := group.Do(key, func() (any, error) {
		return load()
	})
	if err != nil {
		var zero T
		return zero, shared, err
	}
	res, ok := v.(T)
	if !ok {
		var zero T
		return zero, shared, ErrInvalidType
	}
	return res, shared, nil
}

func Forget(key string) {
	group.Forget(key)
}
