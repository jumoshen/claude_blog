package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRotatingWriter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rotating_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.log")
	writer := NewRotatingWriter(filename, 100, 7, 30, true)

	if writer == nil {
		t.Fatal("Expected writer to be created")
	}
	defer writer.Close()
}

func TestRotatingWriter_Write(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rotating_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.log")
	writer := NewRotatingWriter(filename, 100, 7, 30, true)
	defer writer.Close()

	data := []byte("test log entry\n")
	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}
}

func TestRotatingWriter_Close(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rotating_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.log")
	writer := NewRotatingWriter(filename, 100, 7, 30, true)

	// Write some data
	writer.Write([]byte("test\n"))

	// Close should not error
	if err := writer.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Second close should be safe
	if err := writer.Close(); err != nil {
		t.Errorf("Second close failed: %v", err)
	}
}

func TestNewFileWriter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.log")
	writer, err := NewFileWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()
}

func TestNewFileWriter_CreatesDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Try to create writer in non-existent subdirectory
	filename := filepath.Join(tmpDir, "subdir", "test.log")
	writer, err := NewFileWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	// Check directory was created
	if _, err := os.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		t.Error("Expected directory to be created")
	}
}

func TestFileWriter_Write(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.log")
	writer, err := NewFileWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	data := []byte("test log entry\n")
	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}
}

func TestFileWriter_Sync(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.log")
	writer, err := NewFileWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	if err := writer.Sync(); err != nil {
		t.Errorf("Sync failed: %v", err)
	}
}

func TestFileWriter_Close(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filewriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.log")
	writer, err := NewFileWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}

	writer.Write([]byte("test\n"))

	if err := writer.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestNewDailyWriter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dailywriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test")
	writer := NewDailyWriter(filename)
	defer writer.Close()

	if writer == nil {
		t.Fatal("Expected writer to be created")
	}
}

func TestDailyWriter_Write(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dailywriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test")
	writer := NewDailyWriter(filename)
	defer writer.Close()

	data := []byte("test log entry\n")
	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}
}

func TestDailyWriter_Close(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dailywriter_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test")
	writer := NewDailyWriter(filename)

	writer.Write([]byte("test\n"))

	if err := writer.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
