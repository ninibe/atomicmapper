## [atomic map](https://github.com/ninibe/atomicmapper) typewriter for the [gen (v4)](https://clipperhouse.github.io/gen/) code generation tool

Once you have `gen` installed, find the folder with the target type and install the typewriter

```
cd
gen add github.com/ninibe/atomicmapper/gen
```

This creates a `_gen.go` file that you can remove after generation.
Set the tag over your target type.

```go
// +gen atomicmap
type Foo struct {}
```

or

```go
// +gen * atomicmap
type Foo struct {}
```

and run `gen` on the folder to generate the atomic map code.