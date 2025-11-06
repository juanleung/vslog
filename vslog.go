package vslog

import (
  "fmt"
  "io"
  "log"
  "os"
  "path/filepath"
  "runtime"
  "sync"
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
  mu    sync.Mutex
  lErr  *log.Logger // internal logger for errors
}

func (l *Logger) logWithLevel(level, message string) {
  l.mu.Lock()
  defer l.mu.Unlock()

  timestamp := time.Now().Format(DateTimeFormat)

  if l.flags&FILE != 0 {
    f, err := l.openFile()
    if err != nil {
      l.lErr.Printf("an error occurred while logging to file: %v", err)
    } else {
      defer func() {
        _ = f.Close()
      }()
      if _, err = fmt.Fprintf(f, "%s | %s | %s\n", timestamp, level, message); err != nil {
        l.lErr.Printf("an error occurred while writing the log file: %v", err)
      }
    }
  }

  // Escribir a las salidas estándar configuradas
  if l.flags&(STDOUT|STDERR) != 0 {
    _, _ = fmt.Fprintf(
      l.log.Writer(), "%s | %s | %s\n", timestamp, level, message)
  }
}

// Debug registers a log message with the Debug level
func (l *Logger) Debug(message string) {
  l.logWithLevel("DEBUG", message)
}

// Debugf register a log message with string interpolation and Debug level
func (l *Logger) Debugf(message string, a ...interface{}) {
  l.Debug(fmt.Sprintf(message, a...))
}

// Info registers a log message with the Info level
func (l *Logger) Info(message string) {
  l.logWithLevel("INFO", message)
}

// Infof register a log message with string interpolation and Info level
func (l *Logger) Infof(message string, a ...interface{}) {
  l.Info(fmt.Sprintf(message, a...))
}

// Warning registers a log message with the Warning level
func (l *Logger) Warning(message string) {
  l.logWithLevel("WARNING", message)
}

// Warningf register a log message with string interpolation and Warning level
func (l *Logger) Warningf(message string, a ...interface{}) {
  l.Info(fmt.Sprintf(message, a...))
}

// Error registers a log message with Error level
func (l *Logger) Error(message string) {
  l.logWithLevel("ERROR", message)
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
      "an error occurred while opening/creating the log file: %v", err)
  }
  return file, nil
}

// GetLogger is a function that creates a new logger
// If name is not provided, it uses the caller function name
func GetLogger(flags int, name ...string) (*Logger, error) {
  var loggerName string
  if len(name) > 0 && name[0] != "" {
    loggerName = name[0]
  } else {
    loggerName = getCallerName()
  }

  path := fmt.Sprintf("logs/%s", loggerName)
  if _, err := os.Stat(path); os.IsNotExist(err) {
    _ = os.MkdirAll(path, 0755)
  }

  var writers []io.Writer
  if flags&STDOUT != 0 {
    writers = append(writers, os.Stdout)
  }
  if flags&STDERR != 0 {
    writers = append(writers, os.Stderr)
  }
  if len(writers) == 0 {
    writers = append(writers, os.Stdout)
  }
  mw := writers[0]
  if len(writers) > 1 {
    mw = io.MultiWriter(writers...)
  }

  logger := &Logger{
    log:   log.New(mw, "", 0),
    name:  loggerName,
    flags: flags,
    lErr:  log.New(os.Stderr, "vslog error: ", log.Ldate|log.Ltime),
  }
  return logger, nil
}

// getCallerName returns the name of the function that called GetLogger
func getCallerName() string {
  pc, _, _, ok := runtime.Caller(2)
  if !ok {
    return "unknown"
  }

  fn := runtime.FuncForPC(pc)
  if fn == nil {
    return "unknown"
  }

  fullName := fn.Name()
  name := filepath.Base(fullName)
  return name
}
