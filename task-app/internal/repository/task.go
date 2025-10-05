package repository

import (
	"fmt"
    "strconv"
	"task-app/internal/model"
)

var tasks = []model.Task{
    {ID: "1", Title: "Task 1", Description: "D1"},
    {ID: "2", Title: "Task 2", Description: "D2"},
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


func CreateTask(newTask model.Task) model.Task {
    
    maxID := 0
    for _, t := range tasks {
        // Attempt to convert the string ID to an integer
        if currentID, err := strconv.Atoi(t.ID); err == nil {
            if currentID > maxID {
                maxID = currentID
            }
        }
    }

    newID := maxID + 1 
    
    newTask.ID = strconv.Itoa(newID) 

    tasks = append(tasks, newTask)
    
    return newTask
}


func UpdateTaskByID(id string, updatedTaskData model.Task) (model.Task, error) {

    var index int
    found := false

    for i, t := range tasks {
        if t.ID == id {
            index = i
            found = true
            break
        }
    }

    if !found {
        return model.Task{}, fmt.Errorf("task with ID %s not found", id)
    }

    updatedTaskData.ID = id 
    tasks[index] = updatedTaskData 
    return tasks[index], nil

}

func DeleteTaskByID(id string) (string, error) {

    index := -1
    for i, t := range tasks {
        if t.ID == id {
            index = i
            break
        }
    }

    if index == -1 {
        return "", fmt.Errorf("task with ID %s not found", id)
    }

    tasks = append(tasks[:index], tasks[index+1:]...)

    for i := range tasks {
        newID := i + 1
        tasks[i].ID = strconv.Itoa(newID) 
    }

    return id, nil
}

