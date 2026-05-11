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
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroupDo(t *testing.T) {
	t.Parallel()

	var group Group[int]
	got, err := group.Do("key", func() (int, error) {
		return 1, nil
	})
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got != 1 {
		t.Fatalf("Do() = %d, want 1", got)
	}
}

func TestGroupDoReturnsLoaderError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("load failed")
	var group Group[int]
	got, err := group.Do("key", func() (int, error) {
		return 0, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Do() error = %v, want %v", err, wantErr)
	}
	if got != 0 {
		t.Fatalf("Do() = %d, want zero value", got)
	}
}

func TestGroupDoSharedCoalescesConcurrentCalls(t *testing.T) {
	t.Parallel()

	var (
		group Group[int]
		calls atomic.Int64
		wg    sync.WaitGroup
	)

	const goroutines = 8
	results := make(chan int, goroutines)
	sharedResults := make(chan bool, goroutines)
	errs := make(chan error, goroutines)

	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			got, shared, err := group.DoShared("key", func() (int, error) {
				calls.Add(1)
				time.Sleep(20 * time.Millisecond)
				return 7, nil
			})
			results <- got
			sharedResults <- shared
			errs <- err
		}()
	}

	wg.Wait()
	close(results)
	close(sharedResults)
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("DoShared() error = %v", err)
		}
	}
	for got := range results {
		if got != 7 {
			t.Fatalf("DoShared() = %d, want 7", got)
		}
	}
	for shared := range sharedResults {
		if !shared {
			t.Fatalf("DoShared() shared = false, want true")
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("loader calls = %d, want 1", calls.Load())
	}
}

func TestGroupForget(t *testing.T) {
	t.Parallel()

	var (
		group Group[int]
		calls atomic.Int64
	)

	load := func() (int, error) {
		return int(calls.Add(1)), nil
	}

	got, err := group.Do("key", load)
	if err != nil {
		t.Fatalf("first Do() error = %v", err)
	}
	if got != 1 {
		t.Fatalf("first Do() = %d, want 1", got)
	}

	group.Forget("key")

	got, err = group.Do("key", load)
	if err != nil {
		t.Fatalf("second Do() error = %v", err)
	}
	if got != 2 {
		t.Fatalf("second Do() = %d, want 2", got)
	}
}
