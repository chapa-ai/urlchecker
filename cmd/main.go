package main

import (
	"github.com/chapa-ai/urlchecker/config"
	"github.com/chapa-ai/urlchecker/internal/app"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	cfg, err := config.LoadConfig(filepath.Join(dir, "config", "configs.json"))
	if err != nil {
		panic(err)
	}

	app.NewApp(cfg).Serve()
}
