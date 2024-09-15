package server

import (
	"crypto/tls"
	"fmt"
	proxy "main/internal/pkg/http"
	"net/http"
)

func Init(port string, handler proxy.ProxyHandler) *http.Server {
	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handler.HandleTunneling(w, r)
			} else {
				handler.HandleHTTP(w, r)
			}
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	return server
}
