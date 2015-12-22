# atomicmapper

atomicmapper is a code generation tool for creating high-performance, scalable, frequently read, but infrequently updated maps of strings to any given type `map[string]YourType`. It is based on Go's [atomic.Value read mostly example](https://golang.org/pkg/sync/atomic/#example_Value_readMostly).

**Requires Go 1.4+**

## usage with go:generate

Install atomicmapper.

```bash
go install github.com/ninibe/atomicmapper
# move to $PATH if $GOPATH/bin is not in your $PATH
```

Add the generate command.

```go
//go:generate atomicmapper -pointer -type Foo
type Foo struct { ... }
```

Skip the `-pointer` flag to save entire values.
Generate the atomic map code.

```bash
go generate myfoopkg
```

This will create a new `foo_atomicmap.go` file ready to use.

```go
fooMap := NewFooAtomicMap()
fooMap.Set("myKey", &Foo{}) // save pointer to Foo
foo := fooMap.Get("myKey")  // retrieve pointer
fooMap.Delete("myKey")      // remove pointer from map
```

All methods are thread-safe while `Get` is also lock-free.
Check the example [godoc](https://godoc.org/github.com/ninibe/atomicmapper/test)

## usage with gen

See [atomic map typewriter](https://github.com/ninibe/atomicmapper/tree/master/gen)