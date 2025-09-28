package service

import (
    "task-app/internal/model"
    "task-app/internal/repository"
)

func GetAllTasks() []model.Task {
    return repository.FindAllTasks()
}

func GetTaskByID(id string) (model.Task, error) {
    return repository.FindTaskByID(id)
}