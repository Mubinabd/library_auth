package main

import (
	"github.com/Mubinabd/library_auth/config"
	"github.com/Mubinabd/library_auth/pkg/app"
)

func main() {
	cfg := config.Load()
	app.Run(&cfg)
}
