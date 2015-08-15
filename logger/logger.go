package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	db "github.com/dragonrider23/inca/database"
)

var (
	logLevelCodes = map[string]int{
		"emergency": 0,
		"alert":     1,
		"critical":  2,
		"error":     3,
		"warning":   4,
		"notice":    5,
		"info":      6,
		"debug":     7,
	}
)

// Logger type
type Logger struct {
	LastLog  string
	filename string
	Area     string
}

// New creates a new logger object
func New(a string) *Logger {
	return &Logger{
		filename: "logs/all-logs.log",
		Area:     a,
	}
}

// Fatal write and Emergency level log and exits
func (l *Logger) Fatal(f string, v ...interface{}) {
	l.Log("Emergency", f, v...)
	os.Exit(1)
}

// Emergency log level
func (l *Logger) Emergency(f string, v ...interface{}) {
	l.Log("Emergency", f, v...)
}

// Alert log level
func (l *Logger) Alert(f string, v ...interface{}) {
	l.Log("Alert", f, v...)
}

// Critical log level
func (l *Logger) Critical(f string, v ...interface{}) {
	l.Log("Critical", f, v...)
}

// Error log level
func (l *Logger) Error(f string, v ...interface{}) {
	l.Log("Error", f, v...)
}

// Warning log level
func (l *Logger) Warning(f string, v ...interface{}) {
	l.Log("Warning", f, v...)
}

// Notice log level
func (l *Logger) Notice(f string, v ...interface{}) {
	l.Log("Notice", f, v...)
}

// Info log level
func (l *Logger) Info(f string, v ...interface{}) {
	l.Log("Info", f, v...)
}

// Debug log level
func (l *Logger) Debug(f string, v ...interface{}) {
	l.Log("Debug", f, v...)
}

// Log - Generic log function
func (l *Logger) Log(level string, f string, v ...interface{}) {
	log := fmt.Sprintf(f, v...)
	now := time.Now()
	nowString := now.Format(time.RFC822Z)
	message := nowString + " " + level + ": " + l.Area + ": " + log

	l.writeLogToFile(message + "\n")
	l.writeLogToDB(level, now, log)
	fmt.Println(message)
	l.LastLog = message
}

func (l *Logger) writeLogToFile(m string) {
	f, err := os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0665)
	if err != nil {
		fmt.Println("Failed to open log file for writing")
	}
	defer f.Close()
	f.WriteString(m)
}

func (l *Logger) writeLogToDB(level string, t time.Time, m string) {
	if !db.Ready {
		return
	}

	level = strings.ToLower(level)
	levelCode := logLevelCodes[level]

	db.Conn.Exec(`INSERT INTO logs
        VALUES (null, ?, ?, ?, ?)`,
		l.Area,
		levelCode,
		t.Unix(),
		m)
}
