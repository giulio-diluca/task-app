package repository

import (
	"fmt"
	"task-app/internal/model"
)

var tasks = []model.Task{
    {ID: "1", Title: "Task 1", Description: "123"},
    {ID: "2", Title: "Task 2", Description: "321"},
}

func FindAllTasks() []model.Task {
    return tasks
}

func FindTaskByID(id string) (model.Task, error) {
    for _, u := range tasks {
        if u.ID == id {
            return u, nil
        }
    }
    return model.Task{}, fmt.Errorf("task not found")
}

func CreateTask(newTask model.Task) {

    tasks = append(tasks, newTask)

}

func DeleteTaskByID(id string) (string, error) {

    var index int
    for i, u := range tasks {
        if u.ID == id {
            index = i
            break
        }
    }

    tasks = append(tasks[:index], tasks[index+1:]...)
    return id, fmt.Errorf("task not found")

}

