package httpm

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// QFn describes a function that can be member of a chainable Request handling.
type QFn = func(*http.Request) (*http.Request, error)

// ComposeQFn composes a list of QFn into one.
// ComposeQFn(foo, bar) is functionnally equivalent to bar(foo(r))
func ComposeQFn(fns ...QFn) QFn {
	return func(r *http.Request) (*http.Request, error) {
		var err error
		for _, fn := range fns {
			if r, err = fn(r); err != nil {
				return nil, err
			}
		}
		return r, nil
	}
}

// QNew prepares a http request with the standard http.NewRequest method.
func QNew(method, url string) QFn {
	return func(*http.Request) (*http.Request, error) {
		return http.NewRequest(method, url, nil)
	}
}

// QEncodeBody encodes and writes the given input in the body.
func QEncodeBody(e Encoder) func(interface{}) QFn {
	return func(input interface{}) QFn {
		return func(r *http.Request) (*http.Request, error) {
			raw, err := e(input)
			if err != nil {
				return nil, err
			}
			r.ContentLength = int64(len(raw))
			r.GetBody = func() (io.ReadCloser, error) {
				return ioutil.NopCloser(bytes.NewReader(raw)), nil
			}
			r.Body, _ = r.GetBody()
			return r, nil
		}
	}
}

// QDecodeBody decodes the request body.
func QDecodeBody(d Decoder) func(interface{}) QFn {
	return func(input interface{}) QFn {
		return func(r *http.Request) (*http.Request, error) {
			raw, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			return r, d(raw, input)
		}
	}
}
