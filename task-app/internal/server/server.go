package server

import (
    "github.com/gin-gonic/gin"
    "task-app/internal/handler"
    "task-app/internal/repository"
)

func Start() {

    // 1. Initialize the DB Connection Pool ONCE at startup
    repository.Connect()

    // 2. Defer closing the pool until the application exits
    defer repository.CloseDB()

    r := gin.Default()

    //r.LoadHTMLGlob("cmd/templates/*.html")

    // Register routes
    //handler.RegisterIndexRoutes(r)
    handler.RegisterTaskRoutes(r)
    
    r.Run() // default :8080
}