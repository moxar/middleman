# httpm

`httpm` is the middleman subpackage specialized in http.

## Motivation

As golang developer, we often thrive under the weight of error checking. Though error management in go
ensures that the code is resillient, it also makes it very verbose, and bubbeling errors up in the stack
adds a lot of boilerplate.

```go
// Parse content from HTTP request, get header, get params...
raw, err := ioutil.ReadAll(r.Body)
if err != nil {
	// ...
}
err := json.Unmarshal(raw, &payload)
if err != nil {
	// ...
}
err := decodeURIParams(r, &params)
if err != nil {
	// ...
}

// ...

// Start business logic
```

This package provides functions and functors that can be used to prevent this.

## Usage

This section only provides a basic example of usage. Checkout documentation for detailed examples.

```go
type Request struct{
	Foo, Bar string
}

type Response struct{
	Baz string
}

var(

	writeRequest = func(method, url string, payload interface{}) httpm.QFn {
  		return httpm.ComposeQFn(
			httpm.QNew(method, url),
			httpm.QEncode(json.Marshal)(payload),
			// add other function to apply to the request here.
			// Add HTTP header ? basic auth ?
		)
 	}

	readRequest = func(r *http.Request, payload interface{}) *http.Request {
  		return httpm.ComposeQFn(
			httpm.QDecode(json.Unmarshal)(payload),
			// add other function to apply to the request here.
			// Want to parse headers ? request parameters ?
		)(r)
	}

	writeResponse = func(w http.ResponseWriter, payload interface{}, status int) http.ResponseWriter {
  		return httpm.ComposeWFn(
			httpm.WWriteStatus(status),
			httpm.WEncode(json.Marshal)(payload),
		)(w)
  	}

	readResponse = func(payload interface{}) httpm.RFn {
	  	return httpm.ComposeRFn(
			httpm.RDecode(json.Unmarshal)(payload),
		)
 	}
)

func Call() {
	in := Request{Foo:"foo", Bar:"bar"}
	var out Response
	if err := httpm.Send(writeRequest(in), readResponse(&out)); err != nil {
		// ...
	}
	// out is ready for usage.
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload Request
	if err := readRequest(r, &payload); err != nil {
		// ...
	}

  	// the request is ready to be used by business logic.
	var out Response
	writeResponse(w, out, http.StatusOK)
}
```
