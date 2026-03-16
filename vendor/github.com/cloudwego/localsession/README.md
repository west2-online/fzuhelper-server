# LocalSession

## Introduction
LocalSession is used to **implicitly** manage and transmit context **within** or **between** goroutines. In canonical way, Go recommands developers to explicitly pass `context.Context` between functions to ensure the downstream callee get desired information from upstream. However this is tedious and ineffecient, resulting in many developers forget (or just don't want) to follow this practice. We have found many cases like that, especially in framework. Therefore, we design and implement a way to implicitly pass application context from root caller to end callee, without troubling intermediate implementation to always bring context.

## Usage
### Session
Session is an interface to carry and transmit your context. It has `Get()` and `WithValue()` methods to manipulate your data. And `IsValid()` method to tells your if it is valid at present. We provides two implementations by default:
- `SessionCtx`: use std `context.Context` as underlying storage, which means data from different goroutines are isolated.
- `SessionMap`: use std `map` as underlying storage, which means data from different goroutines are shared.

Both implementations are **Concurrent Safe**.

### SessionManager
SessionManager is a global manager of sessions. Through `BindSession()` and `CurSession()` methods it provides, you can transmit your session within the thread implicitly, without using explicit codes like `CallXXX(context.Context, args....)`.
```go
import (
	"context"
    "github.com/cloudwego/localsession"
)

// global manager
var manager = localsession.NewSessionManager(ManagerOptions{
	ShardNumber: 10,
	EnableImplicitlyTransmitAsync: true,
	GCInterval: time.Hour,
})

// global data
var key, v = "a", "b"
var key2, v2 = "c", "d"

func ASSERT(v bool) {
	if !v {
		panic("not true!")
	}
}

func main() {
    // get or initialize your context
    var ctx = context.Background()
    ctx = context.WithValue(ctx, key, v)

	// initialize new session with context
	var session = localsession.NewSessionCtx(ctx) 

	// set specific key-value and update session
	start := session.WithValue(key2, v2)

	// set current session
	manager.BindSession(start)

    // do somethings...
    
    // no need to pass context!
    GetDataX()
}

// read specific key under current session
func GetDataX() {
    // val exists
	val := manager.GetCurSession().Get(key) 
	ASSERT(val == v)

    // val2 exists
	val2 := manager.GetCurSession().Get(key2) 
	ASSERT(val2 == v2)
}
```

We provide a globally default manager to manage session between different goroutines, as long as you set `InitDefaultManager()` first.

### Explicitly Transmit Async Context (Recommended)
You can use `Go()` or `GoSession()` to explicitly transmit your context to other goroutines.

```go

package main

import (
	"context"
    . "github.com/cloudwego/localsession"
)

func init() {
    // initialize default manager first
	InitDefaultManager(DefaultManagerOptions())
}

func GetCurSession() Session {
	s, ok := CurSession()
	if !ok {
		panic("can't get current seession!")
	}
	return s
}

func main() {
	var ctx = context.Background()
	var key, v = "a", "b"
	var key2, v2 = "c", "d"
	var sig = make(chan struct{})
	var sig2 = make(chan struct{})

	// initialize new session with context
	var session = NewSessionCtx(ctx) // implementation...

	// set specific key-value and update session
	start := session.WithValue(key, v)

	// set current session
	BindSession(start)

	// pass to new goroutine...
	Go(func() {
		// read specific key under current session
		val := GetCurSession().Get(key) // val exists
		ASSERT(val == v)
		// doSomething....

		// set specific key-value under current session
		// NOTICE: current session won't change here
		next := GetCurSession().WithValue(key2, v2)
		val2 := GetCurSession().Get(key2) // val2 == nil
		ASSERT(val2 == nil)

		// pass both parent session and new session to sub goroutine
		GoSession(next, func() {
			// read specific key under current session
			val := GetCurSession().Get(key) // val exists
			ASSERT(val == v)

			val2 := GetCurSession().Get(key2) // val2 exists
			ASSERT(val2 == v2)
			// doSomething....

			sig2 <- struct{}{}

			<-sig
			ASSERT(GetCurSession().IsValid() == false) // current session is invalid

			println("g2 done")
			sig2 <- struct{}{}
		})

		Go(func() {
			// read specific key under current session
			val := GetCurSession().Get(key) // val exists
			ASSERT(v == val)

			val2 := GetCurSession().Get(key2) // val2 == nil
			ASSERT(val2 == nil)
			// doSomething....

			sig2 <- struct{}{}

			<-sig
			ASSERT(GetCurSession().IsValid() == false) // current session is invalid

			println("g3 done")
			sig2 <- struct{}{}
		})

		BindSession(next)
		val2 = GetCurSession().Get(key2) // val2 exists
		ASSERT(v2 == val2)

		sig2 <- struct{}{}

		<-sig
		ASSERT(next.IsValid() == false) // next is invalid

		println("g1 done")
		sig2 <- struct{}{}
	})

	<-sig2
	<-sig2
	<-sig2

	val2 := GetCurSession().Get(key2) // val2 == nil
	ASSERT(val2 == nil)

	// initiatively ends the session，
	// then all the inherited session (including next) will be disabled
	session.Disable()
	close(sig)

	ASSERT(start.IsValid() == false) // start is invalid

	<-sig2
	<-sig2
	<-sig2
	println("g0 done")

	UnbindSession()
}
```


## Community
- Email: [conduct@cloudwego.io](conduct@cloudwego.io)
- How to become a member: [COMMUNITY MEMBERSHIP](https://github.com/cloudwego/community/blob/main/COMMUNITY_MEMBERSHIP.md)
- Issues: [Issues](https://github.com/cloudwego/localsession/issues)
- Slack: Join our CloudWeGo community [Slack Channel](https://join.slack.com/t/cloudwego/shared_invite/zt-tmcbzewn-UjXMF3ZQsPhl7W3tEDZboA).
