package mover

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/logutils"
)

var logger *log.Logger

func init() {
	initLogger()
}

func initLogger() {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "" {
		level = "info"
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"debug", "info", "warn", "error", "fatal"},
		MinLevel: logutils.LogLevel(level),
		Writer:   os.Stderr,
	}

	log.SetOutput(filter)
}
