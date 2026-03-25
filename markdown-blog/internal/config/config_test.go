package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	content := `
server:
  host: "127.0.0.1"
  port: 8080
site:
  title: "Test Blog"
  author: "Test Author"
  description: "Test Description"
  beian: ""
  keywords: "test,blog"
database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  database: "testdb"
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
log:
  level: "debug"
  path: "./logs"
  filename: "test.log"
  max_size: 100
  max_backups: 7
  max_age: 30
  compress: true
  console: true
github:
  client_id: "test_client_id"
  client_secret: "test_client_secret"
  callback_url: "http://localhost:8080/callback"
content:
  path: "./content"
  base_url: "http://localhost:8080"
jwt:
  secret: "test_secret"
  expiration: 3600
  issuer: "test_issuer"
`
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test server config
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Expected server host '127.0.0.1', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected server port 8080, got %d", cfg.Server.Port)
	}

	// Test site config
	if cfg.Site.Title != "Test Blog" {
		t.Errorf("Expected site title 'Test Blog', got '%s'", cfg.Site.Title)
	}
	if cfg.Site.Author != "Test Author" {
		t.Errorf("Expected site author 'Test Author', got '%s'", cfg.Site.Author)
	}

	// Test database config
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected database host 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != 3306 {
		t.Errorf("Expected database port 3306, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "root" {
		t.Errorf("Expected database user 'root', got '%s'", cfg.Database.User)
	}
	if cfg.Database.Database != "testdb" {
		t.Errorf("Expected database name 'testdb', got '%s'", cfg.Database.Database)
	}

	// Test redis config
	if cfg.Redis.Host != "localhost" {
		t.Errorf("Expected redis host 'localhost', got '%s'", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 6379 {
		t.Errorf("Expected redis port 6379, got %d", cfg.Redis.Port)
	}

	// Test log config
	if cfg.Log.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.Log.Level)
	}
	if cfg.Log.MaxSize != 100 {
		t.Errorf("Expected log max_size 100, got %d", cfg.Log.MaxSize)
	}

	// Test github config
	if cfg.Github.ClientID != "test_client_id" {
		t.Errorf("Expected github client_id 'test_client_id', got '%s'", cfg.Github.ClientID)
	}
	if cfg.Github.CallbackURL != "http://localhost:8080/callback" {
		t.Errorf("Expected github callback_url 'http://localhost:8080/callback', got '%s'", cfg.Github.CallbackURL)
	}

	// Test jwt config
	if cfg.JWT.Secret != "test_secret" {
		t.Errorf("Expected jwt secret 'test_secret', got '%s'", cfg.JWT.Secret)
	}
	if cfg.JWT.Expiration != 3600 {
		t.Errorf("Expected jwt expiration 3600, got %d", cfg.JWT.Expiration)
	}
	if cfg.JWT.Issuer != "test_issuer" {
		t.Errorf("Expected jwt issuer 'test_issuer', got '%s'", cfg.JWT.Issuer)
	}

	// Test content config
	if cfg.Content.Path != "./content" {
		t.Errorf("Expected content path './content', got '%s'", cfg.Content.Path)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		Database: "testdb",
	}

	expected := "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := cfg.GetDSN()
	if dsn != expected {
		t.Errorf("Expected DSN '%s', got '%s'", expected, dsn)
	}
}

func TestRedisConfig_GetRedisAddr(t *testing.T) {
	cfg := RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	expected := "localhost:6379"
	addr := cfg.GetRedisAddr()
	if addr != expected {
		t.Errorf("Expected addr '%s', got '%s'", expected, addr)
	}
}
