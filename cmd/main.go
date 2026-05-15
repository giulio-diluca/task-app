package main

import (
    "task-app/internal/pkg/logger"
    "task-app/internal/server"
)

func main() {
    logger.Init()
    server.Start()
}