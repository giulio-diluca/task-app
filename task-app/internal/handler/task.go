package handler

import (
	"fmt"
	"task-app/internal/model"
	"task-app/internal/service"
	"github.com/gin-gonic/gin"
)

// func RegisterIndexRoutes(r *gin.Engine) {
//     group := r.Group("")
//     group.GET("", indexHtml)
// }

func RegisterTaskRoutes(r *gin.Engine) {
    group := r.Group("/tasks")
    group.GET("", getAllTasks)
    group.GET("/:id", getTaskByID)
    group.POST("", postTask)
    group.DELETE("/:id", deleteTaskByID)
}

// func indexHtml(c *gin.Context){
//     tasks := service.GetAllTasks()
//     c.HTML(200, "index.html", gin.H{
//         "Title": "Task Service",
//         "Tasks": tasks,
// 	})
// }

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

    service.PostTask(newTask)
    c.IndentedJSON(200, newTask)
}

func deleteTaskByID(c *gin.Context) {
    id := c.Param("id")
    task, err := service.GetTaskByID(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "Task not found"})
    }

    service.DeleteTaskByID(id)
    fmt.Println(task)
    c.IndentedJSON(200, id)
}
