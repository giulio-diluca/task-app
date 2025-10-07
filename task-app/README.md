```
cd ~/projects/golab/task-app
docker build . -t task-app:latest
docker run .d -p 8080:8080 --name task-app task-app:latest
```