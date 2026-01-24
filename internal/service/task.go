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

func PostTask(newTask model.Task) (model.Task, error) {

    updatedTask, err := repository.CreateTask(newTask)

    if err != nil {
        return model.Task{}, err
    }
    
    return updatedTask, nil
}

func UpdateTaskByID(id string, updateTask model.Task) (model.Task, error) {

    updatedTask, err := repository.UpdateTaskByID(id, updateTask)

    if err != nil {
        return model.Task{}, err
    }

    return updatedTask, nil

}

func DeleteTaskByID(id string) (string, error) {
    return repository.DeleteTaskByID(id)
}
