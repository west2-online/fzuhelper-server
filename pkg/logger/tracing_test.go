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

package logger

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func Test_extractSpanContext(t *testing.T) {
	Convey("Test extractSpanContext", t, func() {
		Convey("nil context", func() {
			fields := extractSpanContext(nil) //nolint
			So(fields, ShouldBeNil)
		})

		Convey("invalid span context", func() {
			fields := extractSpanContext(context.Background())
			So(fields, ShouldBeNil)
		})

		Convey("valid span context", func() {
			tp := sdktrace.NewTracerProvider()
			ctx, span := tp.Tracer("logger-context-test").Start(context.Background(), "test-span")
			defer span.End()

			fields := extractSpanContext(ctx)
			So(len(fields), ShouldEqual, 2)
			So(fields[0].Key, ShouldEqual, "trace_id")
			So(fields[0].String, ShouldNotBeBlank)
			So(fields[1].Key, ShouldEqual, "span_id")
			So(fields[1].String, ShouldNotBeBlank)
		})
	})
}

func Test_markSpanError(t *testing.T) {
	Convey("Test markSpanError", t, func() {
		recorder := tracetest.NewSpanRecorder()
		tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
		ctx, span := tp.Tracer("logger-context-test").Start(context.Background(), "test-span")

		markSpanError(ctx, "boom")
		span.End()

		spans := recorder.Ended()
		So(len(spans), ShouldEqual, 1)
		So(spans[0].Status().Code, ShouldEqual, codes.Error)
		So(len(spans[0].Events()), ShouldBeGreaterThan, 0)
	})
}
