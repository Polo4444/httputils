package caching

import (
	"net/http"

	"github.com/Polo4444/httputils"
)

// Caching is the core struct for caching
type Caching struct {
	headers []string // Headers to use for caching
}

// NewCaching creates a new caching struct
func NewCaching(headers []string) *Caching {

	if headers == nil {
		headers = []string{}
	}

	return &Caching{headers: headers}
}

// CahcingMiddleware returns the middleware for caching
func (c *Caching) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cachingMu.Lock()
		cr := NewCachedRequest(w, r, c.headers)
		cachingMu.Unlock()

		if cr.state == cachedRequestStateRunning {

			reqCaching := cachedRequests[cr.key]
			next.ServeHTTP(reqCaching.w, r)

			// After the request is done, we notify other requests waiting for the same cache
			reqCaching = cachedRequests[cr.key]
			for _, v := range reqCaching.notifiers {
				go v.sendData(reqCaching.w)
			}
			delete(cachedRequests, cr.key)
			return
		}

		// ─── Waiting Cache Requests ──────────────────────────────────

		// We create a channel to block
		ch := make(chan bool)

		// We don't serve, we just wait to be notified the cache is ready
		cr.sendData = func(wc *httputils.ResponseWriter) {

			// Set the headers
			for k, v := range wc.Header() {
				w.Header().Set(k, v[0])
			}

			w.WriteHeader(wc.StatusCode)
			w.Write(wc.Body.Bytes())

			// We write on the channel to unblock
			ch <- true
		}

		// reqCaching := CachedRequests[cr.key]
		// reqCaching.notifiers[cr.key] = cr
		// CachedRequests[cr.key] = reqCaching

		// We block until we are notified
		<-ch
	})
}
