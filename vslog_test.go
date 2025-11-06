package vslog

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger, err := GetLogger(STDOUT|STDERR|FILE, "testlog")
	if err != nil {
		t.Fatalf("error ocurred while creating the logger | error: %v", err)
	}
	defer func() {
		err = os.RemoveAll("logs")
		if err != nil {

		}
	}()

	logger.Info("testing Info level log")
	logger.Error("testing Error level log")

	f, err := os.Open(
		fmt.Sprintf("logs/testlog/%s.log", time.Now().Format("02-01-2006")))
	if err != nil {
		t.Fatalf("error ocurred while opening the log file | error: %v", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

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

func TestLogger_Concurrency(t *testing.T) {
	const (
		loggerName = "testlog_concurrent"
		N          = 200
	)
	defer func() {
		_ = os.RemoveAll("logs")
	}()

	logger, err := GetLogger(STDOUT|STDERR|FILE, loggerName)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	infos := N / 2
	errors := N - infos

	go func() {
		defer wg.Done()
		for i := 0; i < infos; i++ {
			logger.Infof("concurrent info %d", i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < errors; i++ {
			logger.Errorf("concurrent error %d", i)
		}
	}()

	wg.Wait()

	// Dar un pequeño margen para flush de IO de archivo
	time.Sleep(50 * time.Millisecond)

	path := fmt.Sprintf("logs/%s/%s.log", loggerName, time.Now().Format("02-01-2006"))
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed opening log file: %v", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)
	count := 0
	hasInfo := false
	hasError := false
	for scanner.Scan() {
		line := scanner.Text()
		count++
		if strings.Contains(line, "INFO | concurrent info") {
			hasInfo = true
		}
		if strings.Contains(line, "ERROR | concurrent error") {
			hasError = true
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}

	if count != N {
		t.Fatalf("unexpected number of lines: got %d, want %d", count, N)
	}
	if !hasInfo {
		t.Fatalf("missing INFO lines")
	}
	if !hasError {
		t.Fatalf("missing ERROR lines")
	}
}

func TestLogger_HighContention(t *testing.T) {
	const (
		loggerName = "testlog_contention"
		G          = 8   // numbers of goroutines
		M          = 150 // messages per goroutine
	)
	defer func() {
		_ = os.RemoveAll("logs")
	}()

	// increase parallelism to stress the lock
	prev := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(biggest(prev, 2))
	defer runtime.GOMAXPROCS(prev)

	logger, err := GetLogger(STDOUT|STDERR|FILE, loggerName)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(G)
	for g := 0; g < G; g++ {
		gid := g
		go func() {
			defer wg.Done()
			for i := 0; i < M; i++ {
				if i%2 == 0 {
					logger.Info(fmt.Sprintf("g%d info %d", gid, i))
				} else {
					logger.Error(fmt.Sprintf("g%d error %d", gid, i))
				}
			}
		}()
	}
	wg.Wait()

	// give some time to flush the IO
	time.Sleep(100 * time.Millisecond)

	expected := G * M
	path := fmt.Sprintf("logs/%s/%s.log", loggerName, time.Now().Format("02-01-2006"))
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed opening log file: %v", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}

	if count != expected {
		t.Fatalf("unexpected number of lines: got %d, want %d", count, expected)
	}
}

func biggest(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestLogger_WithoutName(t *testing.T) {
	defer func() {
		_ = os.RemoveAll("logs")
	}()

	logger, err := GetLogger(STDOUT | FILE)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// The logger name should be the caller function name
	expectedName := "vslog.TestLogger_WithoutName"
	if logger.name != expectedName {
		t.Fatalf("unexpected logger name: got %q, want %q", logger.name, expectedName)
	}
}
