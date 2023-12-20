## HTTPUTILS
This package holds a list of simple yet useful functions for Go standard HTTP package.

## Caching package
This package holds a simple caching system that can be used to cache any kind of data. It is thread safe and can be used in a concurrent environment.
The prymary use case of teh package is to solve react useEffect issues where the useEffect hook is called multiple times therefore making multiple requests to the server. This package can be used to cache the response of the request and return the cached response instead of making a new request avoid same request multiple times.

It caches the method, the url, the body and some headers (if specified) and returns the cached response if the request is the same.

It currently supports those types of caching:
- In memory caching
- File caching (not implemented yet)

Usage:
```go
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/Polo4444/httputils/caching"
    "github.com/gorilla/mux"

    // create http router

    router := mux.NewRouter()
    router.use(NewCaching([]string{"token"}).Middleware)
)
```