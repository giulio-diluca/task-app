# Task-App

[![CodeQL](https://github.com/giulio-diluca/task-app/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/giulio-diluca/task-app/actions/workflows/github-code-scanning/codeql)
[![Dependency Graph](https://github.com/giulio-diluca/task-app/actions/workflows/dependabot/update-graph/badge.svg)](https://github.com/giulio-diluca/task-app/actions/workflows/dependabot/update-graph)

A simple Task-App CRUD REST APIs built with **Go** and the **Gin Gonic** framework. 

## 🏗 System Architecture

The project deploys multiple application instances that connect to a single shared MySQL database.

* **Multi-Instance Support**: By default, the environment runs two application containers (`task-app-1` and `task-app-2`) exposed on ports `8080` and `8081`. Each application container have a bind mount volume to the local `db_config.yaml` file as follow `./db_config.yaml:/app/db_config.yaml`
* **Scalability**: You can easily scale the service by adding more instances in the `docker-compose.yaml` file and mapping them to any available ports on your host machine.
* **Infrastructure**: It includes a MySQL database container named `task-app-sql-db`, exposed on port `3306`. For data persistence, a volume `task-app-sql-db-volume` is created, it contain `/var/lib/mysql` MySQL database container path. All containers are in a network with driver `bridge` named `task-app-network`
* **Database Resilience**: The connection logic features a retry mechanism that attempts to connect up to 10 times with a 5-second delay, ensuring the app handles database startup lag gracefully.
* **Clean Architecture**: The code is organized into distinct layers: `handler`, `service`, `model`, `server`, and `repository`.
* **Containerization**: Includes a multi-stage `Dockerfile` for small image sizes and `docker-compose` for easy orchestration.

## 📂 Project Structure
```text
├── cmd/                # Entry point (main.go) and test.sh
├── internal/
│   ├── handler/        # API routes & logic
│   ├── model/          # Structs & DB config
│   ├── repository/     # DB queries (SQL)
│   ├── server/         # Server initialization
│   └── service/        # Business logic
├── .env                # docker-compose env file
├── db_config.yaml      # Database configuration file
├── docker-compose.yaml # Docker orchestration setup
└── Dockerfile          # Multi-stage build file
```


## 🚀 Getting Started

### 1. Software Requirements
* **Docker & Docker Compose**: Necessary for containerizing the application and the MySQL database. (Developed using **Docker Desktop** on **Windows**).
* **Go (Golang)**: Version **1.25.1** or higher is required if you plan to build or run the application outside of Docker.
* **Git**: Required to clone the repository and manage version control.
* **Bash Shell**: Required to execute the `test.sh` script located in the `cmd/` directory.

### 2. Network & Ports
Ensure the following ports are available on your host machine to avoid conflicts:
* **Port 8080**: Used by `task-app-1`.
* **Port 8081**: Used by `task-app-2`.
* **Port 3306**: Used by the `task-app-sql-db` (MySQL).
* **Note**: If you add more instances, ensure the additional ports are also free.

### 3. Configuration Setup
* **Database Config**: A `db_config.yaml` file must exist in the root directory to define the database connection parameters (user, password, address, etc.). In repository there is a `db_config.example.yaml`, rename this properly.
* **Docker**: Ensure the Docker daemon is running so that the `docker-compose up` command can initialize the containers and the bridge network. Make sure also to have `.env` file in root directory, this define variables used by docker-compose, rename it properly.

### 4. Launch with Docker
Use Docker Compose to build and start the entire stack
```
docker-compose up -d --build
```

### 5. Database Setup
```
docker exec -it task-app-sql-db bash

// in db container, write command and insert password
mysql -u root -p

// in mysql command
show databases;
connect task_app;
create table task_app (ID integer key auto_increment,Title varchar(20),Description varchar(50));
```

## 💻 Running Locally ( without Docker )

If you prefer to run the Go application directly on your host machine for development:

### 1. Adjust Configuration
If you want to use a MySQL locally, you need to change `addr` configuration in `db_config.yaml` to reach it on `localhost`
```yaml
database:
  addr: 127.0.0.1:3306  # Change from task-app-sql-db:3306 to local address
  # ... keep other settings the same
```

### 2. Start the Database
If you want to continue MySQL with Docker, you can still do it by not changing `addr` in configuration `db_config.yaml`, make sure to edit `.env` file properly and then execute `docker-compose` command
```bash
docker-compose up -d task-app-sql-db
```

### 3. Run the Go App
Navigate to the project root and execute the following commands:
```bash
# Download dependencies
go mod download
# Run the application
go run cmd/main.go
```

## 📦 GitHub Packages (GHCR) Integration
This project is configured to be hosted on the GitHub Container Registry (GHCR).
Currently you can't publish your own Docker image in the already existing Docker image `ghcr.io/giulio-diluca/task-app:latest` but of course you can build your own image as follow

### 1. Authenticate with GitHub
Before pushing, you must log in using a Personal Access Token (PAT) with `write:packages` permissions:
```bash
echo $CR_PAT | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin
```

### 1.1 Pull from GitHub Packages
If you want to use the already existing image without proceeding with Build and Push, use following command after authentication step
```bash
docker pull ghcr.io/giulio-diluca/task-app:latest
```

### 2. Build and Tag
Build the image using the tag specified in your `docker-compose.yaml`:
```bash
docker build -t ghcr.io/YOUR_GITHUB_USERNAME/task-app:latest .
```

### 3. Push to GitHub Packages
```
docker push ghcr.io/YOUR_GITHUB_USERNAME/task-app:latest
```

## 📡 API Endpoints
The service provides a simple RESTful interface for Task CRUD operations:
Method | Endpoint | Description                    | Request Body ( JSON ) |
------ | -------- | -----------------------------  | --------------------- |
GET	   | /tasks     | Retrieve all tasks | N/A
GET	   | /tasks/:id | Retrieve a specific task by ID | N/A
POST   | /tasks     | Create a new task              | `{"Title": "string", "Description": "string"}`
PUT	   | /tasks/:id	| Update an existing task | `{"Title": "string", "Description": "string"}`
DELETE | /tasks/:id | Delete a task           | N/A

## 🧪 Testing
A dedicated test script cmd/test.sh is included to verify API functionality across both instances.

**Sequential Mode (Default)**: Runs full CRUD verification on port 8080 and then port 8081.
```bash
./cmd/test.sh
```

**Parallel Mode**: Simulates concurrent POST requests to both instances simultaneously using the --parallel or -p flag.
```bash
./cmd/test.sh --parallel
# or
./cmd/test.sh -p
```

## ⚙️ Configuration
Database settings are managed in **db_config.yaml** and loaded via **Viper**.
You can modify these values to suit your environment
```yaml
database:
  user: root                # Database username
  password: 12345           # Database password
  net: tcp                  # Network type (usually tcp)
  addr: task-app-sql-db:3306 # Address: Use MySQL container name
  dbname: task_app          # The name of the database to connect to
  max_open_conns: 25        # Maximum number of open connections to the database
  max_idle_conns: 10        # Maximum number of connections in the idle connection pool
  max_connect_retries: 10   # How many times the app tries to connect before crashing
  retry_delay_seconds: 5    # Time to wait between each connection attempt
```
<br>

**docker-compose** environment variables are managed in **.env** file
You can modify these values to suit your environment
```yaml
MYSQL_ROOT_USER: root      # Database root username
MYSQL_ROOT_PASSWORD: 12345 # Database root password
MYSQL_DATABASE: task_app   # Database name
```