package server

import (
    "log/slog"
    "github.com/gin-gonic/gin"
    "task-app/internal/handler"
    "task-app/internal/repository"
)

func Start() {

    slog.Info("Starting task-app server", slog.String("component", "server"))

    repository.Connect()
    defer repository.CloseDB()

    r := gin.New()
    r.Use(gin.Recovery())
    handler.RegisterTaskRoutes(r)

    slog.Info("Listening and serving HTTP on: 8080")
    if err:= r.Run(); err != nil {
        slog.Error("Server failed to start", "error", err)
    }

}