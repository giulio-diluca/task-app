package server

import (
    "github.com/gin-gonic/gin"
    "task-app/internal/handler"
)

func Start() {
    r := gin.Default()

    // Register routes
    handler.RegisterTaskRoutes(r)

    r.Run() // default :8080
}