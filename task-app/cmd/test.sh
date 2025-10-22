#!/bin/bash

#{ID: "1", Title: "task t.1", Description: "task d.1"}
echo "GET"
curl -X GET http://localhost:8080/tasks
echo ""
echo "---------------------------------------"

#{ID: "1", Title: "task t.1", Description: "task d.1"}
#{ID: "1", Title: "task t.2", Description: "task d.2"}
echo "POST -> GET"
curl -X POST http://localhost:8080/tasks -d '{"Title": "task t.2", "Description": "task d.2"}'
echo ""
curl -X GET http://localhost:8080/tasks
echo ""
echo "---------------------------------------"

#{ID: "1", Title: "task t.1", Description: "task d.1"}
#{ID: "1", Title: "task t.2", Description: "task d.3"}
echo "PUT -> GET"
curl -X PUT http://localhost:8080/tasks/2 -d '{"Title": "task t.2", "Description": "task d.3"}'
echo ""
curl -X GET http://localhost:8080/tasks
echo ""
echo "---------------------------------------"

#{ID: "2", Title: "task t.2", Description: "task d.3"}
echo "DELETE -> GET"
curl -X DELETE http://localhost:8080/tasks/1
echo ""
curl -X GET http://localhost:8080/tasks
