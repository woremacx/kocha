package kocha

import (
	"io"
	"net/http"
	"net/http/httptest"
)

var _ http.ResponseWriter = &Response{}

// Response represents a response.
type Response struct {
	http.ResponseWriter
	ContentType string
	StatusCode  int

	cookies []*http.Cookie
	resp    *httptest.ResponseRecorder
}

// newResponse returns a new Response that responds to rw.
func newResponse() *Response {
	r := &Response{}
	r.reset()
	return r
}

// Cookies returns a slice of *http.Cookie.
func (r *Response) Cookies() []*http.Cookie {
	return r.cookies
}

// SetCookie adds a Set-Cookie header to the response.
func (r *Response) SetCookie(cookie *http.Cookie) {
	r.cookies = append(r.cookies, cookie)
	http.SetCookie(r, cookie)
}

func (r *Response) writeTo(w http.ResponseWriter) error {
	for key, values := range r.Header() {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(r.resp.Code)
	_, err := io.Copy(w, r.resp.Body)
	return err
}

func (r *Response) reset() {
	r.StatusCode = http.StatusOK
	r.resp = httptest.NewRecorder()
	r.ResponseWriter = r.resp
}
