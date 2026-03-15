package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/mek/go-envdir"
)

const (
	defaultEnvDir = "./env"
	defaultHost   = "127.0.0.1"
	defaultPort   = "8080"
	defaultDir    = "."
)

func main() {
	envdirPath := defaultEnvDir
	if len(os.Args) > 1 {
		envdirPath = os.Args[1]
	}

	env := os.Environ()
	if entries, err := envdir.Read(envdirPath); err == nil {
		env = envdir.Apply(env, entries)
	} else if !os.IsNotExist(err) {
		log.Fatal(err)
	}

	cfg := configFromEnv(env)
	if err := validateConfig(cfg); err != nil {
		log.Fatal(err)
	}

	addr := net.JoinHostPort(cfg.host, cfg.port)
	log.Printf("serving %s on http://%s", cfg.dir, addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(cfg.dir))))
}

type config struct {
	host string
	port string
	dir  string
}

func configFromEnv(env []string) config {
	values := make(map[string]string, len(env))
	for _, item := range env {
		name, value, ok := splitEnv(item)
		if !ok {
			continue
		}
		values[name] = value
	}

	return config{
		host: getOrDefault(values, "SERVER_HOST", defaultHost),
		port: getOrDefault(values, "SERVER_PORT", defaultPort),
		dir:  getOrDefault(values, "SERVER_DIR", defaultDir),
	}
}

func validateConfig(cfg config) error {
	if cfg.host == "" {
		return fmt.Errorf("SERVER_HOST must not be empty")
	}
	if cfg.port == "" {
		return fmt.Errorf("SERVER_PORT must not be empty")
	}
	if _, err := net.LookupPort("tcp", cfg.port); err != nil {
		return fmt.Errorf("invalid SERVER_PORT %q: %w", cfg.port, err)
	}
	info, err := os.Stat(cfg.dir)
	if err != nil {
		return fmt.Errorf("SERVER_DIR %q: %w", cfg.dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("SERVER_DIR %q is not a directory", cfg.dir)
	}
	return nil
}

func getOrDefault(values map[string]string, name, fallback string) string {
	if value, ok := values[name]; ok && value != "" {
		return value
	}
	return fallback
}

func splitEnv(item string) (string, string, bool) {
	for i := 0; i < len(item); i++ {
		if item[i] == '=' {
			return item[:i], item[i+1:], true
		}
	}
	return "", "", false
}
