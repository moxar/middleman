package httpm

import "net/http"

// ResponseWriterFn describes a function that can be member of a chainable ResponseWriter handling.
type ResponseWriterFn = func(w http.ResponseWriter) http.ResponseWriter

// ComposeResponseWriter composes a list of ResponseWriterFn into one.
// ComposeResponseWriter(foo, bar) is functionnally equivalent to bar(foo(r))
func ComposeResponseWriter(fn ...ResponseWriterFn) ResponseWriterFn {
	return func(w http.ResponseWriter) http.ResponseWriter {
		for _, f := range fn {
			w = f(w)
		}
		return w
	}
}

// WriteTextInResponseWriterBody writes the input text in the ResponseWriter.
func WriteTextInResponseWriterBody(txt string) ResponseWriterFn {
	return func(w http.ResponseWriter) http.ResponseWriter {
		w.Write([]byte(txt)) // nolint: errcheck
		return w
	}
}

// WriteResponseWriterStatus writes the input status code as header.
func WriteResponseWriterStatus(status int) ResponseWriterFn {
	return func(w http.ResponseWriter) http.ResponseWriter {
		w.WriteHeader(status)
		return w
	}
}

// WriteResponseWriterBody encodes and writes the given input in the body.
func WriteResponseWriterBody(e Encoder) func(interface{}) ResponseWriterFn {
	return func(input interface{}) ResponseWriterFn {
		return func(w http.ResponseWriter) http.ResponseWriter {
			raw, _ := e(input)
			w.Write(raw) // nolint: errcheck
			return w
		}
	}
}
