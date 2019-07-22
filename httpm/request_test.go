package httpm_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/moxar/middleman/httpm"
)

func TestQNew(t *testing.T) {
	r, err := httpm.QNew("POST", "https://github.com")(nil)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Method != "POST" {
		t.Errorf("Method should be 'POST', '%s' given", r.Method)
		return
	}

	want, err := url.Parse("https://github.com")
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(r.URL, want) {
		t.Error(r.URL, want)
		return
	}
}

func TestQEncodeDecodeBody(t *testing.T) {
	type Payload struct {
		Foo string
		Bar int
	}
	in := Payload{Foo: "foo", Bar: 4}
	r, err := httpm.QEncodeBody(json.Marshal)(in)(new(http.Request))
	if err != nil {
		t.Error(err)
		return
	}

	var out Payload
	r, err = httpm.QDecodeBody(json.Unmarshal)(&out)(r)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(in, out) {
		t.Error(in, out)
		return
	}
	if r.GetBody == nil {
		t.Error("GetBody should not be nil")
		return
	}

	if r.ContentLength == 0 {
		t.Error("ContentLength should not be 0")
		return
	}
}

func TestComposeQFn(t *testing.T) {

	t.Run("on happy case", func(t *testing.T) {
		type Payload struct {
			Foo string
			Bar int
		}
		in := Payload{Foo: "foo", Bar: 4}

		r, err := httpm.ComposeQFn(
			httpm.QNew("POST", "https://github.com"),
			httpm.QEncodeBody(json.Marshal)(in),
		)(nil)
		if err != nil {
			t.Error(err)
			return
		}

		var out Payload
		r, err = httpm.QDecodeBody(json.Unmarshal)(&out)(r)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(in, out) {
			t.Error(in, out)
			return
		}

		if r.Method != "POST" {
			t.Errorf("Method should be 'POST', '%s' given", r.Method)
			return
		}
		want, err := url.Parse("https://github.com")
		if err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(r.URL, want) {
			t.Error(r.URL, want)
			return
		}
	})

	t.Run("on failing QFn", func(t *testing.T) {
		var i int
		failer := func(r *http.Request) (*http.Request, error) {
			i++
			return nil, errors.New("boom")
		}
		_, err := httpm.ComposeQFn(
			failer,
			failer,
			failer,
		)(nil)
		if err == nil {
			t.Fail()
		}
		if i != 1 {
			t.Errorf("too many Fn called, expected 1, having %d", i)
			return
		}
	})
}

func ExampleComposeQFn() {

	// Define a noop function, witness print order.
	Println := func(vs ...interface{}) httpm.QFn {
		return func(r *http.Request) (*http.Request, error) {
			fmt.Println(vs...)
			return r, nil
		}
	}

	// Compose request.
	r, err := httpm.ComposeQFn(
		Println("Foo"),
		httpm.QNew("POST", "https://github.com"),
		Println("Bar"),
		httpm.QEncodeBody(json.Marshal)(map[string]string{"foo": "bar"}),
	)(nil)

	if err != nil {
		log.Println(err)
	}

	// use r
	_ = r

	// Output:
	// Foo
	// Bar
}
