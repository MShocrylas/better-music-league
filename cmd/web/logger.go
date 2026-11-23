package main

import (
	"log/slog"
	"strings"
	"fmt"
	"os"
)


type gooseSlogAdapter struct {
	logger *slog.Logger
}

func (s *gooseSlogAdapter) Printf(format string, v ...any) {
	msg := strings.TrimRight(fmt.Sprintf(format, v...), "\n")
	s.logger.Info(msg)
}

func (s *gooseSlogAdapter) Fatalf(format string, v ...any) {
	msg := strings.TrimRight(fmt.Sprintf(format, v...), "\n")
	s.logger.Error(msg)
	os.Exit(1)
}

