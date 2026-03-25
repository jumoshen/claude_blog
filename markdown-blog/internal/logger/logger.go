package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Level 日志级别
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// Config 日志配置
type Config struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	Path       string `yaml:"path"`       // 日志目录
	Filename   string `yaml:"filename"`   // 日志文件名
	MaxSize    int    `yaml:"max_size"`   // MB
	MaxBackups int    `yaml:"max_backups"` // 保留备份数
	MaxAge     int    `yaml:"max_age"`    // 天
	Compress   bool   `yaml:"compress"`   // 压缩
	Console    bool   `yaml:"console"`    // 是否输出到控制台
}

// Logger 日志器
type Logger struct {
	mu      sync.Mutex
	config  *Config
	output  io.WriteCloser
	file    *os.File
	stdout  io.Writer
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// New 创建日志器
func New(cfg *Config) *Logger {
	once.Do(func() {
		defaultLogger = &Logger{
			config: cfg,
			stdout: os.Stdout,
		}
		defaultLogger.init()
	})
	return defaultLogger
}

// Default 获取默认日志器
func Default() *Logger {
	if defaultLogger == nil {
		defaultLogger = New(&Config{
			Level:    "debug",
			Console: true,
		})
	}
	return defaultLogger
}

func (l *Logger) init() {
	if l.config.Path != "" && l.config.Filename != "" {
		if err := os.MkdirAll(l.config.Path, 0755); err != nil {
			fmt.Printf("Failed to create log dir: %v\n", err)
		}
		filename := filepath.Join(l.config.Path, l.config.Filename)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
		} else {
			l.file = file
		}
	}
}

func (l *Logger) write(level Level, format string, args ...interface{}) {
	if l.config == nil {
		return
	}

	if !l.shouldLog(level) {
		return
	}

	// 获取调用者信息
	_, file, line, _ := runtime.Caller(2)
	_, fname := filepath.Split(file)

	// 格式化时间
	now := time.Now()

	// 格式化日志
	levelName := "DEBUG"
	if int(level) < len(levelNames) {
		levelName = levelNames[level]
	}

	var body string
	if len(args) == 0 {
		body = format
	} else {
		body = fmt.Sprintf(format, args...)
	}

	// JSON 格式日志
	logEntry := map[string]interface{}{
		"timestamp": now.Format("2006-01-02T15:04:05.000Z07:00"),
		"level":     levelName,
		"file":      fmt.Sprintf("%s:%d", fname, line),
		"message":   body,
	}
	jsonBytes, _ := json.Marshal(logEntry)
	logLine := string(jsonBytes) + "\n"

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.config.Console {
		l.stdout.Write([]byte(logLine))
	}

	if l.file != nil {
		l.file.Write([]byte(logLine))
	}
}

func (l *Logger) shouldLog(level Level) bool {
	levelMap := map[string]Level{
		"debug": DEBUG,
		"info":  INFO,
		"warn":  WARN,
		"error": ERROR,
	}

	minLevel, ok := levelMap[l.config.Level]
	if !ok {
		minLevel = INFO
	}

	return level >= minLevel
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.write(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.write(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.write(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.write(ERROR, format, args...)
}

// Fatal 致命日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.write(FATAL, format, args...)
	os.Exit(1)
}

// Close 关闭日志器
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// WithFields 带字段的日志
func (l *Logger) WithFields(fields map[string]interface{}) *Entry {
	return &Entry{
		logger: l,
		fields: fields,
	}
}

// Entry 日志条目
type Entry struct {
	logger *Logger
	fields map[string]interface{}
}

func (e *Entry) formatMessage(format string, args []interface{}) string {
	parts := make([]string, 0, len(e.fields)+1)
	for k, v := range e.fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	var body string
	if len(args) == 0 {
		body = format
	} else {
		body = fmt.Sprintf(format, args...)
	}
	parts = append(parts, body)

	jsonBytes, _ := json.Marshal(parts)
	return string(jsonBytes)
}

// Debug 调试日志
func (e *Entry) Debug(format string, args ...interface{}) {
	e.logger.write(DEBUG, e.formatMessage(format, args))
}

// Info 信息日志
func (e *Entry) Info(format string, args ...interface{}) {
	e.logger.write(INFO, e.formatMessage(format, args))
}

// Warn 警告日志
func (e *Entry) Warn(format string, args ...interface{}) {
	e.logger.write(WARN, e.formatMessage(format, args))
}

// Error 错误日志
func (e *Entry) Error(format string, args ...interface{}) {
	e.logger.write(ERROR, e.formatMessage(format, args))
}

// 快捷函数
func Debug(format string, args ...interface{}) {
	Default().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	Default().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	Default().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	Default().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	Default().Fatal(format, args...)
}
