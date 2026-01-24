package repository

import (
	"fmt"
    "strconv"
    "database/sql"
    "sync"
    "log"
    "time"
    "github.com/go-sql-driver/mysql"
    "github.com/spf13/viper"
	"task-app/internal/model"
)

var (
    db   *sql.DB
    once sync.Once
)

func loadConfig() (model.DBConfig, error) {
	var config model.DBConfig
	
	viper.SetConfigName("db_config")
	viper.SetConfigType("yaml")
    viper.AddConfigPath(".")      
	viper.AddConfigPath("..")

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.UnmarshalKey("database", &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal database config: %w", err)
	}

	return config, nil
}

func Connect() {
	once.Do(func() {
		cfgData, err := loadConfig()
		if err != nil {
			log.Fatalf("FATAL: Could not load database configuration: %v", err)
		}

        // Build the MySQL configuration from config data
		cfg := mysql.NewConfig()
		cfg.User = cfgData.User
		cfg.Passwd = cfgData.Passwd
		cfg.Net = cfgData.Net
		cfg.Addr = cfgData.Addr
		cfg.DBName = cfgData.DBName
        
        // MySQL Connection retry
		for i := 1; i <= cfgData.MaxConnectRetries; i++ {
            // 1. Open the database handle
            var openErr error
            db, openErr = sql.Open("mysql", cfg.FormatDSN())
            if openErr != nil {
                log.Fatalf("FATAL: Error opening database handle: %v", openErr)
            }

			pingErr := db.Ping()
			
			if pingErr == nil {
				db.SetMaxOpenConns(cfgData.MaxOpenConns)
				db.SetMaxIdleConns(cfgData.MaxIdleConns)
				fmt.Println("Database Connected Successfully! (Pool established)")
				return
			}

            db.Close()
            
			log.Printf("WARNING: Database connection failed (Attempt %d/%d): %v", i, cfgData.MaxConnectRetries, pingErr)

			if i < cfgData.MaxConnectRetries {
				log.Printf("Retrying connection in %d seconds...", time.Duration(cfgData.RetryDelaySeconds))
				time.Sleep(time.Duration(cfgData.RetryDelaySeconds) * time.Second)
			}
		}
        
		log.Fatalf("FATAL: Failed to connect to database after %d retries.", cfgData.MaxConnectRetries)
	})
}

func CloseDB() {
    if db != nil {
    db.Close()
    fmt.Println("Database connection pool closed.")
    }
}

func FindAllTasks() []model.Task {
	if db == nil {
		log.Print("ERROR: Database connection not initialized when FindAllTasks was called.")
		return nil
	}

	var tasks = []model.Task{}
	
	rows, err := db.Query("SELECT * FROM task_app")
	if err != nil {
		log.Printf("ERROR: Query failed: %v", err)
		return nil
	}

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
			log.Printf("ERROR: Row scan failed: %v", err)
			return nil
		}
		tasks = append(tasks, task)
	}
	
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Rows iteration error: %v", err)
		return nil
	}
	
	return tasks
}

func FindTaskByID(id string) (model.Task, error) {

    var task = model.Task{}

	if db == nil {
        errorMessage := "ERROR: Database connection not initialized when FindTaskByID was called."
		log.Print(errorMessage)
        return task, fmt.Errorf("%s", errorMessage)
	}
	
	row := db.QueryRow("SELECT ID, Title, Description FROM task_app WHERE ID = ?", id)

    err :=  row.Scan(&task.ID, &task.Title, &task.Description)

	if err != nil {
        if err == sql.ErrNoRows {
            return task, fmt.Errorf("task with ID %s not found", id)
        }
        errorMessage := fmt.Sprintf("ERROR: Failed to scan single row: %v", err)
		log.Printf("%s", errorMessage)
		return task, fmt.Errorf("%s", errorMessage)
	}

    return task, nil
}

func CreateTask(newTask model.Task) (model.Task, error) {

    if db == nil {
        errorMessage := "ERROR: Database connection not initialized when CreateTask was called."
        log.Print(errorMessage)
        return model.Task{}, fmt.Errorf("%s", errorMessage)
    }

    stmt := "INSERT INTO task_app (Title, Description) VALUES (?, ?)"
    
    result, err := db.Exec(stmt, newTask.Title, newTask.Description)

    if err != nil {
        errorMessage := fmt.Sprintf("ERROR: Failed to insert new task: %v", err)
        log.Printf("%s", errorMessage)
        return model.Task{}, fmt.Errorf("%s", errorMessage)
    }

    lastID, err := result.LastInsertId()
    
    if err != nil {
        errorMessage := fmt.Sprintf("CRITICAL ERROR: MySQL LastInsertId failed: %v", err)
        log.Printf("%s", errorMessage)
        return model.Task{}, fmt.Errorf("%s", errorMessage) 
    }

    newTask.ID = strconv.FormatInt(lastID, 10)
    
    log.Printf("DEBUG: Final Task before return: %+v", newTask)
    
    return newTask, nil
}

func UpdateTaskByID(id string, updateTask model.Task) (model.Task, error) {

    if db == nil {
        errorMessage := "ERROR: Database connection not initialized when UpdateTaskByID was called."
        log.Print(errorMessage)
        return model.Task{}, fmt.Errorf("%s", errorMessage) 
    }

    stmt := "UPDATE task_app SET Title = ?, Description = ? WHERE ID = ?"
    result, err := db.Exec(stmt, updateTask.Title, updateTask.Description, id)

    if err != nil {
        errorMessage := fmt.Sprintf("ERROR: Failed to execute update for ID %s: %v", id, err)
        log.Printf("%s", errorMessage)
        return model.Task{}, fmt.Errorf("%s", errorMessage)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        errorMessage := fmt.Sprintf("WARNING: Update succeeded, but could not check rows affected for ID %s: %v", id, err)
        log.Print(errorMessage)
        return model.Task{}, fmt.Errorf("%s", errorMessage)
    }

    if rowsAffected == 0 {
        _, existsErr := FindTaskByID(id) 
        
        if existsErr == sql.ErrNoRows {
            errorMessage := fmt.Sprintf("Task with ID %s not found for update", id)
            return model.Task{}, fmt.Errorf("%s", errorMessage)
            
        } else if existsErr != nil {
            errorMessage := fmt.Sprintf("Error during post-update check: %v", existsErr)
            return model.Task{}, fmt.Errorf("%s", errorMessage)
        }
    }

    updateTask.ID = id
    return updateTask, nil
}

func DeleteTaskByID(id string) (string, error) {

    if db == nil {
        errorMessage := "ERROR: Database connection not initialized when DeleteTaskByID was called."
        log.Print(errorMessage)
        return "", fmt.Errorf("%s", errorMessage)
    }

    stmt := "DELETE FROM task_app WHERE ID = ?"
    
    result, err := db.Exec(stmt, id)

    if err != nil {
        errorMessage := fmt.Sprintf("ERROR: Failed to delete task ID %s: %v", id, err)
        log.Printf("%s", errorMessage)
        return "", fmt.Errorf("%s", errorMessage)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("WARNING: Delete succeeded, but could not check rows affected for ID %s: %v", id, err)
    }

    if rowsAffected == 0 {
        return "", fmt.Errorf("task with ID %s not found", id)
    }

    return id, nil
}
