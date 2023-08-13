package binance

import (
	"fmt"
	"go-binance-api/config"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logChannel chan string
var bInit = false
var logTargetChannel chan string
var bInitTarget = false

func TargetSetup(filename string) {
	logTargetChannel = make(chan string, 100000)

	go func() {
		if !config.IsLogToFile() {
			log.SetFlags(0)
			log.SetOutput(os.Stdout)
			for msg := range logTargetChannel {
				log.Print(msg)
			}
		} else {
			f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			for msg := range logTargetChannel {
				f.WriteString(msg)
			}
			f.Close()
		}
	}()
	bInitTarget = true
}

func TargetOutlog(format string, a ...any) {
	var message string
	if a == nil {
		message = format
	} else {
		message = fmt.Sprintf(format, a...)
	}

	if bInitTarget {
		logTargetChannel <- fmt.Sprintf("%s %s\n", time.Now().UTC().Format(time.RFC3339Nano), message)
	} else {
		println(fmt.Sprintf("%s %s\n", time.Now().UTC().Format(time.RFC3339Nano), message))
	}
}

func Setup(filename string) {
	logChannel = make(chan string, 100000)

	go func() {
		if !config.IsLogToFile() {
			log.SetFlags(0)
			log.SetOutput(os.Stdout)
			for msg := range logChannel {
				log.Print(msg)
			}
		} else {
			f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				println(filename)
				panic(err)
			}
			for msg := range logChannel {
				f.WriteString(msg)
			}
			f.Close()
		}
	}()
	bInit = true
}

func Outlog(format string, a ...any) {
	var message string
	if a == nil {
		message = format
	} else {
		message = fmt.Sprintf(format, a...)
	}

	if bInit {
		logChannel <- fmt.Sprintf("%s %s\n", time.Now().UTC().Format(time.RFC3339Nano), message)
	} else {
		println(fmt.Sprintf("%s %s\n", time.Now().UTC().Format(time.RFC3339Nano), message))
	}
}

func WriteFile(filename string, body []byte) {
	isTestRun := false
	if os.Getenv("GO_ENV") == "testing" {
		isTestRun = true
	}
	if isTestRun {
		filename = filepath.Join("..", "..", "logs", filename)
	} else {
		filename = filepath.Join("logs", filename)
	}
	err := os.WriteFile(filename, body, 0)
	if err != nil {
		panic(err)
	}
}
