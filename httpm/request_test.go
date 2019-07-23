package httpm_test

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/moxar/middleman/httpm"
)

func TestNewRequest(t *testing.T) {
	r, err := httpm.NewRequest("POST", "https://github.com")(nil)
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

func TestWriteRequestBody(t *testing.T) {
	in := "some string"
	r, err := httpm.WriteRequestBody(httpm.EncodeText)(in)(new(http.Request))
	if err != nil {
		t.Error(err)
		return
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		return
	}
	out := string(raw)

	if out != in {
		t.Fail()
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

func TestComposeRequest(t *testing.T) {

	t.Run("on happy case", func(t *testing.T) {

		newRequest := func(path, url string, input interface{}) (*http.Request, error) {
			return httpm.ComposeRequest(
				httpm.NewRequest(path, url),
				httpm.WriteRequestBody(httpm.EncodeText)(input),
			)(nil)
		}

		in := "some payload"
		r, err := newRequest("POST", "https://github.com", in)
		if err != nil {
			t.Error(err)
			return
		}
		raw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		out := string(raw)
		if out != in {
			t.Fail()
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

	t.Run("on failing RequestFn", func(t *testing.T) {
		var i int
		failer := func(r *http.Request) (*http.Request, error) {
			i++
			return nil, errors.New("boom")
		}
		_, err := httpm.ComposeRequest(
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

func ExampleComposeRequest_client() {

	var encode httpm.Encoder // json.Marshal

	NewRequest := func(path, url string, input interface{}) (*http.Request, error) {
		return httpm.ComposeRequest(
			httpm.NewRequest(path, url),
			httpm.WriteRequestBody(encode)(input),
		)(nil)
	}

	type Hero struct {
		Name     string
		Universe string
	}

	// Compose request.
	batman := Hero{Name: "Batman", Universe: "DC"}
	r, err := NewRequest("POST", "https://api.superheroes.com/heroes", batman)
	if err != nil {
		log.Println(err)
	}

	// use r
	_ = r
}

func ExampleComposeRequest_server() {

	var parseParams httpm.ParamParser // gorilla/schema.NewDecoder().Decode
	var check httpm.Checker           // asaskevich/govalidator.ValidateStruct
	var decode httpm.Decoder          // json.Unmarshal

	// parseRequest parses the request params (?foo=bar) and body.
	parseRequest := func(r *http.Request, body, params interface{}) error {
		_, err := httpm.ComposeRequest(
			httpm.ReadRequestBody(httpm.DecodeAndCheck(decode, check))(body),
			httpm.ReadRequestParams(parseParams)(params),
		)(r)
		return err
	}

	type Hero struct {
		Name     string
		Universe string
	}

	type Params struct {
		FastInsert bool
	}

	// HTTP Handler...
	handle := func(w http.ResponseWriter, r *http.Request) {
		var hero Hero
		var params Params
		if err := parseRequest(r, &hero, &params); err != nil {
			// ...
		}

		// use hero and params values
	}

	_ = handle
}
