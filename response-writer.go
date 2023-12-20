package httputils

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, Body: new(bytes.Buffer)}
}

func (crw *ResponseWriter) Write(data []byte) (int, error) {
	crw.Body.Write(data)
	return crw.ResponseWriter.Write(data)
}

func (crw *ResponseWriter) WriteHeader(statusCode int) {
	crw.StatusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
