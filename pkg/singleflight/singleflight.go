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

	"golang.org/x/sync/singleflight"
)

var ErrInvalidType = errors.New("singleflight: invalid type assertion")

type Loader func() (any, error)

type Group struct {
	group singleflight.Group
}

func (g *Group) Do(key string, load Loader) (any, error) {
	// 由于go不支持给方法单独添加泛型，只能返回any了
	v, _, err := g.DoShared(key, load)
	return v, err
}

func (g *Group) DoShared(key string, load Loader) (any, bool, error) {
	// 把shared暴露出来，日志可能会需要
	v, err, shared := g.group.Do(key, load)
	return v, shared, err
}

func (g *Group) Forget(key string) {
	g.group.Forget(key)
}
