package server

import (
    "github.com/gin-gonic/gin"
    "task-app/internal/handler"
)

func Start() {
    r := gin.Default()

    r.LoadHTMLGlob("cmd/templates/*.html")

    // Register routes
    handler.RegisterTaskRoutes(r)
    handler.RegisterIndexRoutes(r)

    r.Run() // default :8080
}