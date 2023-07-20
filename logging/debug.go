package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var DebugLog *log.Logger

func init() {
	DebugLog = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	DebugLog.SetPrefix("DEBUG: ")
	DebugLog.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Debug(format string, v ...interface{}) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000000Z")
	message := fmt.Sprintf(format, v...)
	DebugLog.Printf("%s %s\n", timestamp, message)
}
