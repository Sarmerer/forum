package logger

import (
	"log"
	"os"
)

func InitLogs(instance string, err error) {
	if err != nil {
		log.Printf("|%s|\t%s: %s\n", Red("ERROR"), instance, err)
		os.Exit(1)
	}
	log.Printf("|%s|\t%s\n", Green("OK"), instance)
}

func HTTPLogs(status, elapsed, host, method, path string) {
	log.Printf("|%s|\t%10s | %s |%s %s", status, elapsed, host, method, path)
}

func CheckErrAndLog(instance, message string, err error) {
	if err != nil {
		log.Printf("|%s|\t%s: %s\n", Red("ERROR"), instance, err)
	} else {
		log.Printf("|%s|\t%s: %s\n", Green("OK"), instance, message)
	}
}

func Log(instance, message string) {
	log.Printf("%s: %s\n", instance, message)
}
