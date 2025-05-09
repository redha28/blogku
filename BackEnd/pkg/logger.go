package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Log levels
const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

// Logger represents the logging system
type Logger struct {
	level      int
	errorLog   *log.Logger
	warnLog    *log.Logger
	infoLog    *log.Logger
	debugLog   *log.Logger
	fileHandle *os.File
}

var (
	// Default logger instance
	defaultLogger *Logger
)

// InitLogger creates a new logger with the specified level and log file path
func InitLogger(level int, logFilePath string) (*Logger, error) {
	logger := &Logger{
		level: level,
	}

	// Ensure log directory exists
	if logFilePath != "" {
		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open the log file, creating it if it doesn't exist
		fileHandle, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// Save file handle for later closure
		logger.fileHandle = fileHandle

		// Set up multi-writer for both console and file
		multiWriter := io.MultiWriter(os.Stdout, fileHandle)

		// Initialize loggers with appropriate prefixes
		logger.errorLog = log.New(multiWriter, "\033[31mERROR\033[0m: ", log.Ldate|log.Ltime)
		logger.warnLog = log.New(multiWriter, "\033[33mWARN\033[0m:  ", log.Ldate|log.Ltime)
		logger.infoLog = log.New(multiWriter, "\033[36mINFO\033[0m:  ", log.Ldate|log.Ltime)
		logger.debugLog = log.New(multiWriter, "\033[35mDEBUG\033[0m: ", log.Ldate|log.Ltime)
	} else {
		// Initialize loggers with console output only
		logger.errorLog = log.New(os.Stdout, "\033[31mERROR\033[0m: ", log.Ldate|log.Ltime)
		logger.warnLog = log.New(os.Stdout, "\033[33mWARN\033[0m:  ", log.Ldate|log.Ltime)
		logger.infoLog = log.New(os.Stdout, "\033[36mINFO\033[0m:  ", log.Ldate|log.Ltime)
		logger.debugLog = log.New(os.Stdout, "\033[35mDEBUG\033[0m: ", log.Ldate|log.Ltime)
	}

	// Set as default logger if none exists yet
	if defaultLogger == nil {
		defaultLogger = logger
	}

	return logger, nil
}

// Close closes the log file if one was opened
func (l *Logger) Close() error {
	if l.fileHandle != nil {
		return l.fileHandle.Close()
	}
	return nil
}

// LogWithCaller adds file and line information to the log message
func logWithCaller(msg string) string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	// Extract just the filename without the full path
	file = filepath.Base(file)
	return fmt.Sprintf("%s (%s:%d)", msg, file, line)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error) {
	if l.level >= LevelError {
		if err != nil {
			l.errorLog.Println(logWithCaller(fmt.Sprintf("%s: %v", msg, err)))
		} else {
			l.errorLog.Println(logWithCaller(msg))
		}
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	if l.level >= LevelWarn {
		l.warnLog.Println(logWithCaller(msg))
	}
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	if l.level >= LevelInfo {
		l.infoLog.Println(msg)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	if l.level >= LevelDebug {
		l.debugLog.Println(logWithCaller(msg))
	}
}

// ErrorWithFields logs an error with additional fields
func (l *Logger) ErrorWithFields(msg string, err error, fields map[string]interface{}) {
	if l.level >= LevelError {
		fieldsStr := formatFields(fields)
		if err != nil {
			l.errorLog.Println(logWithCaller(fmt.Sprintf("%s: %v %s", msg, err, fieldsStr)))
		} else {
			l.errorLog.Println(logWithCaller(fmt.Sprintf("%s %s", msg, fieldsStr)))
		}
	}
}

// InfoWithFields logs info with additional fields
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
	if l.level >= LevelInfo {
		fieldsStr := formatFields(fields)
		l.infoLog.Println(fmt.Sprintf("%s %s", msg, fieldsStr))
	}
}

// formatFields formats a map of fields into a string
func formatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	result := "{"
	for k, v := range fields {
		result += fmt.Sprintf("%s:%v, ", k, v)
	}
	// Remove the trailing comma and space
	result = result[:len(result)-2] + "}"
	return result
}

// GetLogger returns the default logger
func GetLogger() *Logger {
	return defaultLogger
}

// Error logs an error message using the default logger
func Error(msg string, err error) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, err)
	}
}

// Warn logs a warning message using the default logger
func Warn(msg string) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg)
	}
}

// Info logs an informational message using the default logger
func Info(msg string) {
	if defaultLogger != nil {
		defaultLogger.Info(msg)
	}
}

// Debug logs a debug message using the default logger
func Debug(msg string) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg)
	}
}

// ErrorWithFields logs an error with fields using the default logger
func ErrorWithFields(msg string, err error, fields map[string]any) {
	if defaultLogger != nil {
		defaultLogger.ErrorWithFields(msg, err, fields)
	}
}

// LogHTTPRequest logs details about an HTTP request with its processing time
func LogHTTPRequest(method, path, ip string, status int, duration time.Duration) {
	if defaultLogger == nil || defaultLogger.level < LevelInfo {
		return
	}

	fields := map[string]interface{}{
		"method":   method,
		"path":     path,
		"ip":       ip,
		"status":   status,
		"duration": duration,
	}
	defaultLogger.InfoWithFields("HTTP Request", fields)
}

// LogPanic logs a panic and recovers from it
func LogPanic() {
	if r := recover(); r != nil {
		// Get stack trace
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		stackTrace := string(buf[:n])

		fields := map[string]any{
			"stack": stackTrace,
		}
		ErrorWithFields("PANIC RECOVERED", fmt.Errorf("%v", r), fields)
	}
}
