package parser

import (
	"net/http"
	"strings"

	"Proxy/pkg/domain/models"
)

func Request(r *http.Request) models.ParsedRequest {
	parsedReq := models.ParsedRequest{
		Method:    r.Method,
		Path:      r.URL.Path,
		GetParams: QueryParams(r.URL.Query()),
		Headers:   Headers(r.Header),
		Cookies:   Cookies(r.Cookies()),
	}

	if r.Method == http.MethodPost && strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err == nil {
			parsedReq.PostParams = FormParams(r.PostForm)
		}
	}

	return parsedReq
}

func Response(resp *http.Response, body string) models.ParsedResponse {
	return models.ParsedResponse{
		Code:    resp.StatusCode,
		Message: resp.Status,
		Headers: Headers(resp.Header),
		Body:    body,
	}
}

func QueryParams(params map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, values := range params {
		result[key] = strings.Join(values, ", ")
	}
	return result
}

func Headers(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		result[key] = strings.Join(values, ", ")
	}
	return result
}

func Cookies(cookies []*http.Cookie) map[string]string {
	result := make(map[string]string)
	for _, cookie := range cookies {
		result[cookie.Name] = cookie.Value
	}
	return result
}

func FormParams(form map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, values := range form {
		result[key] = strings.Join(values, ", ")
	}
	return result
}
