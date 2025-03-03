package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/abel1502/mipt-kp-m-test/internal/app"
)

func main() {
	appName, err := os.Executable()
	if err != nil {
		panic(err)
	}

	err = app.MakeCmdRoot(filepath.Base(appName)).Execute()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
