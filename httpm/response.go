package httpm

import (
	"io/ioutil"
	"net/http"
)

// ResponseFn describes a function that can be member of a chainable Response handling.
type ResponseFn = func(*http.Response) (*http.Response, error)

// ComposeResponse composes a list of ResponseFns into one.
// ComposeResponse(foo, bar) is functionnally equivalent to bar(foo(w))
func ComposeResponse(fns ...ResponseFn) ResponseFn {
	return func(r *http.Response) (*http.Response, error) {
		var err error
		for _, fn := range fns {
			if r, err = fn(r); err != nil {
				return nil, err
			}
		}
		return r, nil
	}
}

// ReadResponseBody decodes the body.
func ReadResponseBody(d Decoder) func(into interface{}) ResponseFn {
	return func(into interface{}) ResponseFn {
		return func(r *http.Response) (*http.Response, error) {
			raw, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			return r, d(raw, into)
		}
	}
}

// ReadResponseStatus into status.
func ReadResponseStatus(status *int) ResponseFn {
	return func(r *http.Response) (*http.Response, error) {
		*status = r.StatusCode
		return r, nil
	}
}

// ReturnErrorFromResponseStatus returns an error depending on the response status.
func ReturnErrorFromResponseStatus(f func(status int) error) ResponseFn {
	return func(r *http.Response) (*http.Response, error) {
		return r, f(r.StatusCode)
	}
}
