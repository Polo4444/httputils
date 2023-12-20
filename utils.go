package httputils

import (
	"encoding/json"
	"io"
	"net/http"
)

type RespResult struct {
	Result  string   `json:"result"`
	Details []string `json:"details"`
}

func errorsToStrings(errs []error) []string {

	result := make([]string, len(errs))
	for i := range errs {
		result[i] = errs[i].Error()
	}
	return result
}

func writeOKStatus(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func writeCustomStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func ReturnBytes(w http.ResponseWriter, v []byte) {

	w.Header().Set("Content-Type", "application/octet-stream")
	writeOKStatus(w)
	w.Write(v)
}

func ReturnBytesWithCode(w http.ResponseWriter, code int, v []byte) {

	w.Header().Set("Content-Type", "application/octet-stream")
	writeCustomStatus(w, code)
	w.Write(v)
}

func ReturnString(w http.ResponseWriter, v string) {

	w.Header().Set("Content-Type", "text/plain")
	writeOKStatus(w)
	w.Write([]byte(v))
}

func ReturnStringWithCode(w http.ResponseWriter, code int, v string) {

	w.Header().Set("Content-Type", "text/plain")
	writeCustomStatus(w, code)
	w.Write([]byte(v))
}

func ReturnReader(w http.ResponseWriter, v io.Reader) {

	w.Header().Set("Content-Type", "application/octet-stream")
	writeOKStatus(w)
	io.Copy(w, v)
}

func ReturnReaderWithCode(w http.ResponseWriter, code int, v io.Reader) {

	w.Header().Set("Content-Type", "application/octet-stream")
	writeCustomStatus(w, code)
	io.Copy(w, v)
}

func ReturnJSON(w http.ResponseWriter, v interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	writeOKStatus(w)
	return json.NewEncoder(w).Encode(v)
}

func ReturnJSONWithCode(w http.ResponseWriter, code int, v interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	writeCustomStatus(w, code)
	return json.NewEncoder(w).Encode(v)
}

func ReturnError(w http.ResponseWriter, code int, message string, err ...error) error {

	w.Header().Set("Content-Type", "application/json")
	writeCustomStatus(w, code)
	return json.NewEncoder(w).Encode(&RespResult{
		Result:  message,
		Details: errorsToStrings(err),
	})
}

func ReturnOK(w http.ResponseWriter) error {

	w.Header().Set("Content-Type", "application/json")
	writeOKStatus(w)
	return json.NewEncoder(w).Encode(&RespResult{
		Result: "ok",
	})
}
