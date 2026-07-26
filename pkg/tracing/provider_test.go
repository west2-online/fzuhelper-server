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
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Test_getResource(t *testing.T) {
	Convey("Test getResource", t, func() {
		Convey("should include serviceName in attributes", func() {
			const serviceName = "trace-test-service"

			res := getResource(context.Background(), serviceName)
			So(res, ShouldNotBeNil)

			found := false
			for _, attr := range res.Attributes() {
				if attr.Key == semconv.ServiceNameKey && attr.Value == attribute.StringValue(serviceName) {
					found = true
					break
				}
			}

			So(found, ShouldBeTrue)
		})
	})
}
