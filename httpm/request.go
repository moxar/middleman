package httpm

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// RequestFn describes a function that can be member of a chainable Request handling.
type RequestFn = func(*http.Request) (*http.Request, error)

// ComposeRequest composes a list of RequestFn into one.
// ComposeRequest(foo, bar) is functionnally equivalent to bar(foo(r))
func ComposeRequest(fns ...RequestFn) RequestFn {
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

// NewRequest prepares a http request with the standard http.NewRequest method.
func NewRequest(method, url string) RequestFn {
	return func(*http.Request) (*http.Request, error) {
		return http.NewRequest(method, url, nil)
	}
}

// WriteRequestBody encodes and writes the given input in the body.
func WriteRequestBody(e Encoder) func(input interface{}) RequestFn {
	return func(input interface{}) RequestFn {
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

// SetRequestContext adds the context to the request.
func SetRequestContext(ctx context.Context) RequestFn {
	return func(r *http.Request) (*http.Request, error) {
		return r.WithContext(ctx), nil
	}
}

// SetRequestHeader sets the header of the request.
func SetRequestHeader(h http.Header) RequestFn {
	return func(r *http.Request) (*http.Request, error) {
		r.Header = h
		return r, nil
	}
}

// ReadRequestBody decodes the request body.
func ReadRequestBody(d Decoder) func(into interface{}) RequestFn {
	return func(into interface{}) RequestFn {
		return func(r *http.Request) (*http.Request, error) {
			raw, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			return r, d(raw, into)
		}
	}
}

// ParamParser designates a function capable of parsing URL params.
type ParamParser = func(interface{}, url.Values) error

// ReadRequestParams decodes the request parameters.
func ReadRequestParams(fn ParamParser) func(interface{}) RequestFn {
	return func(into interface{}) RequestFn {
		return func(r *http.Request) (*http.Request, error) {
			return r, fn(into, map[string][]string(r.URL.Query()))
		}
	}
}

// Checker is a func that checks the validity of an input.
type Checker = func(interface{}) error

// DecodeAndCheck the input.
func DecodeAndCheck(decode Decoder, check Checker) Decoder {
	return func(raw []byte, input interface{}) error {
		if err := decode(raw, input); err != nil {
			return err
		}
		return check(input)
	}
}
