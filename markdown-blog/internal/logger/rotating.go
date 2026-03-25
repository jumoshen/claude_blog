package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// RotatingWriter 旋转日志写入器
type RotatingWriter struct {
	mu       sync.Mutex
	filename string
	writer   io.WriteCloser
}

// NewRotatingWriter 创建旋转日志写入器
func NewRotatingWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) *RotatingWriter {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create log dir: %v\n", err)
	}

	w := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	return &RotatingWriter{
		filename: filename,
		writer:   w,
	}
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.writer.Write(p)
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.writer != nil {
		return w.writer.Close()
	}
	return nil
}

// FileWriter 文件写入器（简单版，不旋转）
type FileWriter struct {
	mu    sync.Mutex
	file  *os.File
	path  string
}

// NewFileWriter 创建文件写入器
func NewFileWriter(path string) (*FileWriter, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file: file,
		path: path,
	}, nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Write(p)
}

func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func (w *FileWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// DailyWriter 按天分割的写入器
type DailyWriter struct {
	mu       sync.Mutex
	filename string
	writer   io.WriteCloser
	date     string
}

// NewDailyWriter 创建按天分割的写入器
func NewDailyWriter(filename string) *DailyWriter {
	w := &DailyWriter{
		filename: filename,
	}
	w.rotate()
	return w
}

func (w *DailyWriter) rotate() {
	// 检查是否需要切换
	today := time.Now().Format("2006-01-02")
	if w.date == today && w.writer != nil {
		return
	}

	// 关闭旧的 writer
	if w.writer != nil {
		w.writer.Close()
	}

	// 创建新的文件名
	newFilename := fmt.Sprintf("%s.%s.log", w.filename, today)
	file, err := os.OpenFile(newFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open daily log file: %v\n", err)
		return
	}

	w.writer = file
	w.date = today
}

func (w *DailyWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 检查是否需要切换日期
	w.rotate()

	if w.writer != nil {
		return w.writer.Write(p)
	}
	return 0, nil
}

func (w *DailyWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.writer != nil {
		return w.writer.Close()
	}
	return nil
}
