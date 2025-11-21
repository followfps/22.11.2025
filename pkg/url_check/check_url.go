package url_check

import (
    "net/http"
    "time"
)

func UrlHealthCheck(url string) bool {
    client := &http.Client{Timeout: 5 * time.Second}

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return false
    }

    resp, err := client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode >= 200 && resp.StatusCode < 400
}
