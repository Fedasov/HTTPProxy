package provider

import (
	"fmt"
	"io"
	"net/http"
)

type Handler struct{}

func CreateHandler() *Handler {
	handler := &Handler{}
	return handler
}

func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	// удаляем заголовок Proxy-Connection
	fmt.Println(r)
	r.Header.Del("Proxy-Connection")
	// заменяем путь на относительный
	r.RequestURI = ""

	// низкоуровнево отправляем запрос на r.Host
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	for name, values := range resp.Header {
		for _, v := range values {
			fmt.Println("name: ", name, " Value: ", v)
			w.Header().Add(name, v)
		}
	}

	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}
