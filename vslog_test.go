package vslog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger, err := GetLogger("testlog", STDOUT|STDERR|FILE)
	if err != nil {
		t.Fatalf("error ocurred while creating the logger | error: %v", err)
	}
	defer os.RemoveAll("logs")

	logger.Info("testing Info level log")
	logger.Error("testing Error level log")

	f, err := os.Open(
		fmt.Sprintf("logs/testlog/%s.log", time.Now().Format("02-01-2006")))
	if err != nil {
		t.Fatalf("error ocurred while opening the log file | error: %v", err)
	}
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) != 2 {
		t.Fatalf("wrong number of logged lines | %d vs 2", len(lines))
	}

	if !strings.Contains(lines[0], "testing Info level log") {
		t.Fatalf("Info level logged message is wrong | %s", lines[0])
	}
	if !strings.Contains(lines[1], "testing Error level log") {
		t.Fatalf("Error level logged message is wrong | %s", lines[1])
	}
}
