package core

import (
	"log"
	"os"
)

var logFile, _ = os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
var Logger = log.New(logFile, "", log.LstdFlags)
