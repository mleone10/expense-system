package domain

import "context"

type Logger interface {
	Print(context.Context, interface{})
}
