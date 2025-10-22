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
    handler.RegisterTaskRoutes(r)
    
    r.Run() // default :8080

}