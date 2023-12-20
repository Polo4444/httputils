package caching

import (
	"bytes"
	"crypto/md5"
	"io"
	"net/http"
	"sync"

	"github.com/Polo4444/httputils"
)

type cachedRequestState int

const (
	cachedRequestStateRunning cachedRequestState = iota
	cachedRequestStateDone
	cachedRequestStateWaitingCache
)

var cachedRequests = map[string]requestsCaching{}
var cachingMu = sync.Mutex{}

type cachedRequest struct {
	key      string
	state    cachedRequestState
	sendData func(w *httputils.ResponseWriter)
}

type requestsCaching struct {
	w         *httputils.ResponseWriter
	notifiers map[string]*cachedRequest
}

func computeCachedRequestKey(r *http.Request, headers []string) string {

	// Grab full URL including query params
	key := []byte(r.URL.String())

	// Add the request method
	key = append(key, []byte(r.Method)...)

	// Add the request body
	body, err := io.ReadAll(r.Body)
	if err == nil {
		key = append(key, body...)
	}
	// Close the original body
	_ = r.Body.Close()

	// Replace r.Body with a new reader containing the same data
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Add the headers
	for _, v := range headers {
		key = append(key, []byte(r.Header.Get(v))...)
	}

	// create md5 hash
	md5New := md5.New()
	md5New.Write(key)
	md5hash := md5New.Sum(nil)

	// return hex encoded hash
	return string(md5hash[:])
}

func NewCachedRequest(w http.ResponseWriter, r *http.Request, headers []string) *cachedRequest {

	key := computeCachedRequestKey(r, headers)
	cr := &cachedRequest{key: key}
	if _, exist := cachedRequests[key]; !exist {
		cr.state = cachedRequestStateRunning
		cachedRequests[key] = requestsCaching{w: httputils.NewResponseWriter(w)}
		return cr
	}

	cr.state = cachedRequestStateWaitingCache
	reqCaching := cachedRequests[key]
	if _, exist := reqCaching.notifiers[key]; exist {
		reqCaching.notifiers[key] = cr
	} else {
		reqCaching.notifiers = map[string]*cachedRequest{key: cr}
	}
	cachedRequests[key] = reqCaching

	return cr
}
