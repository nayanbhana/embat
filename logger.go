package embat

// Logger prints debug messages.
type Logger interface {
	// Debug prints a debug message.
	Debug(format string, args ...any)
}

// noOpLogger is a logger that does nothing.
type noOpLogger struct{}

// Debug does nothing.
func (noOpLogger) Debug(format string, args ...any) {}
