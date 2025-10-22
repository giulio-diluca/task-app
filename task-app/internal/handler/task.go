package handler

import (
	"fmt"
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
        c.JSON(404, gin.H{"error": "Task not found"})
    }
    c.JSON(200, task)

}

func postTask(c *gin.Context) {

    var newTask model.Task
    
    if err := c.ShouldBindJSON(&newTask); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    createdTask, err := service.PostTask(newTask)

    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create task"})
        return
    }

    c.JSON(201, createdTask)

}


func updateTaskByID(c *gin.Context) {
    id := c.Param("id")

    var updateTaskData model.Task 

    if err := c.ShouldBindJSON(&updateTaskData); err != nil {
        c.JSON(400, gin.H{"error": "Invalid JSON format: " + err.Error()})
        return
    }

    updateTaskData.ID = id

    updatedTask, err := service.UpdateTaskByID(id, updateTaskData)
    
    if err != nil {
        c.JSON(404, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, updatedTask)

}

func deleteTaskByID(c *gin.Context) {
    id := c.Param("id")
    task, err := service.GetTaskByID(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "Task not found"})
    }

    service.DeleteTaskByID(id)
    fmt.Println(task)
    c.JSON(200, id)

}
