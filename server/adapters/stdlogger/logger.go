package stdlogger

import (
	"context"
	"encoding/json"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stderr, "", 0),
	}
}

func (l Logger) Print(ctx context.Context, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	log.Println(string(payloadBytes))
	l.logger.Printf("%v", string(payloadBytes))
}
