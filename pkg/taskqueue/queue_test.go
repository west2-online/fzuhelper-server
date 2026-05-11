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

package taskqueue

import (
	"context"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func Test_executeTask(t *testing.T) {
	Convey("Test btq.executeTask", t, func() {
		Convey("should record tracing metadata", func() {
			recorder := tracetest.NewSpanRecorder()
			tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
			originalProvider := otel.GetTracerProvider()
			otel.SetTracerProvider(tp)
			defer otel.SetTracerProvider(originalProvider)

			const (
				traceKeyValue = "trace-key"
				taskTypeValue = "ScheduleQueueTask"
			)
			btq := NewBaseTaskQueue()
			err := btq.executeTask(context.Background(), traceKeyValue, func(ctx context.Context) error {
				So(oteltrace.SpanFromContext(ctx).IsRecording(), ShouldBeTrue)
				return nil
			}, taskTypeValue)
			So(err, ShouldBeNil)

			spans := recorder.Ended()
			So(len(spans), ShouldEqual, 1)

			attrs := spans[0].Attributes()
			assertHasAttribute(t, attrs, tqKeyKey, traceKeyValue)
			assertHasAttribute(t, attrs, tqTypeKey, taskTypeValue)
		})

		Convey("should mark span as error if failed", func() {
			recorder := tracetest.NewSpanRecorder()
			tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
			originalProvider := otel.GetTracerProvider()
			otel.SetTracerProvider(tp)
			defer otel.SetTracerProvider(originalProvider)

			btq := NewBaseTaskQueue()
			wantErr := errors.New("boom")

			err := btq.executeTask(context.Background(), "error-key", func(context.Context) error {
				return wantErr
			}, "ScheduleQueueTask")
			So(errors.Is(err, wantErr), ShouldBeTrue)

			spans := recorder.Ended()
			So(len(spans), ShouldEqual, 1)
			So(spans[0].Status().Code, ShouldEqual, codes.Error)
		})
	})
}

func assertHasAttribute(t *testing.T, attrs []attribute.KeyValue, key, value string) {
	t.Helper()

	for _, attr := range attrs {
		if string(attr.Key) == key && attr.Value == attribute.StringValue(value) {
			return
		}
	}

	So(false, ShouldBeTrue)
}
