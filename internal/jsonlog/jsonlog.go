package jsonlog

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

//Level defines the log level
type Level int8

//Represent specfic severity levels
const (
	LevelInfo Level = iota //value 0
	LevelError
	LevelFatal
	LevelOff
)

//String human readable representation of the level
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelOff:
		return "OFF"
	default:
		return ""
	}
}

//Logger type specifies the output and the severity level of the logs.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(msg string, props map[string]string) {
	l.print(LevelInfo, msg, props)
}

func (l *Logger) PrintFatal(err error, props map[string]string) {
	l.print(LevelFatal, err.Error(), props)
}

func (l *Logger) PrintError(err error, props map[string]string) {
	l.print(LevelError, err.Error(), props)
}

func (l *Logger) print(level Level, msg string, props map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Time       string            `json:"time"`
		Trace      string            `json:"trace,omiitempty"`
	}{
		Level:      level.String(),
		Message:    msg,
		Properties: props,
		Time:       time.Now().UTC().Format(time.RFC3339),
	}

	//include stacktrace for errors
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte("{\"error\":\"error marshalling log message:\"}" + err.Error())
	}

	//adquire lock for writing
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

// We also implement a Write() method on our Logger type so that it satisfies the
// io.Writer interface. This writes a log entry at the ERROR level with no additional
// properties.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
