package httpm

import (
	"errors"
	"net/http"
)

// Doer describes a type capable of sending Requests and returning responses.
// *http.Client is the cannonical implementation of a Doer.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Send a request and handles its response with the http.DefaultClient.
//  If in is nil, Send returns an error.
//  If out is nil, the response is not handled. This can be useful for fire & forget usages.
func Send(in RequestFn, out ResponseFn) error {
	return NewSender(http.DefaultClient)(in, out)
}

// NewSender returns a function that sends a request and handling its
// response with the given doer. See Send for details.
func NewSender(client Doer) func(RequestFn, ResponseFn) error {
	return func(in RequestFn, out ResponseFn) error {
		if in == nil {
			return errors.New("RequestFn is nil")
		}
		req, err := in(&http.Request{})
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if out == nil {
			return nil
		}
		_, err = out(resp)
		return err
	}
}
