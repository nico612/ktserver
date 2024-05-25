package main

import (
	_ "go.uber.org/automaxprocs"
	"ktserver/internal/admin"
)

// go build -ldflags "-X main.Version=x.y.z"
func main() {
	app, clean := admin.NewApp("admin-server")
	defer clean()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
