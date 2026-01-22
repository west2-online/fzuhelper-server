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

package handler

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		expectStatus   int
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "ping success",
			url:            "/ping",
			expectStatus:   consts.StatusOK,
			expectContains: "pong",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/ping", Ping)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, tc.expectStatus, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
