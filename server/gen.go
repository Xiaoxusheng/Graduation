package main

//go:generate   go mod tidy
//go:generate   go build -o server -ldflags "-s -w"
//go:generate   ls
//go:generate   nohup ./server   &
