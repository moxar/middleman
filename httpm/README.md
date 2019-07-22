# httpm

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/moxar/middleman/httpm)

`httpm` is the middleman subpackage specialized in http.

## Motivation

As golang developer, we often thrive under the weight of error checking. Though error management in go
ensures that the code is resillient, it also makes it very verbose, and bubbeling errors up in the stack
adds a lot of boilerplate.

```go
// Parse content from HTTP request, get header, get params...
raw, err := ioutil.ReadAll(r.Body)
if err != nil {
	// ...
}
err := json.Unmarshal(raw, &payload)
if err != nil {
	// ...
}
err := decodeURIParams(r, &params)
if err != nil {
	// ...
}

// ...

// Start business logic
```

This package provides functions and functors that can be used to prevent this.
