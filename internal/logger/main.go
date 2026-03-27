package logger

import (
	"fmt"
	"log"
	"os"
)

func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
}
