package prism

import (
	"fmt"
	"log"
	"os"
)

type PrismLogger struct {
	File        string
	Application string
	LogFile     *os.File
}

func NewPrismLogger(file, app string) {
	tdbRoot := os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}

	//  create and open file
}

func (pl *PrismLogger) Printf(format string, args ...interface{}) {
	pl.log("[INFO]"+format, args...)
}

func (pl *PrismLogger) Errorf(format string, args ...interface{}) {
	pl.log("[ERROR]"+format, args...)
}

func (pl *PrismLogger) log(format string, args ...interface{}) {
	fmt.Fprintf(pl.LogFile, format, args...)
}
