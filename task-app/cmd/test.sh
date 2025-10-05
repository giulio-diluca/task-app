#!/bin/bash

#{ID: "1", Title: "Task 1", Description: "D1"}
#{ID: "2", Title: "Task 2", Description: "D2"},
echo "GET"
curl -X GET http://localhost:8080/tasks
echo ""
echo "---------------------------------------"

#{ID: "1", Title: "Task 1", Description: "D1"}
#{ID: "2", Title: "Task 2", Description: "D2"},
#{ID: "3", Title: "Task 3", Description: "D3"},
echo "POST -> GET"
curl -X POST http://localhost:8080/tasks -d '{"Title": "Task 3", "Description": "D3"}'
echo ""
curl -X GET http://localhost:8080/tasks
echo ""
echo "---------------------------------------"

#{ID: "1", Title: "Task 1", Description: "D1"}
#{ID: "2", Title: "Task 2", Description: "D4"},
#{ID: "3", Title: "Task 3", Description: "D3"},
echo "PUT -> GET"
curl -X PUT http://localhost:8080/tasks/2 -d '{"Title": "Task 2", "Description": "D4"}'
echo ""
curl -X GET http://localhost:8080/tasks
echo ""
echo "---------------------------------------"

#{ID: "1", Title: "Task 2", Description: "D4"},
#{ID: "2", Title: "Task 3", Description: "D3"},
echo "DELETE -> GET"
curl -X DELETE http://localhost:8080/tasks/1
echo ""
curl -X GET http://localhost:8080/tasks
