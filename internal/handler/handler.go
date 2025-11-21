package handler

import (
	"HealthCheck/pkg/local_storage"
	"HealthCheck/pkg/pdf"
	"HealthCheck/pkg/url_check"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	storage     *local_storage.TempLinksFormRequest
	persistPath string
}

func NewHandler(storage *local_storage.TempLinksFormRequest, persistPath string) *Handler {
	return &Handler{storage: storage, persistPath: persistPath}
}

type urlsRequest struct {
	Url  string   `json:"url"`
	Urls []string `json:"urls"`
}

type urlsResponse struct {
	Id       uint64          `json:"id"`
	Statuses map[string]bool `json:"statuses"`
}

func (h *Handler) HandleUrlsCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req urlsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	var urls []string
	if req.Url != "" {
		urls = append(urls, req.Url)
	}
	urls = append(urls, req.Urls...)
	m := map[string]struct{}{}
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		m[u] = struct{}{}
	}
	urls = urls[:0]
	for u := range m {
		urls = append(urls, u)
	}
	if len(urls) == 0 {
		http.Error(w, "no urls", http.StatusBadRequest)
		return
	}

	id := h.storage.AddUrl(urls)
	statuses := make(map[string]bool, len(urls))
	for _, u := range urls {
		statuses[u] = url_check.UrlHealthCheck(u)
	}
	h.storage.UpdateAllStatuses(id, statuses)
	_ = h.storage.SaveToDisk(h.persistPath)

	resp := urlsResponse{Id: id, Statuses: statuses}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	data, ok := h.storage.GetUrls(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     data.Id,
		"urls":   data.Urls,
		"status": data.UrlsStatus,
	})
}

type reportRequest struct {
	Ids []uint64 `json:"ids"`
}

func (h *Handler) HandlePdfReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req reportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Ids) == 0 {
		http.Error(w, "invalid json or empty ids", http.StatusBadRequest)
		return
	}
	var lines []string
	for _, id := range req.Ids {
		data, ok := h.storage.GetUrls(id)
		if !ok {
			lines = append(lines, fmt.Sprintf("ID %d: not found", id))
			continue
		}
		lines = append(lines, fmt.Sprintf("ID %d:", id))
		for _, u := range data.Urls {
			st := data.UrlsStatus[u]
			lines = append(lines, fmt.Sprintf("  %s - %s", u, boolToWord(st)))
		}
	}
	pdfBytes := pdf.Generate("HealthCheck Report", lines)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"report.pdf\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}

func boolToWord(b bool) string {
	if b {
		return "доступен"
	}
	return "не доступен"
}
