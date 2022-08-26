package stdlogger

import (
	"context"
	"log"
)

type Logger struct{}

func (l Logger) Print(ctx context.Context, payload interface{}) {
	log.Printf("%+v", payload)
}
