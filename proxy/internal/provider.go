package internal

import "net/http"

type Handler interface {
	Proxy(w http.ResponseWriter, r *http.Request)
}
