package repository

import (
	"fmt"
    "strconv"
    "database/sql"
    "sync"
    "log"
    "github.com/go-sql-driver/mysql"
	"task-app/internal/model"
)

// The global database handle (connection pool) for the repository package.
var db *sql.DB

// sync.Once is a synchronization primitive used to guarantee that a specific function (usually an expensive initialization step) is executed only one time, even if called simultaneously by multiple goroutines
var once sync.Once

// var tasks = []model.Task{
//     {ID: "1", Title: "Task 1", Description: "D1"},
//     {ID: "2", Title: "Task 2", Description: "D2"},
// }

var tasks = []model.Task{ }

func Connect() {
    
    // sync.Once ensures that the Connect function's body runs exactly once.
    once.Do(func() {
        cfg := mysql.NewConfig()
        cfg.User = "root"
        cfg.Passwd = "12345"
        cfg.Net = "tcp"
        // This uses the HOST's port 3306, which Docker maps to the container.
        cfg.Addr = "127.0.0.1:3306" 
        cfg.DBName = "task_app"

        var err error
        // sql.Open returns a DB handle (connection pool), not a connection.
        db, err = sql.Open("mysql", cfg.FormatDSN())
        if err != nil {
            // A connection error at startup is fatal.
            log.Fatalf("FATAL: Error opening database handle: %v", err)
        }

        // Set connection pool limits for a robust application
        // (Example values - adjust based on your workload and DB server limits)
        db.SetMaxOpenConns(25)
        db.SetMaxIdleConns(10)

        // Verify the connection is actually working
        pingErr := db.Ping()
        if pingErr != nil {
            db.Close() // Close the handle if the ping fails
            log.Fatalf("FATAL: Error pinging database: %v", pingErr)
        }
        fmt.Println("Database Connected Successfully! (Pool established)")
    })
}

func CloseDB() {
    if db != nil {
    db.Close()
    fmt.Println("Database connection pool closed.")
    }
}

func FindAllTasks() []model.Task {
	// A simple check to ensure the connection was initialized
	if db == nil {
		log.Print("ERROR: Database connection not initialized when FindAllTasks was called.")
		return nil
	}

	var tasks = []model.Task{}
	
	// db.Query draws a connection from the pool automatically
	rows, err := db.Query("SELECT * FROM task_app")
	if err != nil {
		log.Printf("ERROR: Query failed: %v", err)
		return nil
	}
	defer rows.Close() // ALWAYS defer rows.Close()

	for rows.Next() {
		var task model.Task
		// Update Scan to match the columns in your table and fields in your model.Task struct
		if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
			log.Printf("ERROR: Row scan failed: %v", err)
			return nil
		}
		tasks = append(tasks, task)
	}
	
	// Check for errors after the loop (important for some drivers)
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Rows iteration error: %v", err)
		return nil
	}
	
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

