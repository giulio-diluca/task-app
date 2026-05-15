package repository

import (
	"context"
    "database/sql"
    "fmt"
    "log/slog"
    "os"
    "strconv"
    "sync"
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
		return config, err
	}
	if err := viper.UnmarshalKey("database", &config); err != nil {
		return config, nil
	}
	return config, nil
}

func Connect() {
	once.Do(func() {
		cfgData, err := loadConfig()
		if err != nil {
			slog.Error("failed to load database configuration", "error", err)
            os.Exit(1)
		}
		cfg := mysql.NewConfig()
		cfg.User = cfgData.User
		cfg.Passwd = cfgData.Passwd
		cfg.Net = cfgData.Net
		cfg.Addr = cfgData.Addr
		cfg.DBName = cfgData.DBName
		for i := 1; i <= cfgData.MaxConnectRetries; i++ {
            db, err = sql.Open("mysql", cfg.FormatDSN())
            if err != nil {
                slog.Error("error opening database handle", "error", err)
                os.Exit(1)
            }
            if err := db.Ping(); err == nil {
                db.SetMaxOpenConns(cfgData.MaxOpenConns)
                db.SetMaxIdleConns(cfgData.MaxIdleConns)
                slog.Info("database connected successfully", "status", "pool_established", "addr", cfgData.Addr)
                return
            }
            db.Close()
            slog.Warn("database connection failed", "attempt", i, "max", cfgData.MaxConnectRetries)
            if i < cfgData.MaxConnectRetries {
                time.Sleep(time.Duration(cfgData.RetryDelaySeconds) * time.Second)
            }
        }

        slog.Error("failed to connect to database after retries", "attempts", cfgData.MaxConnectRetries)
        os.Exit(1)
    })
}

func CloseDB() {
    if db != nil {
        db.Close()
        slog.Info("database connection pool closed")
    }
}

func trackQuery(start time.Time, queryName string) {
    elapsed := time.Since(start)
    slog.Debug("database queri executed",
        "query", queryName,
        "latency_ms", elapsed.Milliseconds(),
    )
}

// stopped here 15/05/2026
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
