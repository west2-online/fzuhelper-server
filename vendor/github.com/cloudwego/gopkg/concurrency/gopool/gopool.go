/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gopool

import (
	"context"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"
)

// Option ...
type Option struct {
	// MaxIdleWorkers is the max idle workers keeping in pool for waiting tasks.
	// There workers will exit after `WorkerMaxAge`
	MaxIdleWorkers int

	// WorkerMaxAge is the max age of a worker in pool.
	WorkerMaxAge time.Duration

	// TaskChanBuffer is the size of task queue length.
	// if it's full, we will fall back to use `go` directly without using pool.
	// normally, the queue length should be small,
	// coz we will create new workers to pick tasks if necessary.
	TaskChanBuffer int
}

// DefaultOption returns the default values of Option.
func DefaultOption() *Option {
	return &Option{
		MaxIdleWorkers: 1000,
		WorkerMaxAge:   time.Minute,
		TaskChanBuffer: 1000,
	}
}

var defaultGoPool = NewGoPool("__default__", nil)

// Go runs the given func in background
func Go(f func()) {
	defaultGoPool.Go(f)
}

// CtxGo runs the given func in background, and it passes ctx to panic handler when happens.
func CtxGo(ctx context.Context, f func()) {
	defaultGoPool.CtxGo(ctx, f)
}

// SetPanicHandler sets a func for handling panic cases.
//
// check the comment of (*GoPool).SetPanicHandler for details
func SetPanicHandler(f func(ctx context.Context, r interface{})) {
	defaultGoPool.SetPanicHandler(f)
}

type task struct {
	ctx context.Context
	f   func()
}

// GoPool represents a simple worker pool which manages goroutines for background tasks.
type GoPool struct {
	name string

	workers int32
	maxIdle int32
	maxage  int64 // milliseconds

	panicHandler func(ctx context.Context, r interface{})

	tasks     chan task
	unixMilli int64

	createWorker func()
}

// NewGoPool create a new instance for goroutine worker
func NewGoPool(name string, o *Option) *GoPool {
	if o == nil {
		o = DefaultOption()
	}
	p := &GoPool{
		name:    name,
		tasks:   make(chan task, o.TaskChanBuffer),
		maxage:  o.WorkerMaxAge.Milliseconds(),
		maxIdle: int32(o.MaxIdleWorkers),
	}

	// fix: func literal escapes to heap
	p.createWorker = func() {
		p.runWorker()
	}
	return p
}

// Go runs the given func in background
func (p *GoPool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

// CtxGo runs the given func in background, and it passes ctx to panic handler when happens.
func (p *GoPool) CtxGo(ctx context.Context, f func()) {
	select {
	case p.tasks <- task{ctx: ctx, f: f}:
	default:
		// full? fall back to use go directly
		go p.runTask(ctx, f)
		return
	}
	// luckily ... it's true when there're many workers.
	if len(p.tasks) == 0 {
		return
	}
	// all worker is busy, create a new one
	go p.createWorker()
}

// SetPanicHandler sets a func for handling panic cases.
//
// Panic handler takes two args, `ctx` and `r`.
// `ctx` is the one provided when calling CtxGo, and `r` is returned by recover()
//
// By default, GoPool will use log.Printf to record the err and stack.
//
// It's recommended to set your own handler.
func (p *GoPool) SetPanicHandler(f func(ctx context.Context, r interface{})) {
	p.panicHandler = f
}

func (p *GoPool) runTask(ctx context.Context, f func()) {
	defer func(p *GoPool, ctx context.Context) {
		if r := recover(); r != nil {
			if p.panicHandler != nil {
				p.panicHandler(ctx, r)
			} else {
				log.Printf("GOPOOL: panic in pool: %s: %v: %s", p.name, r, debug.Stack())
			}
		}
	}(p, ctx)
	f()
}

func (p *GoPool) CurrentWorkers() int {
	return int(atomic.LoadInt32(&p.workers))
}

func (p *GoPool) runWorker() {
	id := atomic.AddInt32(&p.workers, 1)
	defer atomic.AddInt32(&p.workers, -1)

	if id > p.maxIdle {
		// drain task chan and exit without waiting
		for {
			select {
			case t := <-p.tasks:
				p.runTask(t.ctx, t.f)
			default:
				return
			}
		}
	}

	createdAt := time.Now().UnixMilli() // for checking maxage
	for t := range p.tasks {
		p.runTask(t.ctx, t.f)

		now := atomic.LoadInt64(&p.unixMilli)

		// check if ticker is NOT alive
		// p.unixMilli will be set to zero if it's not running
		if now == 0 {
			// cas and create a new ticker
			now = time.Now().UnixMilli()
			if atomic.CompareAndSwapInt64(&p.unixMilli, 0, now) {
				go p.runTicker()
			}
		}

		// check maxage
		if now-createdAt > p.maxage {
			return
		}
	}
}

// noopTask is used by runTicker() to wake up workers and checks their age.
var noopTask = task{f: func() {}}

func (p *GoPool) runTicker() {
	// mark it zero to trigger ticker to be created when we have active workers
	defer atomic.StoreInt64(&p.unixMilli, 0)

	// If p.maxage=1s, it updates `unixMilli` and sends 100 noop tasks per second.
	// As a result, workers may take longer time to exit, and this is expected.
	d := time.Duration(p.maxage) * time.Millisecond / 100

	// set a minimum value to avoid performance issues.
	if d < time.Millisecond {
		d = time.Millisecond
	}

	t := time.NewTicker(d)
	defer t.Stop()

	for now := range t.C {
		if p.CurrentWorkers() == 0 {
			return
		}
		atomic.StoreInt64(&p.unixMilli, now.UnixMilli())
		p.tasks <- noopTask
	}
}
