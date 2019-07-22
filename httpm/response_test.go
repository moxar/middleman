package httpm_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/moxar/middleman/httpm"
)

func TestReadResponseBody(t *testing.T) {
	type Payload struct {
		Foo string
		Bar int
	}

	r := &http.Response{}
	r.Body = ioutil.NopCloser(strings.NewReader(`{"foo":"foo", "bar": 4}`))

	var got Payload
	r, err := httpm.ReadResponseBody(json.Unmarshal)(&got)(r)
	if err != nil {
		t.Error(err)
		return
	}

	want := Payload{Foo: "foo", Bar: 4}
	if !reflect.DeepEqual(got, want) {
		t.Error(got, want)
		return
	}
}

func TestReturnErrorFromResponseStatus(t *testing.T) {
	e := func(s int) error {
		if s <= http.StatusBadRequest {
			return nil
		}
		return errors.New(http.StatusText(s))
	}
	t.Run("on happy case", func(t *testing.T) {

		r := &http.Response{StatusCode: 202}
		_, err := httpm.ReturnErrorFromResponseStatus(e)(r)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("on bad status", func(t *testing.T) {
		r := &http.Response{StatusCode: 500}
		_, err := httpm.ReturnErrorFromResponseStatus(e)(r)
		if err == nil {
			t.Fail()
			return
		}
	})
}

func TestComposeResponse(t *testing.T) {
	t.Run("on happy case", func(t *testing.T) {
		var order []string
		succeeder := func(val string) httpm.ResponseFn {
			return func(r *http.Response) (*http.Response, error) {
				order = append(order, val)
				return r, nil
			}
		}

		_, err := httpm.ComposeResponse(
			succeeder("A"),
			succeeder("B"),
			succeeder("C"),
		)(new(http.Response))
		if err != nil {
			t.Error(err)
			return
		}

		want := []string{"A", "B", "C"}
		if !reflect.DeepEqual(order, want) {
			t.Error(order, want)
			return
		}
	})

	t.Run("on failing ResponseFn", func(t *testing.T) {
		var i int
		failer := func(r *http.Response) (*http.Response, error) {
			i++
			return nil, errors.New("boom")
		}

		_, err := httpm.ComposeResponse(
			failer,
			failer,
			failer,
		)(new(http.Response))
		if err == nil {
			t.Fail()
			return
		}

		if i != 1 {
			t.Errorf("too many Fn called, expected 1, having %d", i)
			return
		}
	})
}

func ExampleComposeResponse() {

	// Define a noop function, witness print order.
	Println := func(vs ...interface{}) httpm.ResponseFn {
		return func(r *http.Response) (*http.Response, error) {
			fmt.Println(vs...)
			return r, nil
		}
	}

	FailOver400 := func(s int) error {
		if s <= 400 {
			return nil
		}
		return errors.New(http.StatusText(s))
	}

	r := &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(`{"foo": "car", "bar": "baz"}`)),
		StatusCode: 305,
	}

	// Compose response.
	var out = make(map[string]string)
	r, err := httpm.ComposeResponse(
		Println("Foo"),
		httpm.ReturnErrorFromResponseStatus(FailOver400),
		Println("Bar"),
		httpm.ReadResponseBody(json.Unmarshal)(&out),
	)(r)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(out["foo"])
	fmt.Println(out["bar"])

	// use r
	_ = r

	// Output:
	// Foo
	// Bar
	// car
	// baz
}
