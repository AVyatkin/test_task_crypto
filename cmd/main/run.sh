#!/bin/bash

exec go get -v "github.com/go-sql-driver/mysql" &
wait
exec go build main.go scheduler.go storage.go crypto.go server.go &
wait
exec ./main
