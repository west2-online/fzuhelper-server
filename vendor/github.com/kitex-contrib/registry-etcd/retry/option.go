// Copyright 2021 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import "time"

// Option is the only struct that can be used to set Retry Config.
type Option struct {
	F func(o *Config)
}

// WithMaxAttemptTimes sets MaxAttemptTimes
func WithMaxAttemptTimes(maxAttemptTimes uint) Option {
	return Option{F: func(o *Config) {
		o.MaxAttemptTimes = maxAttemptTimes
	}}
}

// WithObserveDelay sets ObserveDelay
func WithObserveDelay(observeDelay time.Duration) Option {
	return Option{F: func(o *Config) {
		o.ObserveDelay = observeDelay
	}}
}

// WithRetryDelay sets RetryDelay
func WithRetryDelay(retryDelay time.Duration) Option {
	return Option{F: func(o *Config) {
		o.RetryDelay = retryDelay
	}}
}
