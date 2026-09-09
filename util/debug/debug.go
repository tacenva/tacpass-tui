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
	write("DEBUG", message)
}

func Info(message string) {
	write("INFO", message)
}

func Error(message string) {
	write("ERROR", message)
}

func Event(name string, data any) {
	raw, err := json.Marshal(data)
	if err != nil {
		write(
			"EVENT",
			fmt.Sprintf(
				"%s error=%v",
				name,
				err,
			),
		)

		return
	}

	write(
		"EVENT",
		fmt.Sprintf(
			"%s %s",
			name,
			string(raw),
		),
	)
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

	message = format(message)

	_, _ = fmt.Fprintf(
		file,
		"[%s] %s %s\n",
		time.Now().Format(time.RFC3339Nano),
		level,
		message,
	)
}

func format(message string) string {
	var data any

	if err := json.Unmarshal(
		[]byte(message),
		&data,
	); err != nil {
		return message
	}

	formatted, err := json.MarshalIndent(
		data,
		"",
		"  ",
	)
	if err != nil {
		return message
	}

	return string(formatted)
}
