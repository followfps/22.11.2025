package router

import (
    "HealthCheck/internal/handler"
    "net/http"
)

func NewRouter(h *handler.Handler) http.Handler {
    r := http.NewServeMux()

    r.HandleFunc("/urls", h.HandleUrlsCheck)
    r.HandleFunc("/report/pdf", h.HandlePdfReport)
    r.HandleFunc("/status", h.HandleStatus)
    return r
}
