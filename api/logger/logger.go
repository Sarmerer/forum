package logger

import (
	"forum/api/utils"
	"log"
	"os"
)

func InitLogs(instance string, err error) {
	if err != nil {
		log.Printf("%s:\t|%s| %s\n", instance, utils.Red("ERROR"), err)
		os.Exit(1)
	}
	log.Printf("%s:\t|%s|\n", instance, utils.Green("OK"))
}

func HTTPLogs(status, elapsed, host, method, path string) {
	log.Printf("|%s|\t%10s|\t%s |%s %s", status, elapsed, host, method, path)
}

func ServerLogs(instance string, message string, err error) {
	if err != nil {
		log.Printf("%s:\t|%s| %s\n", instance, utils.Red("ERROR"), err)
	} else {
		log.Printf("%s:\t|%s| %s\n", instance, utils.Green("OK"), message)
	}
}
