package util

import (
	"log"
	"os"
)

var (
	// InfoLogger is a custom logger that adds the [INFO] prefix and logs to stdout
	InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)

	// WarningLogger is a custom logger that adds the [WARNING] prefix and logs to stdout
	WarningLogger = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime)

	// ErrorLogger is a custom logger that adds the [ERROR] prefix and logs to stderr
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
)
