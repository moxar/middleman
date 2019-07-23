# httpm

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/moxar/middleman/contextm)

`contextm` is the middleman subpackage specialized in context. It provides a simple pair of functions to iterate over a context.


```go
ctx := context.Background()
ctx = WithRange(ctx, "first Avenger", "Steve Rogers")
ctx = WithRange(ctx, "expert tinker", "Tony Stark")
ctx = context.WithValue(ctx, ctxKey("not an Avenger"), "Batman")
ctx = WithRange(ctx, "friendly neighbour", "Peter Parker")

Range(ctx, func(key, val interface{}) bool {
	fmt.Println(key, val)
	return true
})
fmt.Println(Keys(ctx))

// Output:
// first Avenger Steve Rogers
// expert tinker Tony Stark
// friendly neighbour Peter Parker
// [first Avenger expert tinker friendly neighbour]

```
