package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/figure/cmd"
)

func main() {
	log.SetReportTimestamp(false)
	if os.Getenv("MODE") == "dev" {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
	cmd.Execute()
}
