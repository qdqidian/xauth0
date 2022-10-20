package xauth0

import (
	"bytes"
	"net/http"
)

type Auth0Writer struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w Auth0Writer) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w Auth0Writer) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}
