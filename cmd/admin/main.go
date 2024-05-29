package main

import (
	_ "go.uber.org/automaxprocs"
	"ktserver/internal/admin"
)

// go build -ldflags "-X main.Version=x.y.z"
func main() {
	admin.NewApp()
}
