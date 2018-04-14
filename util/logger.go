package util

import (
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	WarningLogger = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
}