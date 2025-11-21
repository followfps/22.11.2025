package main

import (
    "HealthCheck/internal/handler"
    "HealthCheck/internal/router"
    "HealthCheck/pkg/local_storage"
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    storage := local_storage.NewLocalStorage()
    persistPath := "storage.json"
    _ = storage.LoadFromDisk(persistPath)

    h := handler.NewHandler(storage, persistPath)
    r := router.NewRouter(h)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = srv.Shutdown(ctx)
    _ = storage.SaveToDisk(persistPath)
}
