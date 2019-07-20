# mhttp

## Usage

mhttp has **four** main purposes:
  * encoding client http request
  * decoding server http request
  * encoding server http response
  * decoding client http response
  
Those purposes are achieved with 3 types: `RequestChainer`, `ResponseChainer` and `ResponseWriterChainer`.

On client side
```go
type Request struct{
  Foo, Bar string
}
  
type Response struct{
  Baz string
}

var(
  // Instructions to prepare the request.
  prepare = func(payload interface{}) httpkit.RequestChainer {
    return Request(
      httpkit.NewRequest("POST", "https://my.api.org/path"),
      httpkit.EncodeRequestBody(json.Marshal)(payload),
      // Change http protocol ?
      // Set headers ? (content-type, auth, etc...)
    )  
  }
  
  // Instruction to consume the response.
  consume = func(payload interface{}) httpkit.ResponseChainer {
    return Response(
      // Check http status ?
      httpkit.DecodeResponseBody(json.Unmarshal)(payload),
    )
  }
  
  // Sefault sending function.
  send = httpkit.Send(prepare, consume)
)

func main() {
  req := Request{Foo: "foo", Bar: "bar"}
  var resp Response
  
  if err := send(prepare(req), consume(&resp)); err != nil {
    // handle error
  }
}
```

On server side
```go
type Request struct{
  Foo, Bar string
}
  
type Response struct{
  Baz string
}

var(
  decode = func(r *http.Request, payload interface{}) error {
    return httpkit.Request(
      // bind request parameters ? (url/{user}/{template})
      // parse request query string array ? (?isset=true)
      httpkit.DecodeRequestBody(json.Unmarshal)(payload),
    )(r)
  }
  
  encode = func(w http.ResponseWriter, payload interface{}){
    return httpkit.ResponseWriter(
      // Set http status ?
      httpkit.EncodeResponseWriterBody(json.Marshal)(payload),
    )(w)
  }
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var payload Request
  if err := decode(r, &payload); err != nil {
    encode(r, err)
    return
  }
  
  var output Response
  output.Baz = "baz"
  encode(w, output)
}
```
