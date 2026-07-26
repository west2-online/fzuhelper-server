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

func TestDo(t *testing.T) {
	t.Parallel()

	got, err := Do(t.Name(), func() (int, error) {
		return 1, nil
	})
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got != 1 {
		t.Fatalf("Do() = %d, want 1", got)
	}
}

func TestDoReturnsLoaderError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("load failed")
	got, err := Do(t.Name(), func() (int, error) {
		return 0, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Do() error = %v, want %v", err, wantErr)
	}
	if got != 0 {
		t.Fatalf("Do() = %d, want zero value", got)
	}
}

func TestDoSharedNotShared(t *testing.T) {
	t.Parallel()

	got, shared, err := DoShared(t.Name(), func() (int, error) {
		return 3, nil
	})
	if err != nil {
		t.Fatalf("DoShared() error = %v", err)
	}
	if got != 3 {
		t.Fatalf("DoShared() = %d, want 3", got)
	}
	if shared {
		t.Fatalf("DoShared() shared = true, want false")
	}
}

func TestDoSharedCoalescesConcurrentCalls(t *testing.T) {
	t.Parallel()

	var (
		calls atomic.Int64
		wg    sync.WaitGroup
	)

	const goroutines = 8
	key := t.Name()
	results := make(chan int, goroutines)
	sharedResults := make(chan bool, goroutines)
	errs := make(chan error, goroutines)

	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			got, shared, err := DoShared(key, func() (int, error) {
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

func TestDoSharedReturnsNilValue(t *testing.T) {
	t.Parallel()

	got, shared, err := DoShared(t.Name(), func() (*int, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("DoShared() error = %v", err)
	}
	if got != nil {
		t.Fatalf("DoShared() = %v, want nil", got)
	}
	if shared {
		t.Fatalf("DoShared() shared = true, want false")
	}
}

func TestForget(t *testing.T) {
	t.Parallel()

	var calls atomic.Int64
	key := t.Name()

	load := func() (int, error) {
		return int(calls.Add(1)), nil
	}

	got, err := Do(key, load)
	if err != nil {
		t.Fatalf("first Do() error = %v", err)
	}
	if got != 1 {
		t.Fatalf("first Do() = %d, want 1", got)
	}

	Forget(key)

	got, err = Do(key, load)
	if err != nil {
		t.Fatalf("second Do() error = %v", err)
	}
	if got != 2 {
		t.Fatalf("second Do() = %d, want 2", got)
	}
}

func TestDoSharedInvalidType(t *testing.T) {
	t.Parallel()

	key := t.Name()
	started := make(chan struct{})
	release := make(chan struct{})
	errs := make(chan error, 1)
	result := make(chan struct {
		got    string
		shared bool
		err    error
	}, 1)

	go func() {
		_, _, err := DoShared(key, func() (int, error) {
			close(started)
			<-release
			return 1, nil
		})
		errs <- err
	}()

	<-started
	go func() {
		got, shared, err := DoShared(key, func() (string, error) {
			return "unexpected", nil
		})
		result <- struct {
			got    string
			shared bool
			err    error
		}{got: got, shared: shared, err: err}
	}()

	time.Sleep(10 * time.Millisecond)
	close(release)
	second := <-result

	if !errors.Is(second.err, ErrInvalidType) {
		t.Fatalf("DoShared() error = %v, want %v", second.err, ErrInvalidType)
	}
	if second.got != "" {
		t.Fatalf("DoShared() = %q, want empty string", second.got)
	}
	if !second.shared {
		t.Fatalf("DoShared() shared = false, want true")
	}
	if err := <-errs; err != nil {
		t.Fatalf("first DoShared() error = %v", err)
	}
}
