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

type Config struct {
	// The maximum number of call attempt times, including the initial call
	MaxAttemptTimes uint

	// The delay time of observing etcd key
	ObserveDelay time.Duration

	// The retry delay time
	RetryDelay time.Duration
}

func (o *Config) Apply(opts []Option) {
	for _, op := range opts {
		op.F(o)
	}
}

func NewRetryConfig(opts ...Option) *Config {
	retryConfig := &Config{
		MaxAttemptTimes: 5,
		ObserveDelay:    30 * time.Second,
		RetryDelay:      10 * time.Second,
	}

	retryConfig.Apply(opts)

	return retryConfig
}
