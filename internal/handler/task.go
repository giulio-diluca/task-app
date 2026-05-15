package handler

import (
    "log/slog"
    "net/http"
	"task-app/internal/model"
	"task-app/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(r *gin.Engine) {
    group := r.Group("/tasks")
    group.GET("", getAllTasks)
    group.GET("/:id", getTaskByID)
    group.POST("", postTask)
    group.PUT("/:id", updateTaskByID)
    group.DELETE("/:id", deleteTaskByID)
}

func getAllTasks(c *gin.Context) {
    tasks := service.GetAllTasks()
    c.JSON(200, tasks)
}

func getTaskByID(c *gin.Context) {
    id := c.Param("id")
    task, err := service.GetTaskByID(id)
    if err != nil {
        slog.Warn("task not found", "id", id, "error", err.Error())
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    c.JSON(http.StatusOK, task)
}

func postTask(c *gin.Context) {
    var newTask model.Task
    if err := c.ShouldBindJSON(&newTask); err != nil {
        slog.Warn("invalid request body", "error", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    createdTask, err := service.PostTask(newTask)
    if err != nil {
        slog.Error("failed to create task", "error", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }
    c.JSON(http.StatusCreated, createdTask)

}


func updateTaskByID(c *gin.Context) {
    id := c.Param("id")
    var updateTaskData model.Task 
    if err := c.ShouldBindJSON(&updateTaskData); err != nil {
        slog.Warn("invalid update data", "id", id, "error", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
        return
    }
    updateTaskData.ID = id
    updatedTask, err := service.UpdateTaskByID(id, updateTaskData)
    if err != nil {
        slog.Warn("failed to update task", "id", id, "error", err.Error())
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, updatedTask)
}

func deleteTaskByID(c *gin.Context) {
    id := c.Param("id")
    _, err := service.GetTaskByID(id)
    if err != nil {
        slog.Warn("attempted to delete non-existent task", "id", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    deletedId, err := service.DeleteTaskByID(id)
    if err != nil {
        slog.Error("failed to delete task", "id", id, "error", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete task"})
        return
    }

    slog.Info("task deleted successfully", "id", deletedId)
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted", "id": deletedId})

}
