package http

import (
	"Proxy/pkg/domain/parser"
	"Proxy/pkg/repository/mongodb"
	"bytes"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func InitDB() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://mongo-container:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Disabling MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB is not available: %v", err)
	}

	collection := client.Database("web").Collection("requests")

	return collection
}

type DataBase struct {
	Repo *mongodb.Repo
}

func (p *DataBase) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleHTTPS(w, r)
	} else {
		p.handleHTTP(w, r)
	}
}

func (p *DataBase) handleHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")

	fullURL := r.URL.String()
	request := parser.Request(r)
	request.Path = fullURL
	response := &http.Response{}

	r.RequestURI = ""

	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	proxyResponse, err := httpClient.Do(r)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		log.Fatalf(err.Error())
	}
	defer proxyResponse.Body.Close()

	response.Header = make(http.Header)
	for header, values := range proxyResponse.Header {
		stringValues := strings.Join(values, ", ")
		w.Header().Set(header, stringValues)
		response.Header.Set(header, stringValues)
	}
	w.WriteHeader(proxyResponse.StatusCode)
	response.StatusCode = proxyResponse.StatusCode

	var buf bytes.Buffer
	mw := io.MultiWriter(&buf, w)
	io.Copy(mw, proxyResponse.Body)

	respBody := buf.String()
	decodedBody := html.UnescapeString(respBody)
	resp := parser.Response(response, "")
	resp.Body = decodedBody
	_, err = p.Repo.AddRequestResponse(r.Context(), request, resp)
	if err != nil {
		log.Printf("Don`t save request")
	}
}

func (p *DataBase) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	fullURL := r.URL.String()
	request := parser.Request(r)
	request.Path = "https:" + strings.Split(fullURL, ":")[0] + "/"
	response := &http.Response{}

	connDest, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf(err.Error())
	}

	response.StatusCode = http.StatusOK
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		response.StatusCode = http.StatusInternalServerError
		log.Fatalf(err.Error())
	}

	connSrc, _, err := hijacker.Hijack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.StatusCode = http.StatusInternalServerError
		log.Fatalf(err.Error())
	}

	go broadcastData(connDest, connSrc, nil)
	body := make(chan string)
	go broadcastData(connSrc, connDest, body)

	resp := parser.Response(response, "")
	resp.Body = <-body
	_, err = p.Repo.AddRequestResponse(r.Context(), request, resp)
	if err != nil {
		log.Printf("Don`t save request")
	}
}

func broadcastData(to io.WriteCloser, from io.ReadCloser, body chan string) {
	defer func() {
		to.Close()
		from.Close()
	}()
	var buf bytes.Buffer
	mw := io.MultiWriter(&buf, to)
	io.Copy(mw, from)
	if body != nil {
		body <- buf.String()
	}
}
