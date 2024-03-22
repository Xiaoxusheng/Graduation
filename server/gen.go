package main

//go:generate   go mod tidy
//go:generate   go build -o "server" -ldflages "-s -w"
//go:generate   ls
//go:generate   nohup ./server   &
