package vslog

import (
  "fmt"
  "io"
  "log"
  "os"
  "time"
)

// Constant that define where the logger will be writing
const (
  STDOUT = 1 << iota
  STDERR
  FILE
)

// DateTimeFormat defines the custom date and time format for logs
const DateTimeFormat = "02-01-2006 15:04:05"

// Logger is the logging object
type Logger struct {
  // file        *os.File
  log   *log.Logger
  name  string
  flags int
}

// Info registers a log message with the Info level
func (l *Logger) Info(message string) {
  timestamp := time.Now().Format(DateTimeFormat)

  if l.flags&FILE != 0 {
    lerr := log.New(os.Stderr, "vslog error", log.Ldate|log.Ltime)
    f, err := l.openFile()
    defer func() {
      _ = f.Close()
    }()
    if err != nil {
      lerr.Printf("an error ocurred while logging to file: %v", err)
    }
    _, err = fmt.Fprintf(f, "%s | INFO | %s\n", timestamp, message)
    if err != nil {
      lerr.Printf("an error ocurred while writing the log file: %v", err)
    }
  }

  if l.flags&STDOUT != 0 || l.flags&STDERR != 0 {
    _, _ = fmt.Fprintf(l.log.Writer(), "%s | INFO | %s\n", timestamp, message)
  }
}

// Infof register a log message with string interpolation and Info level
func (l *Logger) Infof(message string, a ...interface{}) {
  l.Info(fmt.Sprintf(message, a...))
}

// Error registers a log message with Error level
func (l *Logger) Error(message string) {
  timestamp := time.Now().Format(DateTimeFormat)

  if l.flags&FILE != 0 {
    lerr := log.New(os.Stderr, "vslog error", log.Ldate|log.Ltime)
    f, err := l.openFile()
    defer func() {
      _ = f.Close()
    }()
    if err != nil {
      lerr.Printf("an error ocurred while logging to file: %v", err)
    }
    _, err = fmt.Fprintf(f, "%s | ERROR | %s\n", timestamp, message)
    if err != nil {
      lerr.Printf("an error ocurred while writing the log file: %v", err)
    }
  }

  if l.flags&STDOUT != 0 || l.flags&STDERR != 0 {
    _, _ = fmt.Fprintf(l.log.Writer(), "%s | ERROR | %s\n", timestamp, message)
  }
}

// Errorf register a log message with string interpolation and Error level
func (l *Logger) Errorf(message string, a ...interface{}) {
  l.Error(fmt.Sprintf(message, a...))
}

func (l *Logger) openFile() (*os.File, error) {
  path := fmt.Sprintf("logs/%s", l.name)
  var err error
  file, err := os.OpenFile(
    fmt.Sprintf("%s/%s.log", path, time.Now().Format("02-01-2006")),
    os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
  if err != nil {
    return nil, fmt.Errorf(
      "an error ocurred while opening/creating the log file: %v", err)
  }
  return file, nil
}

// GetLogger is a function that creates a new logger
func GetLogger(name string, flags int) (*Logger, error) {
  path := fmt.Sprintf("logs/%s", name)
  if _, err := os.Stat(path); os.IsNotExist(err) {
    _ = os.MkdirAll(path, 0755)
  }
  var mw io.Writer
  logger := new(Logger)
  logger.flags = flags
  logger.name = name

  if STDOUT == logger.flags {
    mw = os.Stdout
  } else if STDERR == flags {
    mw = os.Stderr
  } else if STDOUT|STDERR == flags {
    mw = io.MultiWriter(os.Stdout, os.Stderr)
  } else {
    mw = os.Stdout
  }
  logger.log = log.New(mw, "", 0)
  return logger, nil
}
