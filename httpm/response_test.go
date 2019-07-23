package httpm_test

import (
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

	r := &http.Response{}
	r.Body = ioutil.NopCloser(strings.NewReader(`some string`))

	var got []byte
	r, err := httpm.ReadResponseBody(httpm.DecodeText)(&got)(r)
	if err != nil {
		t.Error(err)
		return
	}

	want := "some string"
	if string(got) != want {
		t.Fail()
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

func ExampleComposeResponse_client() {

	FailOver400 := func(s int) error {
		if s <= 400 {
			return nil
		}
		return errors.New(http.StatusText(s))
	}

	// Compose response.
	ParseResponse := func(r *http.Response, status *int, payload *string) error {

		_, err := httpm.ComposeResponse(
			httpm.ReadResponseStatus(status),
			httpm.ReturnErrorFromResponseStatus(FailOver400),
			httpm.ReadResponseBody(httpm.DecodeText)(payload),
		)(r)
		return err
	}

	r := &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(`some payload`)),
		StatusCode: 305,
	}

	var out string
	var status int
	err := ParseResponse(r, &status, &out)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(out)
	fmt.Println(status)

	// use r
	_ = r

	// Output:
	// some payload
	// 305
}
