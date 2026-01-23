#!/bin/bash

BASE_URL_1="http://localhost:8080/tasks"
BASE_URL_2="http://localhost:8081/tasks"
MODE="SEQUENTIAL"

if [[ "$1" == "-p" || "$1" == "--parallel" ]]; then
    MODE="PARALLEL"
fi

test_instance_sequence() {
    local URL=$1
    local PORT=$2
    
    echo "--- Running Full Sequence on PORT: $PORT"

    local TITLE="task t.$PORT"
    local DESCRIPTION="task d.$PORT"
    curl -s -X POST "$URL" -d "{\"Title\": \"$TITLE\", \"Description\": \"$DESCRIPTION\"}" -H "Content-Type: application/json" > /dev/null

    echo "    GET $PORT/tasks result:"
    curl -s -X GET "$URL"
    echo ""
    
    curl -s -X PUT "$URL/1" -d "{\"Title\": \"$TITLE updated\", \"Description\": \"$DESCRIPTION updated\"}" -H "Content-Type: application/json" > /dev/null

    curl -s -X DELETE "$URL/1" > /dev/null
    
    echo "    GET $PORT/tasks final result:"
    curl -s -X GET "$URL"
    echo ""
}

test_parallel_conflict() {
    echo "======================================="
    echo "   SIMULATING PARALLEL POST REQUESTS   "
    echo "======================================="

    post_task_concurrently() {
        local URL=$1
        local PORT=$2
        
        echo "Sending POST to port $PORT..."
        curl -s -X POST "$URL" -d '{"Title": "Concurrent Task", "Description": "From Port '$PORT'"}' \
          -H "Content-Type: application/json" -w "Port $PORT HTTP Status: %{http_code}\n"
    }

    post_task_concurrently $BASE_URL_1 8080 & 
    PID_1=$!

    post_task_concurrently $BASE_URL_2 8081 &
    PID_2=$!

    wait $PID_1
    wait $PID_2

    echo "---------------------------------------"
    echo "✅ Both POSTs complete. Checking final state (GET from 8080):"

    curl -s -X GET $BASE_URL_1
    echo ""
    
}

if [ "$MODE" == "PARALLEL" ]; then
    echo "***************************************"
    echo "*** RUNNING IN PARALLEL (CONCURRENT) MODE ***"
    echo "***************************************"
    test_parallel_conflict
else
    echo "***************************************"
    echo "*** RUNNING IN SEQUENTIAL MODE ***"
    echo "***************************************"
    test_instance_sequence $BASE_URL_1 8080
    echo "======================================="
    test_instance_sequence $BASE_URL_2 8081
fi