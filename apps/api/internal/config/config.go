// Package config loads and validates the broker's startup configuration.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Defaults live in code per delivery/configuration: no config file, no config
// library — flags and env vars override these with stdlib only.
const (
	DefaultPort      = "4300"
	DefaultLogFormat = "text"
	DefaultLogLevel  = "info"
	persistDirName   = ".matecito-ai"
)

// Sentinels let callers distinguish which value failed validation regardless
// of wrapping (runtime/error-handling: every package exposes its own
// sentinels for the cases callers must distinguish).
var (
	ErrInvalidPort           = errors.New("config: invalid port")
	ErrPersistDirNotWritable = errors.New("config: persist dir not writable")
	ErrInvalidLogFormat      = errors.New("config: invalid log format")
	ErrInvalidLogLevel       = errors.New("config: invalid log level")
)

// Config is the broker's fully-resolved, validated startup configuration.
type Config struct {
	Port       int
	PersistDir string
	LogFormat  string
	LogLevel   string
}

// Load resolves configuration from flags (highest precedence), then env
// vars, then the hardcoded defaults above, and validates the result before
// returning. args/getenv are injected so callers (tests, or a future CLI
// wrapper embedding the broker) control input instead of relying on process
// globals — delivery/configuration: the config is independent and
// embeddable.
func Load(args []string, getenv func(string) string) (Config, error) {
	fs := flag.NewFlagSet("broker", flag.ContinueOnError)
	portFlag := fs.String("port", "", "HTTP port for the broker")
	persistDirFlag := fs.String("persist-dir", "", "directory for broker persistence")
	logFormatFlag := fs.String("log-format", "", "log format: text|json")
	logLevelFlag := fs.String("log-level", "", "log level: debug|info|warn|error")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	defaultDir, err := defaultPersistDir()
	if err != nil {
		return Config{}, fmt.Errorf("config: resolve default persist dir: %w", err)
	}

	portStr := firstNonEmpty(*portFlag, getenv("BROKER_PORT"), DefaultPort)
	persistDir := firstNonEmpty(*persistDirFlag, getenv("BROKER_PERSIST_DIR"), defaultDir)
	logFormat := firstNonEmpty(*logFormatFlag, getenv("BROKER_LOG_FORMAT"), DefaultLogFormat)
	logLevel := firstNonEmpty(*logLevelFlag, getenv("BROKER_LOG_LEVEL"), DefaultLogLevel)

	port, err := validatePort(portStr)
	if err != nil {
		return Config{}, err
	}
	if err := validatePersistDir(persistDir); err != nil {
		return Config{}, err
	}
	if err := validateLogFormat(logFormat); err != nil {
		return Config{}, err
	}
	if err := validateLogLevel(logLevel); err != nil {
		return Config{}, err
	}

	return Config{
		Port:       port,
		PersistDir: persistDir,
		LogFormat:  logFormat,
		LogLevel:   logLevel,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func defaultPersistDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, persistDirName), nil
}

func validatePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("%w: %q (must be 1-65535)", ErrInvalidPort, portStr)
	}
	return port, nil
}

// validatePersistDir confirms the directory exists (creating it if missing)
// and is actually writable, by creating and removing a probe file. A
// permission-bit check alone would miss a non-directory path blocking the
// location and can be bypassed by a root process.
func validatePersistDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("%w: %q: %v", ErrPersistDirNotWritable, dir, err)
	}

	probe, err := os.CreateTemp(dir, ".broker-write-test-*")
	if err != nil {
		return fmt.Errorf("%w: %q: %v", ErrPersistDirNotWritable, dir, err)
	}
	probePath := probe.Name()
	_ = probe.Close()
	_ = os.Remove(probePath)
	return nil
}

func validateLogFormat(format string) error {
	switch format {
	case "text", "json":
		return nil
	default:
		return fmt.Errorf("%w: %q (must be text or json)", ErrInvalidLogFormat, format)
	}
}

func validateLogLevel(level string) error {
	switch level {
	case "debug", "info", "warn", "error":
		return nil
	default:
		return fmt.Errorf("%w: %q (must be debug, info, warn or error)", ErrInvalidLogLevel, level)
	}
}
