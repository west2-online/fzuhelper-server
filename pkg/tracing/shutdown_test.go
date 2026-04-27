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

package tracing

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ProviderShutdownWithContext(t *testing.T) {
	Convey("Test ProviderShutdownWithContext", t, func() {
		ctx := context.WithValue(context.Background(), "key", "value") //nolint
		called := false

		fn := ProviderShutdownWithContext(func(got context.Context) error {
			called = true
			So(got, ShouldEqual, ctx)
			return nil
		}, ctx, "shutdown failed: %v")

		fn()
		So(called, ShouldBeTrue)
	})
}

func Test_ProviderShutdown(t *testing.T) {
	Convey("Test ProviderShutdown", t, func() {
		called := false

		fn := ProviderShutdown(func(got context.Context) error {
			called = true
			So(got, ShouldNotBeNil)
			return nil
		}, "shutdown failed: %v")

		fn()
		So(called, ShouldBeTrue)
	})
}
