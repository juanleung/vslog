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

// Logger is the logging object
type Logger struct {
	file *os.File
	log  *log.Logger
}

// Close the file is exist
func (l *Logger) Close() error {
	if l.file != nil {
		err := l.file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetLogger is a function that create a new logger
func GetLogger(name string, flags int) (*Logger, error) {
	path := fmt.Sprintf("logs/%s", name)
	var mw io.Writer
	logger := new(Logger)

	if flags&FILE != 0 {
		var err error
		logger.file, err = os.OpenFile(
			fmt.Sprintf("%s/%s.log", path, time.Now().Format("02-01-2006")),
			os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)

		if err != nil {
			return nil, fmt.Errorf(
				"an error ourred while opening/creating the file: %v", err)
		}
	}

	if STDOUT == flags {
		mw = os.Stdout
	} else if STDERR == flags {
		mw = os.Stderr
	} else if FILE == flags {
		mw = logger.file
	} else if STDOUT|STDERR == flags {
		mw = io.MultiWriter(os.Stdout, os.Stderr)
	} else if STDOUT|FILE == flags {
		mw = io.MultiWriter(os.Stdout, logger.file)
	} else if STDERR|FILE == flags {
		mw = io.MultiWriter(os.Stderr, logger.file)
	} else if STDOUT|STDERR|FILE == flags {
		mw = io.MultiWriter(os.Stdout, os.Stderr, logger.file)
	} else {
		mw = os.Stdout
	}
	logger.log = log.New(mw, "", log.Ldate|log.Ltime)
	return logger, nil
}
