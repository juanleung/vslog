# Vslog - Very simple log

## Flags
STDOUT: Logs to the OS Stdout output

STDERR: Logs to the OS Stderr output

FILE: Create a Logs directory and inside a folder with the logger name

## Usage

```go
package main

import "github.com/juanleung/vslog"

func main() {
  logger, err := vslog.GetLogger("LogName", vslog.STDOUT|vslog.FILE)
  // If using the FILE flag use Close to close the file
  defer logger.Close()

  logger.Info("A very simple logger")
}
```