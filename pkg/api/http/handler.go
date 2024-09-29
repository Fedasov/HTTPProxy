package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"

	"Proxy/pkg/domain/models"
	"Proxy/pkg/repository/mongodb"
)

var payloads = []string{
	";cat /etc/passwd;",
	"|cat /etc/passwd|",
	"`cat /etc/passwd`",
}

type Handler struct {
	Repo *mongodb.Repo
}

func NewHandler(repo *mongodb.Repo) *Handler {
	return &Handler{
		Repo: repo,
	}
}

// HandleGetAllRequests
// @Summary Get all requests
// @Description Возвращает список всех запросов
// @Tags requests
// @Produce json
// @Success 200 {array} models.Request
// @Failure 500 {string} string "Failed to fetch requests"
// @Router /api/v1/requests [get]
func (h *Handler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.Repo.GetAllRequests(context.TODO())
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

// HandleGetRequestByID
// @Summary Get request by ID
// @Description Возвращает запром по ID
// @Tags requests
// @Param id path string true "Request ID"
// @Produce json
// @Success 200 {object} models.Request
// @Failure 400 {string} string "Invalid request ID"
// @Failure 404 {string} string "Request not found"
// @Router /api/v1/requests/{id} [get]
func (h *Handler) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	request, err := h.Repo.GetRequestByID(context.TODO(), id)
	if err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(request)
}

// HandleRepeatRequest
// @Summary Repeat a request by ID
// @Description Повторно отправляет запрос и возвращает результат
// @Tags requests
// @Param id path string true "Request ID"
// @Produce json
// @Success 200 {object} models.ParsedResponse
// @Failure 400 {string} string "Invalid request ID"
// @Failure 404 {string} string "Request not found"
// @Failure 500 {string} string "Failed to repeat request"
// @Router /api/v1/repeat/{id} [post]
func (h *Handler) RepeatRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	reqResp, err := h.Repo.GetRequestByID(context.TODO(), id)
	if err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	res, err := execute(reqResp)
	if err != nil {
		http.Error(w, "Failed to execute request", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func execute(request *models.Request) (string, error) {
	var curlCommand bytes.Buffer

	curlCommand.WriteString("curl -x http://127.0.0.1:8080 ")
	if request.Request.Method != "CONNECT" {
		curlCommand.WriteString("-X ")
		curlCommand.WriteString(request.Request.Method)
	}

	for key, value := range request.Request.GetParams {
		curlCommand.WriteString(fmt.Sprintf(" -G --data-urlencode \"%s=%s\"", key, value))
	}

	for key, value := range request.Request.PostParams {
		curlCommand.WriteString(fmt.Sprintf(" -d \"%s=%s\"", key, value))
	}

	for key, value := range request.Request.Headers {
		curlCommand.WriteString(fmt.Sprintf(" -H \"%s: %s\"", key, value))
	}

	for key, value := range request.Request.Cookies {
		curlCommand.WriteString(fmt.Sprintf(" --cookie \"%s=%s\"", key, value))
	}

	curlCommand.WriteString(" " + request.Request.Path)

	s := curlCommand.String()
	cmd := exec.Command("bash", "-c", s)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении команды curl: %v, вывод: %s", err, out)
	}

	res := strings.Split(string(out), "<html>")
	result := strings.Join(res[len(res)-1:], "<html>")
	return "<html>" + result, nil
}
