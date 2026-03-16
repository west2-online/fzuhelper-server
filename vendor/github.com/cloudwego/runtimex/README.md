# runtimex

Runtimex package help to expose Go Runtime internals representation safely.

## Usage

### Get Goroutine ID

```go
gid, err := runtimex.GID()
```

### Get Processor ID

```go
pid, err := runtimex.PID()
```

## Note

Since we use a hack way to expose internal representation of the Go runtime, so if Go change some internal variable names, the package will return error.

You should care about the error returned by runtimex and do the fallback logic If necessary.

For now, we only depend on `runtime.g` and `g.goid`.
