package httpm_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moxar/middleman/httpm"
)

func TestComposeResponseWriter(t *testing.T) {

	writeResponse := func(w http.ResponseWriter, payload interface{}, status int) {
		httpm.ComposeResponseWriter(
			httpm.WriteResponseWriterStatus(status),
			httpm.WriteResponseWriterBody(httpm.EncodeText)(payload),
		)(w)
	}

	readResponse := func(r *http.Response, output interface{}, status *int) error {
		_, err := httpm.ComposeResponse(
			httpm.ReadResponseBody(httpm.DecodeText)(output),
			httpm.ReadResponseStatus(status),
		)(r)
		return err
	}

	var payload = "my response"
	var status = 202

	var output []byte
	var outStatus int

	w := httptest.NewRecorder()
	writeResponse(w, payload, status)
	if err := readResponse(w.Result(), &output, &outStatus); err != nil {
		t.Error(err)
		return
	}

	if string(output) != payload {
		t.Errorf("%s != %s", output, payload)
		return
	}

	if status != outStatus {
		t.Error(status, "!=", outStatus)
		return
	}
}

func ExampleComposeResponseWriter_server() {
	writeResponse := func(w http.ResponseWriter, payload interface{}, status int) {
		httpm.ComposeResponseWriter(
			httpm.WriteResponseWriterStatus(status),
			httpm.WriteResponseWriterBody(httpm.EncodeText)(payload),
		)(w)
	}

	var w http.ResponseWriter
	writeResponse(w, "some payload", http.StatusNoContent)

}
