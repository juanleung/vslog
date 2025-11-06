# Vslog - Very simple log

![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)

## Flags
STDOUT: Logs to the OS Stdout output

STDERR: Logs to the OS Stderr output

FILE: Create a Logs directory and inside a folder with the logger name

## Usage

```go
package main

import "github.com/juanleung/vslog"

func main() {
  logger, err := vslog.GetLogger(vslog.STDOUT|vslog.FILE, "loggerName" /* optional */)
  logger.Info("A very simple logger")
}
```