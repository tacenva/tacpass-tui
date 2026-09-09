package debug

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	enabled = true
	path    = "/tmp/debug.log"
)

func Enable() {
	mu.Lock()
	defer mu.Unlock()

	enabled = true
}

func Disable() {
	mu.Lock()
	defer mu.Unlock()

	enabled = false
}

func Print(message string) {
	write("DEBUG", format(message))
}

func Info(message string) {
	write("INFO", format(message))
}

func Error(message string) {
	write("ERROR", format(message))
}

func Event(name string, fields map[string]any) {
	message := name

	for key, value := range fields {
		message += fmt.Sprintf(" %s=%v", key, value)
	}

	write("EVENT", format(message))
}

func format(message string) string {
	var data any

	if err := json.Unmarshal([]byte(message), &data); err != nil {
		return message
	}

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return message
	}

	return string(formatted)
}

func write(level, message string) {
	mu.Lock()
	defer mu.Unlock()

	if !enabled {
		return
	}

	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return
	}
	defer file.Close()

	_, _ = fmt.Fprintf(
		file,
		"[%s] %s %s\n",
		time.Now().Format(time.RFC3339Nano),
		level,
		message,
	)
}
