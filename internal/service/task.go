package service

import (
    "fmt"
	"task-app/internal/model"
	"task-app/internal/repository"
)

func GetAllTasks() []model.Task {
    return repository.FindAllTasks()
}

func GetTaskByID(id string) (model.Task, error) {
    task, err := repository.FindTaskByID(id) 
    if err != nil {
        return model.Task{}, fmt.Errorf("service getting task %s: %w", id, err)
    }
    return task, nil
}

func PostTask(newTask model.Task) (model.Task, error) {
    createdTask, err := repository.CreateTask(newTask)
    if err != nil {
        return model.Task{}, fmt.Errorf("service creating task: %w", err)
    }
    return createdTask, nil
}

func UpdateTaskByID(id string, updateTask model.Task) (model.Task, error) {
    updatedTask, err := repository.UpdateTaskByID(id, updateTask)
    if err != nil {
        return model.Task{}, fmt.Errorf("service updating task %s, %w", id, err)
    }
    return updatedTask, nil
}

func DeleteTaskByID(id string) (string, error) {
    deletedID, err := repository.DeleteTaskByID(id)
    if err != nil {
        return "", fmt.Errorf("service deleting task %s, %w", id, err)
    }
    return deletedID, nil
}
