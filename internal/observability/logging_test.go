package observability

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNewJSONHandlerUsesSecondPrecisionTimestamp(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewJSONHandler(&buf, slog.LevelInfo))

	logger.Info("test.log", "k", "v")

	line := strings.TrimSpace(buf.String())
	if line == "" {
		t.Fatal("empty log output")
	}
	entry := map[string]any{}
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("unmarshal log: %v", err)
	}

	timeVal, ok := entry["time"].(string)
	if !ok {
		t.Fatalf("time type = %T, want string", entry["time"])
	}
	timeRaw := strings.TrimSpace(timeVal)
	if timeRaw == "" {
		t.Fatal("time is empty")
	}
	if strings.Contains(timeRaw, ".") {
		t.Fatalf("time includes sub-second precision: %q", timeRaw)
	}
	parsed, err := time.Parse(time.DateTime, timeRaw)
	if err != nil {
		t.Fatalf("time parse error: %v (value=%q)", err, timeRaw)
	}
	if !parsed.Equal(parsed.UTC()) {
		t.Fatalf("time is not UTC: %q", timeRaw)
	}
}
