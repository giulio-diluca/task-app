# APP

```
cd ~/projects/golab/task-app
docker build . -t task-app:latest
docker run .d -p 8080:8080 --name task-app task-app:latest

```

# DB

```
docker run --name task-app-sql-db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=12345 -e MYSQL_DATABASE=task_app -v task-app-sql-db-vm:/var/lib/mysql -d mysql:latest
docker exec -it task-app-sql-db mysql -u root -p

mysql> show databases;
mysql > connect task_app;
mysql> create table task_app (ID integer key auto_increment,Title varchar(20),Description varchar(50));
mysql> show tables;
mysql > insert into task_app values ('1', 'task t.1', 'task d.1');
mysql > select * from task_app;
mysql > truncate table task_app;

docker rm -f task-app-sql-db
docker volume rm -f task-app-sql-db-vm
```