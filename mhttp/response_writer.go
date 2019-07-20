package mhttp

import "net/http"

// WFn describes a function that can be member of a chainable ResponseWriter handling.
type WFn = func(w http.ResponseWriter) http.ResponseWriter

// ComposeWFn composes a list of WFn into one.
// ComposeWFn(foo, bar) is functionnally equivalent to bar(foo(r))
func ComposeWFn(fn ...WFn) WFn {
	return func(w http.ResponseWriter) http.ResponseWriter {
		for _, f := range fn {
			w = f(w)
		}
		return w
	}
}

// WWriteTextBody writes the input text in the ResponseWriter.
func WWriteTextBody(txt string) WFn {
	return func(w http.ResponseWriter) http.ResponseWriter {
		w.Write([]byte(txt)) // nolint: errcheck
		return w
	}
}

// WWriteStatus writes the input status code as header.
func WWriteStatus(status int) WFn {
	return func(w http.ResponseWriter) http.ResponseWriter {
		w.WriteHeader(status)
		return w
	}
}

// WEncodeBody encodes and writes the given input in the body.
func WEncodeBody(e Encoder) func(interface{}) WFn {
	return func(input interface{}) WFn {
		return func(w http.ResponseWriter) http.ResponseWriter {
			raw, _ := e(input)
			w.Write(raw) // nolint: errcheck
			return w
		}
	}
}
