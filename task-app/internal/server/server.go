package server

import (
    "github.com/gin-gonic/gin"
    "task-app/internal/handler"
    "task-app/internal/repository"
)

func Start() {

    repository.Connect()
    defer repository.CloseDB()

    r := gin.Default()

    //r.LoadHTMLGlob("cmd/templates/*.html")

    // Register routes
    //handler.RegisterIndexRoutes(r)
    handler.RegisterTaskRoutes(r)
    
    r.Run() // default :8080
}