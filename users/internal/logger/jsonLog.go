package logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)


type JSONLogger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *JSONLogger {
	return &JSONLogger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *JSONLogger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *JSONLogger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *JSONLogger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *JSONLogger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}

func (l *JSONLogger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}
