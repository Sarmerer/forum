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

func Log(logType string, args ...interface{}) {

}
