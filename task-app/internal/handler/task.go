package handler

import (
    "github.com/gin-gonic/gin"
    "task-app/internal/service"
)

func RegisterTaskRoutes(r *gin.Engine) {
    group := r.Group("/tasks")
    group.GET("", getAllTasks)
    group.GET("/:id", getTaskByID)
}

func getAllTasks(c *gin.Context) {
    tasks := service.GetAllTasks()
    c.JSON(200, tasks)
}

func getTaskByID(c *gin.Context) {
    id := c.Param("id")
    task, err := service.GetTaskByID(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "Task not found"})
        return
    }
    c.JSON(200, task)
}