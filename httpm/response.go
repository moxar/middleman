package httpm

import (
	"io/ioutil"
	"net/http"
)

// RFn describes a function that can be member of a chainable Response handling.
type RFn = func(*http.Response) (*http.Response, error)

// ComposeRFn composes a list of RFns into one.
// ComposeRFn(foo, bar) is functionnally equivalent to bar(foo(w))
func ComposeRFn(fns ...RFn) RFn {
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

// RDecodeBody decodes the request body.
func RDecodeBody(d Decoder) func(interface{}) RFn {
	return func(input interface{}) RFn {
		return func(r *http.Response) (*http.Response, error) {
			raw, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			return r, d(raw, input)
		}
	}
}

// RErrorFromStatus returns an error depending on the response status.
func RErrorFromStatus(f func(status int) error) RFn {
	return func(r *http.Response) (*http.Response, error) {
		return r, f(r.StatusCode)
	}
}
